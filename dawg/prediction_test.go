package dawg

// status ok
import (
	"encoding/binary"
	"testing"
	"unicode/utf8"

	. "github.com/smartystreets/goconvey/convey"
)

// TestPrediction ...
func TestPrediction(t *testing.T) {
	Convey("Suite setup", t, func() {
		replaces := map[rune]rune{'Е': 'Ё'}

		type testSuiteData struct {
			word       string
			prediction []string
		}
		suite := []testSuiteData{
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
				So(dawg.SimilarKeys(d.word, replaces),
					ShouldResemble, d.prediction)
			}
		})

		Convey("record dawg prediction", func() {
			for _, d := range suite {
				So(recordDAWG.SimilarKeys(d.word, replaces),
					ShouldResemble, d.prediction)
			}
		})

		Convey("record dawg items", func() {
			for _, d := range suite {
				predictionItem := make([]StringUints16Arr, len(d.prediction))
				for i, word := range d.prediction {
					lenWord := uint16(utf8.RuneCountInString(word))
					predictionItem[i] = StringUints16Arr{word, [][]uint16{{lenWord}}}
				}
				So(recordDAWG.SimilarItems(d.word, replaces),
					ShouldResemble, predictionItem)
			}
		})

		Convey("record dawg items values", func() {
			for _, d := range suite {
				predictionItem := make([][][]uint16, len(d.prediction))
				for i, word := range d.prediction {
					lenWord := uint16(utf8.RuneCountInString(word))
					predictionItem[i] = [][]uint16{{lenWord}}
				}
				So(recordDAWG.SimilarItemsValues(d.word, replaces),
					ShouldResemble, predictionItem)
			}
		})
	})
}
