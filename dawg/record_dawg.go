package dawg

// status ok
import (
	"encoding/binary"
	"unicode/utf8"
)

// StringUints16 ...
type StringUints16 struct {
	Key   string
	Value []uint16
}

// StringUints16Arr ...
type StringUints16Arr struct {
	Key   string
	Value [][]uint16
}

// RecordDAWG ...
type RecordDAWG struct {
	*BytesDAWG
	// "HH" == 2, "3H" == 3
	fmt uint8
	// ">" - binary.BigEndian, "<" - binary.LittleEndian (default)
	order binary.ByteOrder
}

// Get - Returns a list of payloads (as uint16 objects) for a given key
func (d *RecordDAWG) Get(key string) ([][]uint16, bool) {
	index, ok := d.followKey([]byte(key))
	if !ok {
		return [][]uint16{}, false
	}
	res := d.valueForIndex(index)

	return res, len(res) != 0
}

func (d *RecordDAWG) bytesToUints16(src []byte) []uint16 {
	if len(src) != int(d.fmt)*2 {
		panic("source len error")
	}
	res := make([]uint16, d.fmt)
	for i := range res {
		res[i] = d.order.Uint16(src[2*i:])
	}
	return res
}

func (d *RecordDAWG) valueForIndex(index uint32) [][]uint16 {
	value := d.BytesDAWG.valueForIndex(index)
	res := make([][]uint16, len(value))
	for i, val := range value {
		res[i] = d.bytesToUints16(val)
	}

	return res
}

// Items ...
func (d *RecordDAWG) Items(prefix string) []StringUints16 {
	items := d.BytesDAWG.Items(prefix)
	res := make([]StringUints16, len(items))
	for i, item := range items {
		res[i].Key = item.Key
		res[i].Value = d.bytesToUints16(item.Value)
	}

	return res
}

func (d *RecordDAWG) similarItems(currentPrefix string, key string, index uint32, replaceChars map[rune]rune) []StringUints16Arr {
	res := []StringUints16Arr{}
	exitByBreak := false
	startPos := len(currentPrefix)

	for curPos, bStep := range key[startPos:] {
		ReplaceChar, ok := replaceChars[bStep]

		if ok {
			nextIndex, ok := d.dct.followBytes([]byte(string(ReplaceChar)), index)
			if ok {
				prefix := currentPrefix + key[startPos:curPos] + string(ReplaceChar)
				extraItems := d.similarItems(prefix, key, nextIndex, replaceChars)
				res = append(res, extraItems...)
			}
		}

		index, ok = d.dct.followBytes([]byte(string(bStep)), index)
		if !ok {
			exitByBreak = true
			break
		}
	}

	if !exitByBreak {
		index, ok := d.dct.followChar(constPayloadSeparatorUint, index)
		if ok {
			foundKey := currentPrefix + key[startPos:]
			value := d.valueForIndex(index)
			item := StringUints16Arr{foundKey, value}
			res = append([]StringUints16Arr{item}, res...)
		}
	}

	return res
}

// SimilarItems -
// Returns a list of (key, value) tuples for all variants of 'key'
// in this DAWG according to 'replaces'.
func (d *RecordDAWG) SimilarItems(key string, replaceChars map[rune]rune) []StringUints16Arr {
	return d.similarItems("", key, constRoot, replaceChars)
}

func (d *RecordDAWG) similarItemsValues(startPos int, key string, index uint32, replaceChars map[rune]rune) [][][]uint16 {
	res := [][][]uint16{}
	exitByBreak := false

	for curPos, bStep := range key[startPos:] {
		ReplaceChar, ok := replaceChars[bStep]

		if ok {
			nextIndex, ok := d.dct.followBytes([]byte(string(ReplaceChar)), index)
			if ok {
				extraItems := d.similarItemsValues(curPos+utf8.RuneLen(bStep), key, nextIndex, replaceChars)
				res = append(res, extraItems...)
			}
		}

		index, ok = d.dct.followBytes([]byte(string(bStep)), index)
		if !ok {
			exitByBreak = true
			break
		}
	}

	if !exitByBreak {
		index, ok := d.dct.followChar(constPayloadSeparatorUint, index)
		if ok {
			value := d.valueForIndex(index)
			res = append([][][]uint16{value}, res...)
		}
	}

	return res
}

// SimilarItemsValues -
// Returns a list of values tuples for all variants of 'key'
// in this DAWG according to 'replaces'.
func (d *RecordDAWG) SimilarItemsValues(key string, replaceChars map[rune]rune) [][][]uint16 {
	return d.similarItemsValues(0, key, constRoot, replaceChars)
}

// NewRecordDAWG - constructor for RecordDAWG
func NewRecordDAWG(fmt uint8, order binary.ByteOrder) *RecordDAWG {
	dawg := &RecordDAWG{BytesDAWG: NewBytesDAWG(), fmt: fmt, order: order}
	dawg.vDAWG = dawg

	return dawg
}
