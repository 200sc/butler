package blockpool

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-errors/errors"
	"github.com/itchio/wharf/tlc"
	"golang.org/x/crypto/sha3"
)

// DiskSink stores blocks on disk by their hash and length. It's hard-coded to
// use shake128-32 as a hashing algorithm.
// If `BlockHashes` is set, will store block hashes there.
type DiskSink struct {
	BasePath string

	Container      *tlc.Container
	BlockAddresses BlockAddressMap
	BlockHashes    BlockHashMap

	hashBuf []byte
	shake   sha3.ShakeHash
}

var _ Sink = (*DiskSink)(nil)

// Store should not be called concurrently, as it will result in corrupted hashes
func (ds *DiskSink) Store(loc BlockLocation, data []byte) error {
	if ds.hashBuf == nil {
		ds.hashBuf = make([]byte, 32)
	}

	if ds.shake == nil {
		ds.shake = sha3.NewShake128()
	}

	ds.shake.Reset()
	_, err := ds.shake.Write(data)
	if err != nil {
		return errors.Wrap(err, 1)
	}

	_, err = io.ReadFull(ds.shake, ds.hashBuf)
	if err != nil {
		return errors.Wrap(err, 1)
	}

	if ds.BlockHashes != nil {
		ds.BlockHashes.Set(loc, append([]byte{}, ds.hashBuf...))
	}

	addr := fmt.Sprintf("shake128-32/%x/%d", ds.hashBuf, len(data))
	if ds.BlockAddresses != nil {
		ds.BlockAddresses.Set(loc, addr)
	}

	path := filepath.Join(ds.BasePath, addr)

	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return errors.Wrap(err, 1)
	}

	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		return errors.Wrap(err, 1)
	}

	return nil
}

func (ds *DiskSink) GetContainer() *tlc.Container {
	return ds.Container
}
