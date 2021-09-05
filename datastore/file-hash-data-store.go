package datastore

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
)

// FileHashDataStore implements HashDataStore with file-driven persistence.
//
// It uses MemoryHashDataStore for base functionality, and wraps it with
// reading and writing via the argument io.ReadWriteSeeker.
//
// The file is expected to be a tab-separated value file with two unnamed columns: id, and hash value.
// Invalid lines will be logged and skipped:
type FileHashDataStore struct {
	mds     *MemoryHashDataStore
	rws     io.ReadWriteSeeker
	rwsLock *sync.Mutex
}

func NewFileHashDataStore(rws io.ReadWriteSeeker, mds *MemoryHashDataStore) (*FileHashDataStore, error) {
	ds := FileHashDataStore{mds, rws, &sync.Mutex{}}

	// Read if there is content
	scanner := bufio.NewScanner(rws)
	lastId := uint64(0)
	for scanner.Scan() {
		line := scanner.Text()

		vals := strings.Split(line, "\t")
		if len(vals) != 2 {
			log.Printf("WARN: Invalid hashes data found: (%s). Skipping line\n", line)
		} else {
			id, err := strconv.ParseUint(vals[0], 10, 64)
			if err != nil {
				log.Printf("WARN: Invalid id found: (%s). Skipping line\n", line)
			}

			hash := vals[1]

			lastId = id
			mds.StoreHash(id, hash)
		}
	}

	// Finally, initialize the MemoryDataStore's next id assigned, since we loaded data from outside
	if err := ds.mds.SetLastId(lastId + 1); err != nil {
		return nil, fmt.Errorf("failed to set last id in MemoryDataStore: %s", err.Error())
	}

	// Future improvement: Limit the number of bad lines in a file before aborting.
	// Unnecessary if we're the only ones touching the file
	return &ds, nil
}

func (ds *FileHashDataStore) GetNextId() (uint64, error) {
	return ds.mds.GetNextId()
}

func (ds *FileHashDataStore) StoreHash(id uint64, hash string) error {
	line := fmt.Sprintf("%d\t%s\n", id, hash)

	if err := ds.storeHash(line); err != nil {
		return err
	}

	return ds.mds.StoreHash(id, hash)
}

func (ds *FileHashDataStore) storeHash(line string) error {
	ds.rwsLock.Lock()
	defer ds.rwsLock.Unlock()

	// Seek to end of file (in case this is our first time writing)
	if _, err := ds.rws.Seek(0, 2); err != nil {
		return fmt.Errorf("ERROR: Cannot seek end of hashes file: %s. (%s) will not be persisted.", err.Error(), line)
	}

	if _, err := fmt.Fprintf(ds.rws, line); err != nil {
		return fmt.Errorf("ERROR: Cannot write to hashes file: %s. (%s) will not be persisted", err.Error(), line)
	}
	return nil
}

func (ds *FileHashDataStore) GetHash(id uint64) (string, error) {
	return ds.mds.GetHash(id)
}

func (ds *FileHashDataStore) GetAllHashes() map[uint64]string {
	return ds.mds.GetAllHashes()
}
