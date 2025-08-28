package matrix

import (
	"encoding/binary"
	"golang.org/x/exp/mmap"
)

type M struct {
	n   int
	rdr *mmap.ReaderAt
	off int64
}

func Open(path string) (*M, error) {
	r, err := mmap.Open(path)
	if err != nil {
		return nil, err
	}
	var hdr [4]byte
	if _, err := r.ReadAt(hdr[:], 0); err != nil {
		return nil, err
	}
	n := int(binary.LittleEndian.Uint32(hdr[:]))
	return &M{n: n, rdr: r, off: 4}, nil
}

func (m *M) Get(i, j int) uint16 {
	// row offset = off + i*(2*n)
	off := m.off + int64(i*2*m.n+2*j)
	var b [2]byte
	_, _ = m.rdr.ReadAt(b[:], off)
	return binary.LittleEndian.Uint16(b[:])
}

func (m *M) N() int { return m.n }
