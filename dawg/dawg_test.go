package dawg

// status ok
import (
	"archive/zip"
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"testing"

	"github.com/ReanGD/go-morphy/std"
	. "github.com/smartystreets/goconvey/convey"
)

func testFullPath(path string) string {
	return "../DAWG-Python/dev_data/" + path
}

func testSafeClose(buf io.Closer) {
	So(buf.Close(), ShouldBeNil)
}

func testSafeRemoveFile(path string) {
	So(os.Remove(path), ShouldBeNil)
}

func testWords100k() []string {
	r, err := zip.OpenReader(testFullPath("words100k.txt.zip"))
	So(err, ShouldBeNil)
	defer testSafeClose(r)

	f := r.File[0]
	rc, err := f.Open()
	So(err, ShouldBeNil)
	defer testSafeClose(rc)

	var result []string
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	So(scanner.Err(), ShouldBeNil)

	return result
}

// TestDAWG ...
func TestDAWG(t *testing.T) {
	Convey("Suite setup", t, func() {
		Convey("Invalid file path", func() {
			dawgErr := NewDAWG()
			err := dawgErr.Load("err.dawg")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "open err.dawg: no such file or directory")
		})
	})
}

// TestCompletionDAWG ...
func TestCompletionDAWG(t *testing.T) {
	Convey("Suite setup", t, func() {
		keys := []string{"f", "bar", "foo", "foobar"}

		dawg := NewCompletionDAWG()
		err := dawg.Load(testFullPath("small/completion.dawg"))
		So(err, ShouldBeNil)

		Convey("Contains", func() {
			for _, key := range keys {
				if !dawg.Contains(key) {
					So(key, ShouldEqual, "not contains")
				}
			}
		})

		Convey("Keys", func() {
			sortedKeys := keys
			sort.Strings(sortedKeys)
			So(dawg.Keys(""), ShouldResemble, sortedKeys)
		})

		Convey("Completion", func() {
			So(dawg.Keys("z"), ShouldResemble, []string{})
			So(dawg.Keys("b"), ShouldResemble, []string{"bar"})
			So(dawg.Keys("foo"), ShouldResemble, []string{"foo", "foobar"})
		})

		Convey("Invalid file", func() {
			tmpfile, err := ioutil.TempFile("", "dawg_error")
			So(err, ShouldBeNil)

			defer testSafeRemoveFile(tmpfile.Name())

			_, err = tmpfile.Write([]byte("foo"))
			So(err, ShouldBeNil)
			err = tmpfile.Close()
			So(err, ShouldBeNil)

			dawgErr := NewCompletionDAWG()
			err = dawgErr.Load(tmpfile.Name())
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "unexpected EOF")
		})

		Convey("Empty dawg", func() {
			d := NewCompletionDAWG()
			err := d.Load(testFullPath("small/completion-empty.dawg"))
			So(err, ShouldBeNil)

			So(d.Keys(""), ShouldResemble, []string{})
		})

		Convey("Prefixes", func() {
			So(dawg.Prefixes("foobar"), ShouldResemble,
				[]string{"f", "foo", "foobar"})
			So(dawg.Prefixes("x"), ShouldResemble, []string{})
			So(dawg.Prefixes("bar"), ShouldResemble, []string{"bar"})
		})

		Convey("Invalid file path", func() {
			dawgErr := NewCompletionDAWG()
			err := dawgErr.Load("err.dawg")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "open err.dawg: no such file or directory")
		})
	})
}

// TestIntDAWG ...
func TestIntDAWG(t *testing.T) {
	Convey("getitem", t, func() {
		payload := map[string]uint32{"foo": 1, "bar": 5, "foobar": 3}

		dawg := NewIntCompletionDAWG()
		err := dawg.DAWG.Load(testFullPath("small/int_dawg.dawg"))
		So(err, ShouldBeNil)

		for key, value := range payload {
			val, ok := dawg.Get(key)
			So(ok, ShouldBeTrue)
			So(val, ShouldEqual, value)
		}

		_, ok := dawg.Get("fo")
		So(ok, ShouldBeFalse)
	})
}

// TestIntCompletionDawg ...
func TestIntCompletionDawg(t *testing.T) {
	Convey("Suite setup", t, func() {
		payload := []std.StrUint32{{"bar", 5}, {"foo", 1}, {"foobar", 3}}

		dawg := NewIntCompletionDAWG()
		err := dawg.Load(testFullPath("small/int_completion_dawg.dawg"))
		So(err, ShouldBeNil)

		Convey("getitem", func() {
			for _, it := range payload {
				val, ok := dawg.Get(it.Key)
				So(ok, ShouldBeTrue)
				So(val, ShouldEqual, it.Value)
			}

			_, ok := dawg.Get("fo")
			So(ok, ShouldBeFalse)
		})

		Convey("completion keys", func() {
			So(dawg.Keys(""), ShouldResemble, []string{"bar", "foo", "foobar"})
		})

		Convey("completion keys with prefix", func() {
			So(dawg.Keys("fo"), ShouldResemble, []string{"foo", "foobar"})
			So(dawg.Keys("foo"), ShouldResemble, []string{"foo", "foobar"})
			So(dawg.Keys("foob"), ShouldResemble, []string{"foobar"})
			So(dawg.Keys("z"), ShouldResemble, []string{})
			So(dawg.Keys("b"), ShouldResemble, []string{"bar"})
		})

		Convey("completion items", func() {
			So(dawg.Items(""), ShouldResemble, payload)
		})

		Convey("completion items with prefix", func() {
			So(dawg.Items("foo"), ShouldResemble, []std.StrUint32{{"foo", 1}, {"foobar", 3}})
			So(dawg.Items("x"), ShouldResemble, []std.StrUint32{})
		})

		Convey("Invalid file path", func() {
			dawgErr := NewIntCompletionDAWG()
			err := dawgErr.Load("err.dawg")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "open err.dawg: no such file or directory")
		})
	})
}
