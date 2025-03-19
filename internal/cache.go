package internal

import (
	"log"

	v6d "github.com/v6d-io/v6d/go/vineyard/pkg/client"
)

// Cache is a abstraction for v6d cluster
type Cache interface {
	Put(value []byte) (id uint64, err error)
	Get(id uint64) (value []byte, err error)
}

type V6dCache struct {
	client *v6d.IPCClient
}

func NewV6dCache(v6dSocket string) (*V6dCache, error) {
	client, err := v6d.NewIPCClient(v6dSocket)
	if err != nil {
		return nil, err
	}
	return &V6dCache{client: client}, nil
}

func (c *V6dCache) Put(value []byte) (id uint64, err error) {
	log.Printf("V6dCache Put: %v, len: %d", string(value), len(value))
	id, err = c.client.BuildBuffer(value, uint64(len(value)))
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c *V6dCache) Get(id uint64) (value []byte, err error) {
	log.Printf("V6dCache Get: %d", id)
	blob, err := c.client.GetBuffer(id, false)
	if err != nil {
		return nil, err
	}
	return blob.Buffer.Buf(), nil
}

// client, err := v6d.NewIPCClient(*v6dSocket)
// if err != nil {
// 	log.Fatalf("Failed to create client: %v", err)
// }
// defer client.Disconnect()
// buffer := []byte("hello")
// id, err := client.BuildBuffer(buffer, uint64(len(buffer)))
// if err != nil {
// 	log.Fatalf("Failed to create blob: %v", err)
// }

// fmt.Printf("Created data with id: %v\n", id)

// blob, err := client.GetBuffer(id, false)
// if err != nil {
// 	log.Fatalf("Failed to get buffer: %v", err)
// }

// fmt.Printf("Buffer: %v\n", blob.Buffer)
