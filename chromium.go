// Package simplecache provides support for reading Chromium simple cache.
// http://www.chromium.org/developers/design-documents/network-stack/disk-cache/very-simple-backend
package simplecache

import (
	"encoding/binary"
	"log"
)

const (
	indexMagicNumber uint64 = 0x656e74657220796f
	indexVersion     uint32 = 9

	indexHeaderSize int64 = 40
	indexEntrySize  int64 = 24
)

const (
	initialMagicNumber uint64 = 0xfcfb6d1ba7725c30
	finalMagicNumber   uint64 = 0xf4fa6f45970d41d8
	entryVersion       uint32 = 5

	entryHeaderSize int64 = 24
	entryEOFSize    int64 = 24

	flagCRC32  uint32 = 1
	flagSHA256 uint32 = 2 // (1U << 1)
)

const (
	sparseMagicNumber     uint64 = 0xeb97bf016553676b
	sparseRangeHeaderSize int64  = 28
)

// fakeIndex is the content of the index file.
type fakeIndex struct {
	Magic   uint64
	Version uint32
	_       uint64
}

// indexHeader is the header of the the-real-index file.
type indexHeader struct {
	Payload    uint32
	CRC        uint32
	Magic      uint64
	Version    uint32
	EntryCount uint64
	CacheSize  uint64
	Reason     uint32
}

// indexEntry is an entry in the the-real-index file.
type indexEntry struct {
	Hash     uint64
	LastUsed int64
	Size     uint64
}

// entryHeader is the header of an entry file.
type entryHeader struct {
	Magic   uint64
	Version uint32
	KeyLen  int32
	KeyHash uint32
	Pad     int32
}

// entryEOF ends a stream in an entry file.
type entryEOF struct {
	Magic      uint64
	Flag       uint32
	CRC        uint32
	StreamSize int32
	Pad     int32
}

// HasCRC32
func (e entryEOF) HasCRC32() bool {
	return e.Flag&flagCRC32 != 0
}

// HasSHA256
func (e entryEOF) HasSHA256() bool {
	return e.Flag&flagSHA256 != 0
}

// sparseRangeHeader is the header of a stream range in a sparse file.
type sparseRangeHeader struct {
	Magic  uint64
	Offset int64
	Len    int64
	CRC    uint32
}

// sparseRange is a stream range in a sparse file.
type sparseRange struct {
	Offset     int64
	Len        int64
	CRC        uint32
	FileOffset int64
}

type sparseRanges []sparseRange

func (ranges sparseRanges) Len() int {
	return len(ranges)
}
func (ranges sparseRanges) Swap(i, j int) {
	ranges[i], ranges[j] = ranges[j], ranges[i]
}
func (ranges sparseRanges) Less(i, j int) bool {
	var rng0, rng1 = ranges[i], ranges[j]
	return rng0.Offset < rng1.Offset
}

func init() {
	var index indexHeader
	if n := binary.Size(index); int64(n) != indexHeaderSize {
		log.Fatalf("IndexHeader size error: %d, want: %d", n, indexHeaderSize)
	}

	var entry indexEntry
	if n := binary.Size(entry); int64(n) != indexEntrySize {
		log.Fatalf("IndexEntry size error: %d, want: %d", n, indexEntrySize)
	}

	var entryHead entryHeader
	if n := binary.Size(entryHead); int64(n) != entryHeaderSize {
		log.Fatalf("EntryHeader size error: %d, want: %d", n, entryHeaderSize)
	}

	var entryEnd entryEOF
	if n := binary.Size(entryEnd); int64(n) != entryEOFSize {
		log.Fatalf("EntryEOF size error: %d, want: %d", n, entryEOFSize)
	}

	var rangeHeader sparseRangeHeader
	if n := binary.Size(rangeHeader); int64(n) != sparseRangeHeaderSize {
		log.Fatalf("SparseHeader size error: %d, want: %d", n, sparseRangeHeaderSize)
	}
}
