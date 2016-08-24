package dawg

// status ok
import "github.com/ReanGD/go-morphy/std"

// IntCompletionDAWG - Dict-like class based on DAWG.
// It can store integer values for unicode keys and support key completion.
type IntCompletionDAWG struct {
	DAWG
	CompletionDAWG
}

// Get - Return value for the given key.
func (d *IntCompletionDAWG) Get(key string) (uint32, bool) {
	return d.dct.Find([]byte(key))
}

// Items ...
func (d *IntCompletionDAWG) Items(prefix string) []std.StrUint32 {
	res := []std.StrUint32{}
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
		res = append(res, std.StrUint32{Key: string(completer.key), Value: completer.value()})
	}

	return res
}

// Load ...
func (d *IntCompletionDAWG) Load(path string) error {
	err := d.CompletionDAWG.Load(path)
	if err != nil {
		return err
	}
	d.DAWG.dct = d.CompletionDAWG.dct

	return nil
}

func (d *IntCompletionDAWG) initIntCompletionDAWG() {
	d.initDAWG()
	d.initCompletionDAWG()
}

// NewIntCompletionDAWG - constructor for IntCompletionDAWG
func NewIntCompletionDAWG() *IntCompletionDAWG {
	dawg := &IntCompletionDAWG{}
	dawg.DAWG.vDAWG = dawg
	dawg.CompletionDAWG.vDAWG = dawg
	dawg.initIntCompletionDAWG()

	return dawg
}
