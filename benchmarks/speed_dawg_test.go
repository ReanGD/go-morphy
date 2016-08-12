package benchmarks

import (
	"fmt"
	"strings"
	"testing"
)

type testOp struct {
	opName    string
	descr     string
	opCoeff   int
	repeats   int
	runs      int
	benchFunc func(int)
}

type benchContext struct {
	Words100k    []string
	NonWords100k []string
	NonWords10k  []string
	data         benchInterface
}

var benchCtx *benchContext

func initBenchCtx() {
	words100k := words100k()
	nonWords100k := randomWords(100000)
	benchCtx = &benchContext{
		Words100k:    words100k,
		NonWords100k: nonWords100k,
		NonWords10k:  nonWords100k[:10000],
	}
}

func getHitBench(repeats int) {
	data := benchCtx.data
	for i := 0; i < repeats; i++ {
		for _, word := range benchCtx.Words100k {
			if !data.GetWrap(word) {
				panic("not get " + word)
			}
		}
	}
}

func getMissesBench(repeats int) {
	data := benchCtx.data
	for i := 0; i < repeats; i++ {
		for _, word := range benchCtx.NonWords10k {
			_ = data.GetWrap(word)
		}
	}
}

func containsHitBench(repeats int) {
	data := benchCtx.data
	for i := 0; i < repeats; i++ {
		for _, word := range benchCtx.Words100k {
			if !data.Contains(word) {
				panic("not contains " + word)
			}
		}
	}
}

func containsMissesBench(repeats int) {
	data := benchCtx.data
	for i := 0; i < repeats; i++ {
		for _, word := range benchCtx.NonWords100k {
			_ = data.Contains(word)
		}
	}
}

func itemsBench(repeats int) {
	data := benchCtx.data
	for i := 0; i < repeats; i++ {
		_ = data.ItemsWrap()
	}
}

func keysBench(repeats int) {
	data := benchCtx.data
	for i := 0; i < repeats; i++ {
		_ = data.keysWrap()
	}
}

func bench(data benchInterface, containerName string, test testOp) {
	op := test.opName
	container := strings.TrimSpace(containerName)
	skip := container == "DAWG" &&
		(op == "get() (hits)" || op == "get() (misses)" || op == "items()" || op == "keys()")
	skip = skip || (container == "dict" && (op == "items()" || op == "keys()"))
	skip = skip || (container == "IntDAWG" && (op == "items()" || op == "keys()"))

	benchCtx.data = data
	fullName := fmt.Sprintf("%s %s:\t", containerName, test.opName)
	if len(test.opName) < 21 {
		fullName += "\t"
	}
	if len(test.opName) < 8 {
		fullName += "\t"
	}
	if skip {
		fmt.Printf("%snot supported\n", fullName)
	} else {
		sec := runBenchmark(test.runs, test.repeats, test.benchFunc)

		cnt := float64(test.repeats) / (float64(sec) * float64(test.opCoeff))
		fmt.Printf("%s%0.3f%s\n", fullName, cnt, test.descr)
	}
}

// BenchmarkDAWG ...
func BenchmarkDAWG(b *testing.B) {
	b.N = 0
	initBenchCtx()
	dict := loadDict(benchCtx.Words100k)
	DAWG := loadDAWG()
	bytesDAWG := loadBytesDAWG()
	recordDAWG := loadRecordDAWG()
	intDAWG := loadIntDAWG()

	tests := []testOp{
		testOp{"get() (hits)", "M ops/sec", 10, 3, 5, getHitBench},
		testOp{"get() (misses)", "M ops/sec", 100, 5, 5, getMissesBench},
		testOp{"__contains__ (hits)", "M ops/sec", 10, 3, 5, containsHitBench},
		testOp{"__contains__ (misses)", "M ops/sec", 10, 3, 5, containsMissesBench},
		testOp{"items()", " ops/sec", 1, 1, 5, itemsBench},
		testOp{"keys()", " ops/sec", 1, 1, 5, keysBench},
	}

	fmt.Printf("\n====== Benchmarks (100k unique unicode words) =======\n\n")

	for _, test := range tests {
		bench(dict, "      dict", test)
		bench(DAWG, "      DAWG", test)
		bench(bytesDAWG, " BytesDAWG", test)
		bench(recordDAWG, "RecordDAWG", test)
		bench(intDAWG, "   IntDAWG", test)
		fmt.Println()
	}
	b.SkipNow()
}
