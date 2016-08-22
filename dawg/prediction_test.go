package dawg

// status ok
import (
	"encoding/binary"
	"testing"
	"unicode/utf8"

	"github.com/ReanGD/go-morphy/std"
	. "github.com/smartystreets/goconvey/convey"
)

// TestPrediction ...
func TestPrediction(t *testing.T) {
	Convey("Suite setup", t, func() {
		replaces := map[rune]rune{'Е': 'Ё'}

		suite := []std.StrStrs{
			{"УЖ", []string{}},
			{"ЕМ", []string{"ЕМ"}},
			{"ЁМ", []string{}},
			{"ЁЖ", []string{"ЁЖ"}},
			{"ЕЖ", []string{"ЁЖ"}},
			{"ЁЖИК", []string{"ЁЖИК"}},
			{"ЕЖИКЕ", []string{"ЁЖИКЕ"}},
			{"ДЕРЕВНЯ", []string{"ДЕРЕВНЯ", "ДЕРЁВНЯ"}},
			{"ДЕРЁВНЯ", []string{"ДЕРЁВНЯ"}},
			{"ОЗЕРА", []string{"ОЗЕРА", "ОЗЁРА"}},
			{"ОЗЕРО", []string{"ОЗЕРО"}},
		}
		recordDAWG := NewRecordDAWG(1, binary.LittleEndian)
		err := recordDAWG.Load(testFullPath("small/prediction-record.dawg"))
		So(err, ShouldBeNil)

		dawg := NewDAWG()
		err = dawg.Load(testFullPath("small/prediction.dawg"))
		So(err, ShouldBeNil)

		Convey("dawg prediction", func() {
			for _, d := range suite {
				So(dawg.SimilarKeys(d.Key, replaces),
					ShouldResemble, d.Value)
			}
		})

		Convey("record dawg prediction", func() {
			for _, d := range suite {
				So(recordDAWG.SimilarKeys(d.Key, replaces),
					ShouldResemble, d.Value)
			}
		})

		Convey("record dawg items", func() {
			for _, d := range suite {
				predictionItem := make([]std.StrUints16Arr, len(d.Value))
				for i, word := range d.Value {
					lenWord := uint16(utf8.RuneCountInString(word))
					predictionItem[i] = std.StrUints16Arr{word, [][]uint16{{lenWord}}}
				}
				So(recordDAWG.SimilarItems(d.Key, replaces),
					ShouldResemble, predictionItem)
			}
		})

		Convey("record dawg items values", func() {
			for _, d := range suite {
				predictionItem := make([][][]uint16, len(d.Value))
				for i, word := range d.Value {
					lenWord := uint16(utf8.RuneCountInString(word))
					predictionItem[i] = [][]uint16{{lenWord}}
				}
				So(recordDAWG.SimilarItemsValues(d.Key, replaces),
					ShouldResemble, predictionItem)
			}
		})
	})
}
