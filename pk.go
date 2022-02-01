package pk

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/huandu/xstrings"
	"github.com/uber/h3-go"
)

const (
	RESOLUTION        = 10
	BASE_RESOLUTION   = 12
	EARTH_RADIUS      = 6371 // km
	ALPHABET_BASE     = "23456789BCDFGHJKMNPQRSTVWXYZ"
	CODE_LENGTH       = 9
	TUPLE_LENGTH      = 3
	PADDING_CHAR      = "a"
	REPLACEMENT_CHARS = "eu"
)

var (
	BASE_CELL_SHIFT          = int64(math.Pow(2, 45)) // Adding this will increment the base cell value by 1
	UNUSED_RESOLUTION_FILLER = int64(math.Pow(2, (3*(15-BASE_RESOLUTION))-1))
	ALPHABET                 string
	ALPHABET_LENGTH          int64
	HEADER_BITS              string
	HEADER_INT               int64 = 0
	FIRST_TUPLE_REGEX        string
	TUPLE_REGEX              string
	REPLACEMENT_MAP          = map[string]string{
		"prn":   "pre",
		"f4nny": "f4nne",
		"tw4t":  "tw4e",
		"ngr":   "ngu", // 'u' avoids introducing 'gey'
		"dck":   "dce",
		"vjn":   "vju", // 'u' avoids introducing 'jew'
		"fck":   "fce",
		"pns":   "pne",
		"sht":   "she",
		"kkk":   "kke",
		"fgt":   "fgu", // 'u' avoids introducing 'gey'
		"dyk":   "dye",
		"bch":   "bce",
	}
)

func init() {
	init_ALPHABET_LENGTH()
	init_H3_HEADER()
}

func init_ALPHABET_LENGTH() {
	ALPHABET = strings.ToLower(ALPHABET_BASE)
	ALPHABET_LENGTH = int64(len(ALPHABET))

	FIRST_TUPLE_REGEX = "[" + ALPHABET + REPLACEMENT_CHARS + PADDING_CHAR + "]{3}"
	TUPLE_REGEX = "[" + ALPHABET + REPLACEMENT_CHARS + "]{3}"
}

func zfill(rawString string, padString string, expectLength int) string {
	diff := expectLength - len(rawString)
	if diff > 0 {
		rawString = strings.Repeat(padString, diff) + rawString
	}
	return rawString
}

func init_H3_HEADER() {
	idx := strconv.FormatInt(
		int64(
			h3.FromGeo(
				h3.GeoCoord{Latitude: 0, Longitude: 0},
				RESOLUTION,
			)),
		2)
	// Python bin(xxx) result has prefix "0bxxxxx"
	// Golang FormatInt does not
	// filled := xstrings.LeftJustify(idx, 64, "0")
	filled := zfill(idx, "0", 64)
	bits := filled[:12]
	if bits != "000010001010" {
		panic(errors.New(bits))
	}
	HEADER_BITS = bits

	for index, char := range xstrings.Reverse(HEADER_BITS) {
		intV, err := strconv.ParseInt(string(char), 10, 0)
		if err != nil {
			panic(err)
		}
		HEADER_INT = HEADER_INT + intV*int64(math.Pow(2, float64(index)))
	}
	HEADER_INT = HEADER_INT * int64(math.Pow(2, 52))
	if HEADER_INT != 621496748577128448 {
		panic(HEADER_INT)
	}
}

func GeoDistance(lat1, long1, lat2, long2 float64) float64 {
	hav_lat := 0.5 * (1 - math.Cos(lat1-lat2))
	hav_long := 0.5 * (1 - math.Cos(long1-long2))
	radical := math.Sqrt(hav_lat + math.Cos(lat1)*math.Cos(lat2)*hav_long)
	return 2 * EARTH_RADIUS * math.Asin(radical) * 1000
}

// ShortenH3Integer shorten an H3 integer to only include location data up to the base resolution
func ShortenH3Integer(h3Int int64) int64 {
	// Cuts off the 12 left-most bits that don't code location
	out := (h3Int + BASE_CELL_SHIFT) % int64(math.Pow(2, 52))
	// Cuts off the rightmost bits corresponding to resolutions greater than the base resolution
	out = out >> (3 * (15 - BASE_RESOLUTION))
	return out
}

func UnshortenH3Integer(shortInt int64) int64 {
	unshifted_int := shortInt << (3 * (15 - BASE_RESOLUTION))
	rebuilt_int := HEADER_INT + UNUSED_RESOLUTION_FILLER - BASE_CELL_SHIFT + unshifted_int
	return rebuilt_int
}

func EncodeH3Int(h3Int int64) string {
	shortH3Int := ShortenH3Integer(h3Int)
	encoedH3Int := EncodeShortInt(shortH3Int)
	cleanEncoedShortH3 := CleanString(encoedH3Int)
	if len(cleanEncoedShortH3) <= CODE_LENGTH {
		cleanEncoedShortH3 = xstrings.RightJustify(cleanEncoedShortH3, CODE_LENGTH, PADDING_CHAR)
	}
	parts := []string{}
	for i := 0; i < len(cleanEncoedShortH3); i += TUPLE_LENGTH {
		parts = append(parts, cleanEncoedShortH3[i:i+TUPLE_LENGTH])
	}
	return "@" + strings.Join(parts, "-")
}

func EncodeShortInt(shortInt int64) string {
	if shortInt == 0 {
		return fmt.Sprintf("%c", ALPHABET[0])
	}
	res := ""
	for shortInt > 0 {
		remainder := shortInt % ALPHABET_LENGTH
		res = fmt.Sprintf("%c", ALPHABET[remainder]) + res
		shortInt = int64(shortInt / ALPHABET_LENGTH)
	}
	return res
}

func CleanString(s string) string {
	for k, v := range REPLACEMENT_MAP {
		s = strings.Replace(s, k, v, -1)
	}
	return s
}

func GeoToPlacekey(lat, long float64) string {
	return EncodeH3Int(int64(h3.FromGeo(
		h3.GeoCoord{Latitude: lat, Longitude: long},
		RESOLUTION,
	)))
}
