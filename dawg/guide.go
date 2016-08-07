package dawg

// status ok
import (
	"encoding/binary"
	"io"
)

// Guide ...
type Guide struct {
	units []uint8
}

func (g *Guide) child(index uint32) uint8 {
	return g.units[index*2]
}

func (g *Guide) sibling(index uint32) uint8 {
	return g.units[index*2+1]
}

func (g *Guide) read(buf io.Reader) error {
	var baseSize uint32
	err := binary.Read(buf, binary.LittleEndian, &baseSize)
	if err != nil {
		return err
	}

	g.units = make([]uint8, baseSize*2)
	return binary.Read(buf, binary.LittleEndian, &g.units)
}

func (g *Guide) size() int {
	return len(g.units)
}

// NewGuide - constructor for Guide
func NewGuide() *Guide {
	return &Guide{}
}
