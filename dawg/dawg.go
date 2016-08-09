package dawg

// status ok
import "os"

type virtDAWG interface {
	hasValue(index uint32) bool
}

// DAWG - Base DAWG wrapper.
type DAWG struct {
	vDAWG virtDAWG
	dct   *Dictionary
}

// Contains - Exact matching.
func (d *DAWG) Contains(key string) bool {
	return d.dct.Contains([]byte(key))
}

// Load - Loads DAWG from a file.
func (d *DAWG) Load(path string) error {
	d.dct = NewDictionary()

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	err = d.dct.Read(f)
	closeErr := f.Close()
	if err == nil {
		err = closeErr
	}

	return err
}

func (d *DAWG) hasValue(index uint32) bool {
	return d.dct.hasValue(index)
}

func (d *DAWG) similarKeys(currentPrefix string, key string, index uint32, replaceChars map[rune]rune) []string {
	res := []string{}
	exitByBreak := false
	startPos := len(currentPrefix)

	for curPos, bStep := range key[startPos:] {
		ReplaceChar, ok := replaceChars[bStep]

		if ok {
			nextIndex, ok := d.dct.followBytes([]byte(string(ReplaceChar)), index)
			if ok {
				prefix := currentPrefix + key[startPos:curPos] + string(ReplaceChar)
				extraKeys := d.similarKeys(prefix, key, nextIndex, replaceChars)
				res = append(res, extraKeys...)
			}
		}

		index, ok = d.dct.followBytes([]byte(string(bStep)), index)
		if !ok {
			exitByBreak = true
			break
		}
	}

	if !exitByBreak {
		if d.vDAWG.hasValue(index) {
			foundKey := currentPrefix + key[startPos:]
			res = append([]string{foundKey}, res...)
		}
	}

	return res
}

// SimilarKeys - Returns all variants of 'key' in this DAWG according to 'replaces'.
// This may be useful e.g. for handling single-character umlauts.
func (d *DAWG) SimilarKeys(key string, replaceChars map[rune]rune) []string {
	return d.similarKeys("", key, constRoot, replaceChars)
}

// Prefixes - Returns a list with keys of this DAWG that are prefixes of the 'key'.
func (d *DAWG) Prefixes(key string) []string {
	res := []string{}
	var ok bool
	index := constRoot

	for pos, ch := range []byte(key) {
		index, ok = d.dct.followChar(uint32(ch), index)
		if !ok {
			break
		}

		if d.vDAWG.hasValue(index) {
			res = append(res, string(key[:pos+1]))
		}
	}

	return res
}

// NewDAWG - constructor for DAWG
func NewDAWG() *DAWG {
	dawg := &DAWG{dct: nil}
	dawg.vDAWG = dawg

	return dawg
}
