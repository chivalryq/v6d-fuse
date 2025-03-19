package v6dfs

import (
	"context"
	"log"
	"path"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/chivalryq/v6d-fuse/backend"
	"github.com/chivalryq/v6d-fuse/internal"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

// V6dRoot is the root node of the v6d file system.
// Its functionality is to store a global backend.
type V6dRoot struct {
	V6dNode
	backend backend.Backend
	cache   internal.Cache
}

func NewV6dRoot(backend backend.Backend, cache internal.Cache) *V6dRoot {
	r := &V6dRoot{
		V6dNode: V6dNode{
			path:  "/",
			isDir: true,
		},
		backend: backend,
		cache:   cache,
	}

	r.root = r

	return r
}

type V6dNode struct {
	fs.Inode

	root *V6dRoot
	path string

	mu sync.Mutex

	// is directory
	isDir bool

	cacheID uint64

	// info fuse.Attr
	mtime time.Time
}

func (v *V6dNode) NewChild(name string, isDir bool) *fs.Inode {
	ops := V6dNode{
		root:  v.root,
		mu:    sync.Mutex{},
		path:  v.subPath(name),
		isDir: isDir,
	}

	mode := uint32(syscall.S_IFDIR)
	if !isDir {
		mode = uint32(syscall.S_IFREG)
	}

	ch := v.NewInode(context.Background(), &ops, fs.StableAttr{
		Mode: mode,
	})
	ok := v.AddChild(name, ch, false)
	if !ok {
		log.Printf("Add child failed: %v", name)
	}

	return ch
}

// Node types must be InodeEmbedders
var _ = (fs.InodeEmbedder)((*V6dNode)(nil))

func (n *V6dNode) EmbeddedInode() *fs.Inode {
	return &n.Inode
}

// Node types should implement some file system operations, eg. Lookup
var _ = (fs.NodeLookuper)((*V6dNode)(nil))

func (n *V6dNode) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	// Check if it's a file
	exists, err := n.root.backend.Exists(n.subPath(name))
	if err != nil {
		return nil, syscall.ENOENT
	}

	if child, ok := n.Children()[name]; ok {
		return child, 0
	}

	if exists {
		ch := n.NewChild(name, false)
		return ch, 0
	}

	// Check if it's a directory
	children, err := n.root.backend.List(n.subPath(name))
	if err != nil {
		return nil, syscall.ENOENT // io error?
	}
	if len(children) == 0 {
		return nil, syscall.ENOENT
	}

	ch := n.NewChild(name, true)
	return ch, 0
}

// Ensure we are implementing the NodeReaddirer interface
var _ = (fs.NodeReaddirer)((*V6dNode)(nil))

// Readdir is part of the NodeReaddirer interface
func (n *V6dNode) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	r := make([]fuse.DirEntry, 0)
	children, err := n.root.backend.List(n.path)
	if err != nil {
		return nil, syscall.ENOENT
	}
	for _, child := range children {
		firstPart, isDir := getFirstPathPart(strings.TrimPrefix(child, n.path))
		if firstPart == "" {
			continue
		}

		if existChild, ok := n.Children()[firstPart]; ok {
			r = append(r, fuse.DirEntry{
				Name: firstPart,
				Ino:  existChild.StableAttr().Ino,
				Mode: existChild.StableAttr().Mode,
			})
		} else {
			ch := n.NewChild(firstPart, isDir)
			r = append(r, fuse.DirEntry{
				Name: firstPart,
				Ino:  ch.StableAttr().Ino,
				Mode: ch.StableAttr().Mode,
			})
		}
	}
	return fs.NewListDirStream(r), 0
}

func getFirstPathPart(path string) (part string, isDir bool) {
	clean := strings.TrimPrefix(path, "/")
	parts := strings.Split(clean, "/")
	if len(parts) == 0 {
		return "", false
	}
	part = parts[0]
	if len(parts) == 1 {
		isDir = true
	}
	return
}

func (n *V6dNode) subPath(name string) string {
	return path.Join(n.path, name)
}

// Implement GetAttr to provide size and mtime
var _ = (fs.NodeGetattrer)((*V6dNode)(nil))

func (bn *V6dNode) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	bn.mu.Lock()
	defer bn.mu.Unlock()
	bn.getattr(out)
	log.Printf("Getattr: %v", out.Size)
	return 0
}

func (bn *V6dNode) getattr(out *fuse.AttrOut) {
	if bn.isDir {
		out.Mode = fuse.S_IFDIR
		out.Size = 0
	} else {
		out.Mode = fuse.S_IFREG
		content, err := bn.root.backend.Get(bn.path)
		if err != nil {
			return
		}
		out.Size = uint64(len(content))
	}

	out.SetTimes(nil, &bn.mtime, nil)
}

// Implement Setattr to support truncation
var _ = (fs.NodeSetattrer)((*V6dNode)(nil))

func (bn *V6dNode) Setattr(ctx context.Context, fh fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	bn.mu.Lock()
	defer bn.mu.Unlock()

	// if sz, ok := in.GetSize(); ok {
	// 	bn.resize(sz)
	// }
	bn.getattr(out)
	return 0
}

// Implement (handleless) Open
var _ = (fs.NodeOpener)((*V6dNode)(nil))

func (f *V6dNode) Open(ctx context.Context, openFlags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	return nil, 0, 0
}

// Implement tF read.
var _ = (fs.NodeReader)((*V6dNode)(nil))

func (bn *V6dNode) Read(ctx context.Context, fh fs.FileHandle, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	bn.mu.Lock()
	defer bn.mu.Unlock()

	var cacheErr error
	if bn.cacheID != 0 {
		content, err := bn.root.cache.Get(bn.cacheID)
		if err == nil {
			return fuse.ReadResultData(content), 0
		}
		log.Printf("Retrive from cache failed: %v, skip", err)
		cacheErr = err
	}

	content, err := bn.root.backend.Get(bn.path)
	if err != nil {
		return nil, syscall.ENOENT
	}

	if bn.cacheID == 0 || cacheErr != nil { // no cache or cache failed, put to cache
		bn.cacheID, cacheErr = bn.root.cache.Put(content)
		if cacheErr != nil {
			log.Printf("Put to cache failed: %v, skip", cacheErr)
		}
	}

	end := off + int64(len(dest))
	if end > int64(len(content)) {
		end = int64(len(content))
	}

	// We could copy to the `dest` buffer, but since we have a
	// []byte already, return that.
	return fuse.ReadResultData(content[off:end]), 0
}