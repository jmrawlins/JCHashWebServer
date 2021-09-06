package datastore

import (
	"bytes"
	"testing"
)

func TestNewFileHashDataStore(t *testing.T) {
	t.Run("Reads the input file", func(t *testing.T) {
		t.Parallel()

		buf := []byte("1	firsthash\n2	secondhash\n")
		r := bytes.NewReader(buf)
		w := bytes.NewBuffer(buf)

		mds := NewMemoryHashDataStore()
		ds, err := NewFileHashDataStore(r, w, mds)
		if err != nil {
			t.Fail()
		}

		if vals := ds.GetAllHashes(); vals[1] != "firsthash" || vals[2] != "secondhash" {
			t.Fail()
		}
	})
}

// TODO: Ideally we'd mock MemoryHashDataStore and test that we're calling its functions properly
// instead of using the real class, but for now we just verify functionality by using the class

func TestFileHashDataStore(t *testing.T) {
	t.Parallel()

	t.Run("StoreHash", func(t *testing.T) {
		t.Parallel()

		t.Run("StoreHash stores the hash to file and to its MemoryHashDataStore", func(t *testing.T) {
			t.Parallel()

			buf := []byte("")
			r := bytes.NewReader(buf)
			w := bytes.NewBuffer(buf)

			mds := NewMemoryHashDataStore()
			ds, err := NewFileHashDataStore(r, w, mds)
			if err != nil {
				t.Fail()
			}

			if err := ds.StoreHash(1, "firsthash"); err != nil {
				t.Fail()
			}

			if err := ds.StoreHash(2, "secondhash"); err != nil {
				t.Fail()
			}

			// Check that it went through mds okay
			if vals := ds.GetAllHashes(); len(vals) != 2 || vals[1] != "firsthash" || vals[2] != "secondhash" {
				t.Fail()
			}

			// Check that it wrote to io correctly
			written := string(w.Bytes())
			if written != "1\tfirsthash\n2\tsecondhash\n" {
				t.Fail()
			}
		})
	})

	t.Run("GetHash", func(t *testing.T) {
		t.Parallel()
		// Skipping because this function is just a call through to mds
	})

	t.Run("GetAllHashes", func(t *testing.T) {
		t.Parallel()
		// Skipping because this function is just a call through to mds
	})
}
