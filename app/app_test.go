package app_test

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/cockroachdb/pebble"
	dbm "github.com/cometbft/cometbft-db"
	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
)

var moduleRe = regexp.MustCompile(`s\/k:(\w+)\/`)

func main() {

	fmt.Println("trying to open db")

	pebbleDB, err := dbm.NewPebbleDBWithOpts("application", "/Users/tanmay/Downloads/dataM/", &pebble.Options{
		ReadOnly: true,
	})
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}

	//db, err := dbm.NewDB("application", dbm.PebbleDBBackend, "/Users/tanmay/Downloads/data/")
	//if err != nil {
	//	log.Fatalf("failed to open DB: %v", err)
	//}
	//fmt.Println("db opened")
	//
	//pebbleDB, ok := db.(*dbm.PebbleDB)
	//if !ok {
	//	log.Fatalf("invalid logical DB type; expected: %T, got: %T", &dbm.GoLevelDB{}, db)
	//}

	//for key, val := range pebbleDB.Stats() {
	//	fmt.Printf("%s: %v\n", key, val)
	//
	//}

	//levelDBStats, err := goLevelDB.DB().GetProperty("leveldb.stats")
	//if err != nil {
	//	log.Fatalf("failed to get LevelDB stats: %v", err)
	//}
	//
	//fmt.Printf("%s\n", levelDBStats)

	var (
		totalKeys    int
		totalKeySize int
		totalValSize int

		moduleStats = make(map[string][]int)
	)

	iter, err := pebbleDB.DB().NewIter(nil)
	if err != nil {
		log.Fatalf("failed to create iterator: %v", err)
	}

	// Iterate through all keys
	//for iter.First(); iter.Valid(); iter.Next() {
	//	key := iter.Key()
	//	value := iter.Value()
	//	fmt.Printf("Key: %s, Value: %s\n", key, value)
	//}

	if err := iter.Error(); err != nil {
		log.Fatalf("iterator error: %v", err)
	}
	count := 0
	for iter.First(); iter.Valid(); iter.Next() {
		keySize := len(iter.Key())
		valSize := len(iter.Value())

		totalKeys++
		totalKeySize += keySize
		totalValSize += valSize

		var statKey string
		//if strings.Contains(string(iter.Value()), "begin_block_events") {
		//	fmt.Println(string(iter.Value()))
		//	os.Exit(0)
		//}

		keyStr := string(iter.Key())
		if strings.HasPrefix(keyStr, "s/k:") {
			tokens := moduleRe.FindStringSubmatch(keyStr)
			statKey = tokens[1]
		} else {
			statKey = "misc"
		}

		if statKey == "staking" {
			fmt.Println(keyStr)
			count++
			if count == 30 {
				os.Exit(0)
			}
		}

		if moduleStats[statKey] == nil {
			// XXX/TODO: Move this into a struct
			//
			// 0: total set size
			// 1: total key size
			// 2: total value size
			moduleStats[statKey] = make([]int, 3)
		}

		moduleStats[statKey][0]++
		moduleStats[statKey][1] += keySize
		moduleStats[statKey][2] += valSize
	}

	// print application-specific stats
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Module", "Avg Key Size", "Avg Value Size", "Total Key Size", "Total Value Size", "Total Key Pairs"})

	modules := maps.Keys(moduleStats)
	SortSlice(modules)

	for _, m := range modules {
		stats := moduleStats[m]
		t.AppendRow([]interface{}{
			m,
			ByteCountDecimal(stats[1] / stats[0]),
			ByteCountDecimal(stats[2] / stats[0]),
			ByteCountDecimal(stats[1]),
			ByteCountDecimal(stats[2]),
			stats[0],
		})
	}

	t.AppendFooter(table.Row{"Total", "", "", ByteCountDecimal(totalKeySize), ByteCountDecimal(totalValSize), totalKeys})

	t.Render()
}

func SortSlice[T constraints.Ordered](s []T) {
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
}

func ByteCountDecimal(b int) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := int64(b) / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
