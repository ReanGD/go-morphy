package dawg

// status ok
import (
	"encoding/binary"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestBytesDAWG ...
func TestBytesDAWG(t *testing.T) {
	Convey("Suite setup", t, func() {
		data := []StringBytes{
			StringBytes{"bar", []byte("data2")},
			StringBytes{"foo", []byte("data1")},
			StringBytes{"foo", []byte("data3")},
			StringBytes{"foobar", []byte("data4")}}

		dawg := NewBytesDAWG()
		err := dawg.Load(testFullPath("small/bytes.dawg"))
		So(err, ShouldBeNil)

		Convey("Contains", func() {
			for _, item := range data {
				if !dawg.Contains(item.Key) {
					So(item.Key, ShouldEqual, "not contains")
				}
			}
			So(dawg.Contains("food"), ShouldBeFalse)
			So(dawg.Contains("x"), ShouldBeFalse)
			So(dawg.Contains("fo"), ShouldBeFalse)
		})

		Convey("Getitem", func() {
			res, ok := dawg.Get("foo")
			So(ok, ShouldBeTrue)
			So(res, ShouldResemble, [][]byte{[]byte("data1"), []byte("data3")})

			res, ok = dawg.Get("bar")
			So(ok, ShouldBeTrue)
			So(res, ShouldResemble, [][]byte{[]byte("data2")})

			res, ok = dawg.Get("foobar")
			So(ok, ShouldBeTrue)
			So(res, ShouldResemble, [][]byte{[]byte("data4")})
		})

		Convey("Getitem_missing", func() {
			_, ok := dawg.Get("x")
			So(ok, ShouldBeFalse)

			_, ok = dawg.Get("food")
			So(ok, ShouldBeFalse)

			_, ok = dawg.Get("foobarz")
			So(ok, ShouldBeFalse)

			_, ok = dawg.Get("f")
			So(ok, ShouldBeFalse)
		})

		Convey("Keys", func() {
			So(dawg.Keys(""), ShouldResemble, []string{"bar", "foo", "foo", "foobar"})
		})

		Convey("Key completion", func() {
			So(dawg.Keys("fo"), ShouldResemble, []string{"foo", "foo", "foobar"})
		})

		Convey("Items", func() {
			So(dawg.Items("xxx"), ShouldResemble, []StringBytes{})
			So(dawg.Items("fo"), ShouldResemble,
				[]StringBytes{
					StringBytes{"foo", []byte("data1")},
					StringBytes{"foo", []byte("data3")},
					StringBytes{"foobar", []byte("data4")}})
			So(dawg.Items(""), ShouldResemble, data)
		})

		Convey("Items completion", func() {
			So(dawg.Items("foob"), ShouldResemble,
				[]StringBytes{StringBytes{"foobar", []byte("data4")}})
		})

		Convey("Prefixes", func() {
			So(dawg.Prefixes("foobarz"), ShouldResemble, []string{"foo", "foobar"})
			So(dawg.Prefixes("bar"), ShouldResemble, []string{"bar"})
			So(dawg.Prefixes("x"), ShouldResemble, []string{})
		})
	})
}

// TestRecordDAWG ...
func TestRecordDAWG(t *testing.T) {
	Convey("Suite setup", t, func() {
		dawg := NewRecordDAWG(3, binary.BigEndian)
		err := dawg.Load(testFullPath("small/record.dawg"))
		So(err, ShouldBeNil)

		Convey("Getitem", func() {
			res, ok := dawg.Get("foo")
			So(ok, ShouldBeTrue)
			So(res, ShouldResemble, [][]uint16{[]uint16{3, 2, 1}, []uint16{3, 2, 256}})

			res, ok = dawg.Get("bar")
			So(ok, ShouldBeTrue)
			So(res, ShouldResemble, [][]uint16{[]uint16{3, 1, 0}})

			res, ok = dawg.Get("foobar")
			So(ok, ShouldBeTrue)
			So(res, ShouldResemble, [][]uint16{[]uint16{6, 3, 0}})
		})

		Convey("Getitem missing", func() {
			_, ok := dawg.Get("x")
			So(ok, ShouldBeFalse)

			_, ok = dawg.Get("food")
			So(ok, ShouldBeFalse)

			_, ok = dawg.Get("foobarz")
			So(ok, ShouldBeFalse)

			_, ok = dawg.Get("f")
			So(ok, ShouldBeFalse)
		})

		Convey("Record items", func() {
			data := []StringUints16{
				StringUints16{"bar", []uint16{3, 1, 0}},
				StringUints16{"foo", []uint16{3, 2, 1}},
				StringUints16{"foo", []uint16{3, 2, 256}},
				StringUints16{"foobar", []uint16{6, 3, 0}}}

			So(dawg.Items(""), ShouldResemble, data)
		})

		Convey("Record keys", func() {
			So(dawg.Keys(""), ShouldResemble,
				[]string{"bar", "foo", "foo", "foobar"})
		})

		Convey("Record keys prefix", func() {
			So(dawg.Keys("fo"), ShouldResemble, []string{"foo", "foo", "foobar"})
			So(dawg.Keys("bar"), ShouldResemble, []string{"bar"})
			So(dawg.Keys("barz"), ShouldResemble, []string{})
		})

		Convey("Prefixes", func() {
			So(dawg.Prefixes("foobarz"), ShouldResemble, []string{"foo", "foobar"})
			So(dawg.Prefixes("x"), ShouldResemble, []string{})
			So(dawg.Prefixes("bar"), ShouldResemble, []string{"bar"})
		})
	})
}
