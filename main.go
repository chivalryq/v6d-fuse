package main

import (
	"log"


	"github.com/spf13/pflag"

	"github.com/chivalryq/v6d-fuse/backend"
	"github.com/chivalryq/v6d-fuse/internal"
	"github.com/chivalryq/v6d-fuse/v6dfs"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

func main() {
	internal.AddFlags(pflag.CommandLine)
	pflag.Parse()
	args := internal.GetArgs()

	if len(pflag.Args()) < 1 {
		log.Fatal("Usage: v6d-fuse <mountpoint>")
	}

	mountPoint := pflag.Args()[0]

	mockBackend := backend.NewMockBackendWithData(map[string]string{
		"/dir/file": "Hello, World!",
	})

	cache, err := internal.NewV6dCache(args.V6dSocket)
	if err != nil {
		log.Fatalf("Failed to create cache: %v", err)
	}

	// root node of fs
	root := v6dfs.NewV6dRoot(mockBackend, cache)
	server, err := fs.Mount(mountPoint, root, &fs.Options{
		MountOptions: fuse.MountOptions{
			Debug:              args.Debug,
			DisableReadDirPlus: true,
		},
	})
	if err != nil {
		log.Panicf("Mount fail: %v\n", err)
	}
	log.Printf("Mounted on %s\n", mountPoint)
	log.Printf("Unmount by calling 'fusermount -u %s'", mountPoint)

	server.Wait()
}
