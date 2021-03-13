// Package flakeid generator yielding k-ordered, conflict-free identifiers.
// Each identifier is a 64 bits unsigned integer, consisting of:
// timestamp, a 42 bit long number of milliseconds elapsed since 1 January 1970 00:00:00 UTC;
// datacenter, a 5 bit long identifier, which can take up to 32 unique values (including 0);
// worker, a 5 bit long worker indentifier, which can take up to 32 unique values (including 0);
// counter, a 12 bit long counter of identifiers in the same millisecond, which can take up to 4096 unique values (including 0).
//
// Breakdown of bits for an id e.g. 5828128208445124609 (counter is 1, datacenter is 7 and worker 3) is as follows:
//   010100001110000110101011101110100001000111 00111 00011 000000000001
//                                                         |------------| 12 bit counter
//                                                   |-----|               5 bit worker
//                                             |-----|                     5 bit datacenter
//  |------------------------------------------|                          42 bit timestamp
package flakeid

import (
	"fmt"
	"sync"
	"time"
)

// Epoc1Jan2000 is the time in milliseconds on 1st of January 2000 at midnight UTC
const Epoc1Jan2000 = 946684800000

const (
	workerIDBits     = 5  // Number of bits allocated for a worker id in the generated identifier. 5 bits indicates values from 0 to 31
	datacenterIDBits = 5  // Datacenter identifier this worker belongs to. 5 bits indicates values from 0 to 31
	sequenceBits     = 12 // Number of bits allocated for sequence in the generated identifier

	workerIDShift      = sequenceBits
	datacenterIDShift  = workerIDShift + workerIDBits
	timestampLeftShift = datacenterIDShift + datacenterIDBits

	workerMask     = ^(-1 << workerIDBits)     // Maximum worker identifier
	datacenterMask = ^(-1 << datacenterIDBits) // Maximum datacenter identifier
	sequenceMask   = ^(-1 << sequenceBits)
)

func currentTime() uint64 {
	return uint64(time.Now().UnixNano() / int64(time.Millisecond))
}

func tillNextMills(lastTimestamp uint64) uint64 {
	timestamp := currentTime()

	for lastTimestamp == timestamp {
		time.Sleep(time.Millisecond)
		timestamp = currentTime()
	}

	return timestamp
}

// FlakeID is an identifier generator
type FlakeID struct {
	WorkerID, DatacenterID, Epoc uint64
	lastTimestamp, sequence      uint64
	mutex                        sync.Mutex
}

func (f *FlakeID) nextID() uint64 {
	timestamp := currentTime()

	switch {
	case timestamp < f.lastTimestamp:
		panic(fmt.Sprintf("Clock moved backwards. Refusing to generate id for %d milliseconds", f.lastTimestamp-timestamp))
	case timestamp == f.lastTimestamp:
		f.sequence++
		if f.sequence > sequenceMask {
			f.sequence = 0
			timestamp = tillNextMills(timestamp)
		}
	default:
		f.sequence = 0
	}

	f.lastTimestamp = timestamp

	return (timestamp-f.Epoc)<<timestampLeftShift |
		(f.DatacenterID & datacenterMask << datacenterIDShift) |
		(f.WorkerID & workerMask << workerIDShift) |
		f.sequence
}

// NextID yields k-ordered, conflict-free identifier.
func (f *FlakeID) NextID() uint64 {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.nextID()
}
