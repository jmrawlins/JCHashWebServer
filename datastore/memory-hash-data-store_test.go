package datastore

import (
	"sort"
	"sync"
	"testing"
)

func TestGetNextId(t *testing.T) {
	t.Run("Is thread safe", func(t *testing.T) {
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
}
