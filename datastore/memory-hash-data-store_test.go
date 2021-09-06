package datastore

import (
	"fmt"
	"sort"
	"sync"
	"testing"
)

func TestMemoryHashDataStore(t *testing.T) {
	t.Parallel()

	t.Run("GetNextId", func(t *testing.T) {
		t.Parallel()
		t.Run("Is thread safe", func(t *testing.T) {
			t.Parallel()
			ds := NewMemoryHashDataStore()

			const goRoutines = 10
			const iterations = 1000
			ids := [goRoutines][iterations]uint64{}

			// Verify that GetNextId never assigns the same id
			wg := sync.WaitGroup{}
			for i := range [goRoutines]int{} {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()

					for j := range [iterations]int{} {
						ids[i][j], _ = ds.GetNextId()
					}
				}(i)
			}
			wg.Wait()

			flatids := []uint64{}
			for i := range [goRoutines]int{} {
				flatids = append(flatids[:], ids[i][:]...)
			}
			sort.Slice(flatids, func(i, j int) bool { return flatids[i] < flatids[j] })
			for i := range [goRoutines*iterations - 1]int{} {
				if flatids[i] == flatids[i+1] {
					t.Failed()
				}
			}
		})

		t.Run("GetNextId starts at 1 and goes up 1 at a time.", func(t *testing.T) {
			t.Parallel()
			ds := NewMemoryHashDataStore()
			if id, err := ds.GetNextId(); id != 1 || err != nil {
				t.Fail()
			}
			if id, err := ds.GetNextId(); id != 2 || err != nil {
				t.Fail()
			}
		})
	})

	t.Run("SetLastId", func(t *testing.T) {
		t.Parallel()

		t.Run("SetLastId sets the last id", func(t *testing.T) {
			t.Parallel()
			ds := NewMemoryHashDataStore()
			if err := ds.SetLastId(42); err != nil {
				t.Fail()
			}
			if id, err := ds.GetNextId(); id != 43 || err != nil {
				t.Fail()
			}
		})

		t.Run("SetLastId returns error after GetNextId is called", func(t *testing.T) {
			t.Parallel()
			ds := NewMemoryHashDataStore()
			ds.GetNextId()
			if err := ds.SetLastId(42); nil == err {
				t.Fail()
			}
		})
	})

	t.Run("StoreHash", func(t *testing.T) {
		t.Parallel()

		t.Run("StoreHash stores the requested hash if the id is available, errors if it isn't", func(t *testing.T) {
			t.Parallel()
			ds := NewMemoryHashDataStore()
			id, _ := ds.GetNextId()

			if err := ds.StoreHash(id, "SomeValue"); err != nil {
				t.Fail()
			}

			if vals := ds.GetAllHashes(); len(vals) != 1 || (vals)[id] != "SomeValue" {
				t.Fail()
			}

			if err := ds.StoreHash(id, "SomeOtherValue"); err == nil {
				t.Fail()
			}

			if vals := ds.GetAllHashes(); len(vals) != 1 || (vals)[id] != "SomeValue" {
				t.Fail()
			}
		})
	})

	t.Run("GetHash", func(t *testing.T) {
		t.Parallel()

		t.Run("Gets the requested hash if it exists", func(t *testing.T) {
			t.Parallel()

			ds := NewMemoryHashDataStore()
			id, _ := ds.GetNextId()

			if err := ds.StoreHash(id, "SomeValue"); err != nil {
				t.Fail()
			}

			if val, err := ds.GetHash(id); val != "SomeValue" || err != nil {
				t.Fail()
			}
		})

		t.Run("returns error if the requested id hash isn't set", func(t *testing.T) {
			t.Parallel()

			ds := NewMemoryHashDataStore()

			id, _ := ds.GetNextId()
			if err := ds.StoreHash(id, "SomeValue"); err != nil {
				t.Fail()
			}

			if _, err := ds.GetHash(42); err == nil {
				t.Fail()
			}
		})
	})

	t.Run("GetAllHashes", func(t *testing.T) {
		t.Parallel()

		t.Run("Works when no hashes are set yet", func(t *testing.T) {
			t.Parallel()

			ds := NewMemoryHashDataStore()

			if vals := ds.GetAllHashes(); len(vals) != 0 {
				t.Fail()
			}
		})

		t.Run("Returns correct values after several sets", func(t *testing.T) {
			t.Parallel()

			ds := NewMemoryHashDataStore()

			for range [10]int{} {
				id, _ := ds.GetNextId()
				if err := ds.StoreHash(id, fmt.Sprintf("SomeValue%d", id)); err != nil {
					t.Fail()
				}
			}

			vals := ds.GetAllHashes()

			if len(vals) != 10 {
				t.Fail()
			}

			for i := range [10]uint64{} {
				id := uint64(i + 1)
				if val, ok := vals[id]; val != fmt.Sprintf("SomeValue%d", id) || !ok {
					t.Fail()
				}
			}
		})
	})

	t.Run("StoreRequestTime", func(t *testing.T) {
		t.Parallel()

		t.Run("calculates request times correctly", func(t *testing.T) {
			t.Parallel()

			uri := "/test"

			ds := NewMemoryHashDataStore()

			// Sum(0..9) = 45, so Average = 4.5
			for i := range [10]int{} {
				ds.StoreRequestTime(uri, int64(i))
			}

			if stat, err := ds.GetUriStats(uri); err != nil || stat.Average != 4.5 || stat.URI != uri || stat.Total != 10 {
				t.Fail()
			}

			// Sum(0..9) = 45, so Average = 4.5
			for range [10]int{} {
				ds.StoreRequestTime("/other", int64(10))
			}

			if stat, err := ds.GetUriStats("/other"); err != nil || stat.Average != 10 || stat.URI != "/other" || stat.Total != 10 {
				t.Fail()
			}
		})
	})

	t.Run("GetStats", func(t *testing.T) {
		t.Parallel()

	})

	t.Run("GetUriStats", func(t *testing.T) {
		t.Parallel()

		t.Run("Returns valid response when the requested uri hasn't been accessed yet", func(t *testing.T) {
			t.Parallel()
			ds := NewMemoryHashDataStore()

			if stat, err := ds.GetUriStats("foo"); err != nil || stat.Average != 0 || stat.URI != "foo" || stat.Total != 0 {
				t.Fail()
			}
		})
	})

}
