package dawg

// status ok
import (
	"os"
	"testing"
	"unicode/utf8"

	. "github.com/smartystreets/goconvey/convey"
)

// TestDictionary ...
func TestDictionary(t *testing.T) {
	Convey("Suite setup", t, func() {
		words := testWords100k()
		So(len(words), ShouldEqual, 100000)

		f, err := os.Open(testFullPath("large/int_dawg.dawg"))
		So(err, ShouldBeNil)
		defer testSafeClose(f)

		dict := NewDictionary()
		So(dict.Read(f), ShouldBeNil)

		Convey("Contains", func() {
			for _, word := range words {
				if !dict.Contains([]byte(word)) {
					So(word, ShouldEqual, "not contains")
				}
			}
		})

		Convey("Find", func() {
			for _, word := range words {
				index, ok := dict.Find([]byte(word))
				if !ok || index != uint32(utf8.RuneCountInString(word)) {
					So(word, ShouldEqual, "not found")
				}
			}
		})
	})
}
