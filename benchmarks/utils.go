package benchmarks

import (
	"archive/zip"
	"bufio"
	"encoding/binary"
	"io"
	"math/rand"
	"time"

	"github.com/ReanGD/go-morph/dawg"
)

func fullPath(path string) string {
	return "../DAWG-Python/dev_data/" + path
}

func safeClose(buf io.Closer) {
	err := buf.Close()
	if err != nil {
		panic(err)
	}
}

type benchInterface interface {
	Contains(key string) bool
	GetWrap(key string) bool
	ItemsWrap() bool
	keysWrap() bool
}

type dictWrap struct {
	data map[string]bool
}

func (d *dictWrap) Contains(key string) bool {
	_, ok := d.data[key]
	return ok
}

func (d *dictWrap) GetWrap(key string) bool {
	_, ok := d.data[key]
	return ok
}

func (d *dictWrap) ItemsWrap() bool {
	return true
}

func (d *dictWrap) keysWrap() bool {
	return true
}

func loadDict(words []string) *dictWrap {
	data := make(map[string]bool, len(words))
	for _, key := range words {
		data[key] = true
	}
	return &dictWrap{data: data}
}

type wrapDAWG struct {
	*dawg.DAWG
}

func (d *wrapDAWG) GetWrap(key string) bool {
	return true
}

func (d *wrapDAWG) ItemsWrap() bool {
	return true
}

func (d *wrapDAWG) keysWrap() bool {
	return true
}

func loadDAWG() *wrapDAWG {
	data := dawg.NewDAWG()
	err := data.Load(fullPath("large/dawg.dawg"))
	if err != nil {
		panic(err)
	}

	return &wrapDAWG{DAWG: data}
}

type wrapBytesDAWG struct {
	*dawg.BytesDAWG
}

func (d *wrapBytesDAWG) GetWrap(key string) bool {
	_, ok := d.Get(key)
	return ok
}

func (d *wrapBytesDAWG) ItemsWrap() bool {
	_ = d.Items("")
	return true
}

func (d *wrapBytesDAWG) keysWrap() bool {
	_ = d.Keys("")
	return true
}

func loadBytesDAWG() *wrapBytesDAWG {
	data := dawg.NewBytesDAWG()
	err := data.Load(fullPath("large/bytes_dawg.dawg"))
	if err != nil {
		panic(err)
	}

	return &wrapBytesDAWG{BytesDAWG: data}
}

type wrapRecordDAWG struct {
	*dawg.RecordDAWG
}

func (d *wrapRecordDAWG) GetWrap(key string) bool {
	_, ok := d.Get(key)
	return ok
}

func (d *wrapRecordDAWG) ItemsWrap() bool {
	_ = d.Items("")
	return true
}

func (d *wrapRecordDAWG) keysWrap() bool {
	_ = d.Keys("")
	return true
}

func loadRecordDAWG() *wrapRecordDAWG {
	data := dawg.NewRecordDAWG(1, binary.LittleEndian)
	err := data.Load(fullPath("large/record_dawg.dawg"))
	if err != nil {
		panic(err)
	}

	return &wrapRecordDAWG{RecordDAWG: data}
}

type wrapIntCompletionDAWG struct {
	*dawg.IntCompletionDAWG
}

func (d *wrapIntCompletionDAWG) GetWrap(key string) bool {
	_, ok := d.Get(key)
	return ok
}

func (d *wrapIntCompletionDAWG) ItemsWrap() bool {
	return true
}

func (d *wrapIntCompletionDAWG) keysWrap() bool {
	_ = d.Keys("")
	return true
}

func loadIntDAWG() *wrapIntCompletionDAWG {
	data := dawg.NewIntCompletionDAWG()
	err := data.DAWG.Load(fullPath("large/int_dawg.dawg"))
	if err != nil {
		panic(err)
	}

	return &wrapIntCompletionDAWG{IntCompletionDAWG: data}
}

func words100k() []string {
	result := []string{}

	r, err := zip.OpenReader(fullPath("words100k.txt.zip"))
	if err != nil {
		panic(err)
	}
	defer safeClose(r)

	f := r.File[0]
	rc, err := f.Open()
	if err != nil {
		panic(err)
	}
	defer safeClose(rc)

	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	return result
}

func randomWords(num int) []string {
	rand.Seed(time.Now().UTC().UnixNano())

	alphabet := []rune{}
	letters := "абвгдеёжзиклмнопрстуфхцчъыьэюяabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for _, r := range letters {
		alphabet = append(alphabet, r)
	}
	alphabetLen := int32(len(alphabet))

	res := make([]string, num)
	for i := 0; i != num; i++ {
		wordLen := int(rand.Int31n(15) + 1)
		s := ""
		for j := 0; j != wordLen; j++ {
			s += string(alphabet[rand.Int31n(alphabetLen)])
		}
		res[i] = s
	}
	return res
}
