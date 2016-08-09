package dawg

// status ok
import "os"

// CompletionDAWG - DAWG with key completion support.
type CompletionDAWG struct {
	*DAWG
	guide *Guide
}

// Keys ...
func (d *CompletionDAWG) Keys(prefix string) []string {
	bPrefix := []byte(prefix)
	res := []string{}

	index, ok := d.dct.followBytes(bPrefix, constRoot)
	if !ok {
		return res
	}

	completer := NewCompleter(d.dct, d.guide)
	completer.start(index, bPrefix)

	for completer.next() {
		res = append(res, string(completer.key))
	}

	return res
}

// Load - Loads DAWG from a file.
func (d *CompletionDAWG) Load(path string) error {
	d.dct = NewDictionary()
	d.guide = NewGuide()

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	err = d.dct.Read(f)
	if err == nil {
		err = d.guide.read(f)
	}
	closeErr := f.Close()
	if err == nil {
		err = closeErr
	}

	return err
}

// NewCompletionDAWG - constructor for CompletionDAWG
func NewCompletionDAWG() *CompletionDAWG {
	dawg := &CompletionDAWG{DAWG: NewDAWG(), guide: nil}
	dawg.vDAWG = dawg

	return dawg
}
