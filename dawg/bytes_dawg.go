package dawg

// status ok
import (
	"bytes"
	"encoding/base64"

	"github.com/ReanGD/go-morphy/std"
)

var (
	constPayloadSeparator = []byte("\x01")
)

// BytesDAWG - DAWG that is able to transparently store extra binary payload in keys;
// there may be several payloads for the same key.
// In other words, this class implements read-only DAWG-based
// {unicode -> list of bytes objects} mapping.
type BytesDAWG struct {
	CompletionDAWG
}

// Contains ...
func (d *BytesDAWG) Contains(key string) bool {
	_, ok := d.followKey([]byte(key))
	return ok
}

// Get - Returns a list of payloads (as byte objects) for a given key
func (d *BytesDAWG) Get(key string) ([][]byte, bool) {
	index, ok := d.followKey([]byte(key))
	if !ok {
		return [][]byte{}, false
	}
	res := d.valueForIndex(index)

	return res, len(res) != 0
}

func (d *BytesDAWG) followKey(key []byte) (uint32, bool) {
	index, ok := d.dct.followBytes(key, constRoot)
	if !ok {
		return index, ok
	}
	return d.dct.followChar(constPayloadSeparatorUint, index)
}

func (d *BytesDAWG) decode(src []byte) []byte {
	dstLen := base64.StdEncoding.DecodedLen(len(src))
	dst := make([]byte, dstLen)
	n, err := base64.StdEncoding.Decode(dst, src)
	if err != nil {
		panic(err)
	}
	return dst[:n]
}

func (d *BytesDAWG) valueForIndex(index uint32) [][]byte {
	res := [][]byte{}
	completer := NewCompleter(d.dct, d.guide)
	completer.start(index, []byte{})
	for completer.next() {
		res = append(res, d.decode(completer.key))
	}

	return res
}

func (d *BytesDAWG) stringBySeparator(src []byte) string {
	pos := bytes.IndexByte(src, constPayloadSeparator[0])
	if pos == -1 {
		panic("Separator is not in array")
	}
	return string(src[:pos])
}

// Keys ...
func (d *BytesDAWG) Keys(prefix string) []string {
	bPrefix := []byte(prefix)
	res := []string{}

	index, ok := d.dct.followBytes(bPrefix, constRoot)
	if !ok {
		return res
	}

	completer := NewCompleter(d.dct, d.guide)
	completer.start(index, bPrefix)
	for completer.next() {
		res = append(res, d.stringBySeparator(completer.key))
	}

	return res
}

// Items ...
func (d *BytesDAWG) Items(prefix string) []std.StrBytes {
	res := []std.StrBytes{}
	bPrefix := []byte(prefix)
	index := constRoot
	var ok bool
	if len(prefix) != 0 {
		index, ok = d.dct.followBytes(bPrefix, index)
		if !ok {
			return res
		}
	}

	completer := NewCompleter(d.dct, d.guide)
	completer.start(index, bPrefix)

	for completer.next() {
		parts := bytes.Split(completer.key, constPayloadSeparator)
		res = append(res, std.StrBytes{Key: string(parts[0]), Value: d.decode(parts[1])})
	}

	return res
}

func (d *BytesDAWG) hasValue(index uint32) bool {
	_, ok := d.dct.followChar(constPayloadSeparatorUint, index)
	return ok
}

func (d *BytesDAWG) initBytesDAWG() {
	d.initCompletionDAWG()
}

// NewBytesDAWG - constructor for BytesDAWG
func NewBytesDAWG() *BytesDAWG {
	dawg := &BytesDAWG{}
	dawg.vDAWG = dawg
	dawg.initBytesDAWG()

	return dawg
}
