package dawg

// status ok
import (
	"encoding/binary"
	"io"
)

// Dictionary - class for retrieval and binary I/O.
type Dictionary struct {
	units []uint32
}

// hasValue - Checks if a given index is related to the end of a key.
func (d *Dictionary) hasValue(index uint32) bool {
	return hasLeaf(d.units[index])
}

// value - Gets a value from a given index.
func (d *Dictionary) value(index uint32) uint32 {
	offset := getOffset(d.units[index])
	valueIndex := (index ^ offset) & constPrecisionMask
	return getValue(d.units[valueIndex])
}

// read - Reads a dictionary from an input stream.
func (d *Dictionary) Read(buf io.Reader) error {
	var baseSize uint32
	err := binary.Read(buf, binary.LittleEndian, &baseSize)
	if err != nil {
		return err
	}

	d.units = make([]uint32, baseSize)
	return binary.Read(buf, binary.LittleEndian, &d.units)
}

// Contains - Exact matching.
func (d *Dictionary) Contains(key []byte) bool {
	index, ok := d.followBytes(key, constRoot)
	if !ok {
		return false
	}
	return d.hasValue(index)
}

// Find - Exact matching (returns value)
func (d *Dictionary) Find(key []byte) (uint32, bool) {
	index, ok := d.followBytes(key, constRoot)
	if !ok || !d.hasValue(index) {
		return 0, false
	}

	return d.value(index), true
}

// followChar - Follows a transition.
func (d *Dictionary) followChar(label uint32, index uint32) (uint32, bool) {
	offset := getOffset(d.units[index])
	nextIndex := (index ^ offset ^ label) & constPrecisionMask
	return nextIndex, getLabel(d.units[nextIndex]) == label
}

// followBytes - Follows transitions.
func (d *Dictionary) followBytes(s []byte, index uint32) (uint32, bool) {
	var ok bool
	for _, ch := range s {
		index, ok = d.followChar(uint32(ch), index)
		if !ok {
			return index, ok
		}
	}

	return index, true
}

// NewDictionary - constructor for Dictionary
func NewDictionary() *Dictionary {
	return &Dictionary{}
}
