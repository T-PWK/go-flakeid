# Flake - decentralized, k-ordered identifier generation in Go

[![Build Status](https://travis-ci.com/T-PWK/go-flakeid.svg?branch=master)](https://travis-ci.com/T-PWK/go-flakeid)
[![GitHub issues](https://img.shields.io/github/issues/T-PWK/go-flakeid.svg)](https://github.com/T-PWK/go-flakeid/issues)
[![Go Report Card](https://goreportcard.com/badge/github.com/T-PWK/go-flakeid)](https://goreportcard.com/report/github.com/T-PWK/go-flakeid)
[![Coverage Status](https://coveralls.io/repos/github/T-PWK/go-flakeid/badge.svg?branch=master)](https://coveralls.io/github/T-PWK/go-flakeid?branch=master)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](http://blog.abelotech.com/mit-license)

Flake ID generator produces 64-bit, k-ordered, conflict-free ids in a distributed environment.

To install

```
go get -u github.com/t-pwk/go-flakeid
```

## Flake Numbers Format

The Flake ID is made up of: `timestamp`, `datacenter`, `worker` and `counter`. Examples in the following table:

```
+-------------+------------+--------+---------+--------------------+
|  Timestamp  | Datacenter | Worker | Counter | Flake ID           |
+-------------+------------+--------+---------+--------------------+
| 0x8c20543b0 |   00000b   | 00000b |  0x000  | 0x02308150ec000000 |
+-------------+------------+--------+---------+--------------------+
| 0x8c20543b1 |   00000b   | 00000b |  0x000  | 0x02308150ec400000 |
+-------------+------------+--------+---------+--------------------+
| 0x8c20543b1 |   00000b   | 00000b |  0x001  | 0x02308150ec400001 |
+-------------+------------+--------+---------+--------------------+
| 0x8c20543b1 |   00000b   | 00000b |  0x002  | 0x02308150ec400002 |
+-------------+------------+--------+---------+--------------------+
| 0x8c20543b1 |   00000b   | 00000b |  0x003  | 0x02308150ec400003 |
+-------------+------------+--------+---------+--------------------+
| 0x8c20c0335 |   00011b   | 00001b |  0x000  | 0x02308300cd461000 |
+-------------+------------+--------+---------+--------------------+
| 0x8c20c0335 |   00011b   | 00001b |  0x001  | 0x02308300cd461001 |
+-------------+------------+--------+---------+--------------------+
```

- `timestamp`, a 42 bit long number of milliseconds elapsed since 1 January 1970 00:00:00 UTC
- `datacenter`, a 5 bit long datacenter identifier. It can take up to 32 unique values (including 0)
- `worker`, a 5 bit long worker identifier. It can take up to 32 unique values (including 0)
- `counter`, a 12 bit long counter of ids in the same millisecond. It can take up to 4096 unique values.

Example of a breakdown of bits for an identifier `5828128208445124609` (counter is `1`, datacenter is `7` and worker `3`) is as follows:

```
 010100001110000110101011101110100001000111 00111 00011 000000000001
                                                       |------------| 12 bit counter
                                                 |-----|               5 bit worker
                                           |-----|                     5 bit datacenter
                                           |----- -----|              10 bit generator identifier
|------------------------------------------|                          42 bit timestamp
```

Note that composition of `datacenter id` and `worker id` makes 1024 unique generator identifiers. By modifying datacenter and worker id we can get up to 1024 id generators on a single machine (e.g. each running in a separate process) or have 1024 machines with a single id generator on each.

## Usage

Flake ID Generator returns 64-bit long unsigned integer. Every time you call `FlakeID.NextID` you get a k-ordered, conflict-free identifier.

```go
package  main

import (
  "fmt"
  "github.com/t-pwk/go-flakeid"
)

func main() {
  g := flakeid.FlakeID{}

  fmt.Printf("%x\n", g.NextID()) // 1530d02cb005000
  fmt.Printf("%x\n", g.NextID()) // 1530d02cb005001

  fmt.Printf("%b\n", g.NextID()) // 101010011000011010000001011001011000000000101000000000010
  fmt.Printf("%b\n", g.NextID()) // 101010011000011010000001011001011000000000101000000000011
}
```

If you want to assign different worker or datacenter identifiers, you can do it during a generator creation or after.

```go
package  main

import (
  "fmt"
  "github.com/t-pwk/go-flakeid"
)

func main() {
  g1 := flakeid.FlakeID{WorkerID: 7, DatacenterID: 7}

  g2 := flakeid.FlakeID{}
  g2.WorkerID = 7
  g2.DatacenterID = 7
}
```

If your generator works in a single environment and you would like to use 1024 unique workers, you can convert the worker identifier into the generator's datacenter and worker id in the following way:

```go
package main

import (
  "fmt"
  "github.com/t-pwk/go-flakeid"
)

const (
  mask = 0x1F
)

func main() {
  var worker uint64 = 1022

  g := flakeid.FlakeID{WorkerID: worker & mask, DatacenterID: (worker >> 5) & mask}
  fmt.Printf("W: %b, D: %b\n", g.WorkerID, g.DatacenterID) // W: 11110, D: 11111
}
```

You can also slightly reduce range of the generated identifiers by providing the `Epoc` parameter value. That value is to reduce timestamp (number of milliseconds elapsed since 1 January 1970 00:00:00 UTC) value when building identifiers.

```go
package main

import (
  "fmt"
  "time"
  "github.com/t-pwk/go-flakeid"
)

func main() {
  g1 := flakeid.FlakeID{}

  g2 := flakeid.FlakeID{Epoc: flakeid.Epoc1Jan2000}

  epoc := time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC).UnixNano() / int64(time.Millisecond)
  g3 := flakeid.FlakeID{Epoc: uint64(epoc)}

  fmt.Printf("%x\n", g1.NextID()) // 597edbca22c00000
  fmt.Printf("%b\n", g1.NextID()) // 101100101111110110110111100101000100010110000000000000000000001

  fmt.Printf("%x\n", g2.NextID()) // 226427df22c00000
  fmt.Printf("%b\n", g2.NextID()) // 10001001100100001001111101111100100010110000000000000000000001

  fmt.Printf("%x\n", g3.NextID()) // 1531aa622c00000
  fmt.Printf("%b\n", g3.NextID()) // 101010011000110101010011000100010110000000000000000000001
}
```

As you can see, the range values varied depending on the value added to the `Epoc` parameter. Please note that the `epoc` parameter must be the same for all identifiers. Otherwise, a generator will not be able to generate k-ordered, conflict-free ids. Hence, you should never use current time or a value that changes from execution to execution for an epoch. Always use some constants, or do not use that feature at all.

### Formatting

FlakeID generator returns uint64 number. You can use different formats, using for example `fmt` package or convert an identifier to Base64 format.

```js
package  main

import (
  "encoding/base64"
  "encoding/binary"
  "fmt"
  "github.com/t-pwk/go-flakeid"
)

func main() {
  g := flakeid.FlakeID{WorkerID: 7, DatacenterID: 7}

  fmt.Printf("%d\n", g.NextID())   // 6448828961128345600
  fmt.Printf("%x\n", g.NextID())   // 597ed7c5d54e7001
  fmt.Printf("0x%x\n", g.NextID()) // 0x597ed7c5d54e7002
  fmt.Printf("%X\n", g.NextID())   // 597ED7C5D54E7003
  fmt.Printf("%b\n", g.NextID())   // 101100101111110110101111100010111010101010011100111000000000100

  b := make([]byte, 8)
  binary.LittleEndian.PutUint64(b, g.NextID())
  fmt.Println(base64.StdEncoding.EncodeToString(b)) // BXDOxnfZflk=
}
```

## Documentation

[Documentation](https://godoc.org/github.com/T-PWK/go-flakeid) is hosted at GoDoc project.

## Author

Written by Tom Pawlak - [Blog](https://blog.abelotech.com)

## License

Copyright (c) 2018 Tom Pawlak

MIT License : https://blog.abelotech.com/mit-license/
