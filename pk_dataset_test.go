package pk_test

import (
	"strconv"
	"strings"
	"testing"

	_ "embed"

	"github.com/ringsaturn/pk"
)

//go:embed example_geos.csv
var exampleGeosCSV []byte

func TestDataset(t *testing.T) {
	lines := strings.Split(string(exampleGeosCSV), "\n")
	for index, line := range lines {
		if index == 0 {
			continue
		}
		rawparts := strings.Split(line, ",")
		if len(rawparts) != 8 {
			continue
		}
		lat_str := rawparts[0]
		long_str := rawparts[1]
		// h3_r10_str := rawparts[2]
		// h3_int_r10_str := rawparts[3]
		placekey_str := rawparts[4]
		// h3_lat_str := rawparts[5]
		// h3_long_str := rawparts[6]
		// info_str := rawparts[7]
		// log.Println(lat_str, long_str, h3_r10_str, h3_int_r10_str, placekey_str, h3_lat_str, h3_long_str, info_str)

		lat_float, err := strconv.ParseFloat(lat_str, 64)
		if err != nil {
			panic(err)
		}
		long_float, err := strconv.ParseFloat(long_str, 64)
		if err != nil {
			panic(err)
		}
		placeKey, err := pk.GeoToPlacekey(lat_float, long_float)
		if err != nil {
			panic(err)
		}
		if placeKey != placekey_str {
			t.Errorf("bad result %v got %v\n", line, placeKey)
		}
	}
}
