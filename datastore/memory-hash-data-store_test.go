package datastore

import (
	"sort"
	"sync"
	"testing"
)

func TestGetNextId(t *testing.T) {
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
}

func TestSetLastId(t *testing.T) {
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
}

func TestStoreHash(t *testing.T) {
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
}

func TestGetHash(t *testing.T) {
	t.Parallel()

	t.Run("Gets the requested hash if it exists", func(t *testing.T) {
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
		ds := NewMemoryHashDataStore()

		id, _ := ds.GetNextId()
		if err := ds.StoreHash(id, "SomeValue"); err != nil {
			t.Fail()
		}

		if _, err := ds.GetHash(42); err == nil {
			t.Fail()
		}
	})

}

func TestGetAllHashes(t *testing.T) {
	t.Parallel()

}

func TestStoreRequestTime(t *testing.T) {
	t.Parallel()

}

func TestGetStats(t *testing.T) {
	t.Parallel()

}

func TestGetUriStats(t *testing.T) {
	t.Parallel()

}
