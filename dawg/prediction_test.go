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
			testSuiteData{"УЖ", []string{}},
			testSuiteData{"ЕМ", []string{"ЕМ"}},
			testSuiteData{"ЁМ", []string{}},
			testSuiteData{"ЁЖ", []string{"ЁЖ"}},
			testSuiteData{"ЕЖ", []string{"ЁЖ"}},
			testSuiteData{"ЁЖИК", []string{"ЁЖИК"}},
			testSuiteData{"ЕЖИКЕ", []string{"ЁЖИКЕ"}},
			testSuiteData{"ДЕРЕВНЯ", []string{"ДЕРЕВНЯ", "ДЕРЁВНЯ"}},
			testSuiteData{"ДЕРЁВНЯ", []string{"ДЕРЁВНЯ"}},
			testSuiteData{"ОЗЕРА", []string{"ОЗЕРА", "ОЗЁРА"}},
			testSuiteData{"ОЗЕРО", []string{"ОЗЕРО"}},
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
					predictionItem[i] = StringUints16Arr{word, [][]uint16{[]uint16{lenWord}}}
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
					predictionItem[i] = [][]uint16{[]uint16{lenWord}}
				}
				So(recordDAWG.SimilarItemsValues(d.word, replaces),
					ShouldResemble, predictionItem)
			}
		})
	})
}
