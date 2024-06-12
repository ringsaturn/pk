package pk_test

import (
	"math"
	"strconv"
	"strings"
	"testing"

	_ "embed"

	"github.com/ringsaturn/pk"
)

//go:embed example_geos.csv
var exampleGeosCSV []byte

//go:embed example_distances.tsv
var exampleDistanceTSV []byte

func TestPlaceKeyDataset(t *testing.T) {
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
		lng_str := rawparts[1]
		placekey_str := rawparts[4]

		lat_float, err := strconv.ParseFloat(lat_str, 64)
		if err != nil {
			panic(err)
		}
		lng_float, err := strconv.ParseFloat(lng_str, 64)
		if err != nil {
			panic(err)
		}
		placeKey, err := pk.GeoToPlacekey(lat_float, lng_float)
		if err != nil {
			panic(err)
		}
		if placeKey != placekey_str {
			t.Errorf("bad result %v got %v\n", line, placeKey)
		}
	}
}

func TestGeoDistanceDataset(t *testing.T) {
	lines := strings.Split(string(exampleDistanceTSV), "\n")
	for index, line := range lines {
		if index == 0 {
			continue
		}
		rawparts := strings.Split(line, "\t")
		if len(rawparts) != 8 {
			continue
		}
		// placekey_1	geo_1	placekey_2	geo_2	distance(km)	error
		pk1 := rawparts[0]
		pk2 := rawparts[2]
		expectDistStr := rawparts[4]
		expectDistErrStr := rawparts[5]

		expectDist, err := strconv.ParseFloat(expectDistStr, 64)
		if err != nil {
			panic(err)
		}

		expectDistErr, err := strconv.ParseFloat(expectDistErrStr, 64)
		if err != nil {
			panic(err)
		}

		dist, err := pk.PlacekeyDistance(pk1, pk2)
		if err != nil {
			t.Error(err)
			continue
		}
		if math.Abs(dist/1000-expectDist) > expectDistErr {
			t.Error("bad")
		}
	}
}
