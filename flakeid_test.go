package flakeid

import (
	"sort"
	"testing"
)

const size = 100000

type uint64Sort []uint64

func (p uint64Sort) Len() int           { return len(p) }
func (p uint64Sort) Less(i, j int) bool { return p[i] < p[j] }
func (p uint64Sort) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func TestSingleRoutine(t *testing.T) {
	gen := new(FlakeID)
	ids := make([]uint64, size)

	for i := 0; i < size; i++ {
		ids[i] = gen.NextID()
	}

	verify(t, ids)
}

func TestMultipleRoutines(t *testing.T) {
	gen := new(FlakeID)
	ids := make([]uint64, size)
	done := make(chan bool)

	routine := func(ids []uint64) {
		for i := range ids {
			ids[i] = gen.NextID()
		}

		verify(t, ids)

		done <- true
	}

	go routine(ids[0:25000])
	go routine(ids[25000:50000])
	go routine(ids[50000:75000])
	go routine(ids[75000:])

	for i := 0; i < 4; i++ {
		<-done
	}

	sort.Sort(uint64Sort(ids))

	verify(t, ids)
}

func verify(t *testing.T, ids []uint64) {
	size := len(ids)
	for i := 1; i < size; i++ {
		if ids[i] <= ids[i-1] {
			t.Errorf("Invalid order of identifiers in iteration %d", i)
		}
	}
}

func BenchmarkNextID(b *testing.B) {
	gen := new(FlakeID)

	for n := 0; n < b.N; n++ {
		gen.NextID()
	}
}
