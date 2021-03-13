package flakeid_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/t-pwk/go-flakeid"
)

const count = 10_000

func BenchmarkNextID(b *testing.B) {
	gen := flakeid.FlakeID{}
	for i := 0; i < b.N; i++ {
		gen.NextID()
	}
}

func TestNoDuplicates(t *testing.T) {
	var ids [count]uint64
	dup := make(map[uint64]int, count)
	gen := flakeid.FlakeID{}

	for i := 0; i < count; i++ {
		ids[i] = gen.NextID()
	}

	for _, id := range ids {
		dup[id]++
		if dup[id] > 1 {
			t.Errorf("Duplicate identifier found: %d", id)
		}
	}
}

func TestDatacenter(t *testing.T) {
	gen := flakeid.FlakeID{DatacenterID: 7}
	id := fmt.Sprintf("%b", gen.NextID())

	if id[41:51] != "0011100000" {
		t.Error("Invalid id bits")
	}
}

func TestWorker(t *testing.T) {
	gen := flakeid.FlakeID{WorkerID: 7}
	id := fmt.Sprintf("%b", gen.NextID())

	if id[41:51] != "0000000111" {
		t.Error("Invalid id bits")
	}
}

func Example() {
	gen := flakeid.FlakeID{}

	fmt.Printf("%x", gen.NextID()) // output like 597ed3f490000000
	fmt.Printf("%b", gen.NextID()) // output like 101100101111110110100111111010010010000000000000000000000000001
}

func Example_datacenter() {
	gen := flakeid.FlakeID{DatacenterID: 5}

	fmt.Printf("%x", gen.NextID()) // output like 597ed40a654a0000
	fmt.Printf("%b", gen.NextID()) // output like 101100101111110110101000000101001100101010010100000000000000001
}

func Example_worker() {
	gen := flakeid.FlakeID{WorkerID: 7}

	fmt.Printf("%x", gen.NextID()) // output like 597ed425adc07000
	fmt.Printf("%b", gen.NextID()) // output like 101100101111110110101000010010110101101110000000111000000000001
}

func Example_epoc_1Jan2000() {
	gen := flakeid.FlakeID{Epoc: flakeid.Epoc1Jan2000}

	fmt.Printf("%x", gen.NextID()) // output like 2264205ec2c00000
	fmt.Printf("%b", gen.NextID()) // output like 10001001100100001000000101111011000010110000000000000000000001
}

func Example_epoc_1Jan2018() {
	epoc := time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC).UnixNano() / int64(time.Millisecond)
	gen := flakeid.FlakeID{Epoc: uint64(epoc)}

	fmt.Printf("%x", gen.NextID()) // output like 15313f357000000
	fmt.Printf("%b", gen.NextID()) // output like 101010011000100111111001101010111000000000000000000000001
}
