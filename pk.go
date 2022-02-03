// Package pk provide part of https://github.com/Placekey/placekey-py features in pure Go.
// Most of codes translated from Python SDK.
// See details in https://docs.placekey.io/Placekey_Encoding_Specification_White_Paper.pdf
//
// Placekey has two parts: `what` and  `where`.
//
// The `where` part.
// The core of `where` part is cut 21 bits of 64-bit integer which is a H3 integer, only use 43 bits left.
// Because the resolution is a const `10`, so each H3 id under the same resolution
// could remove common info.
// Then use alphabet represent the 43-bit integer as the `where` part of the whole place key.
//
// The `what` part need use Placekey's API.
// Check API doc https://docs.placekey.io
package pk

import (
	"errors"
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
	ALPHABET                 string // build in init
	ALPHABET_LENGTH          int64  // build in init
	HEADER_BITS              string // build in init
	HEADER_INT               int64  // build in init
	FIRST_TUPLE_REGEX        string // build in init
	TUPLE_REGEX              string // build in init
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

// Degrees × π / 180

func radians(degree float64) float64 {
	return degree * math.Pi / 180
}

func GeoDistance(lat1, long1, lat2, long2 float64) float64 {
	lat1 = radians(lat1)
	long1 = radians(long1)
	lat2 = radians(lat2)
	long2 = radians(long2)

	hav_lat := 0.5 * (1 - math.Cos(lat1-lat2))
	hav_long := 0.5 * (1 - math.Cos(long1-long2))
	radical := math.Sqrt(hav_lat + math.Cos(lat1)*math.Cos(lat2)*hav_long)
	return 2 * EARTH_RADIUS * math.Asin(radical) * 1000
}

// shortenH3Integer shorten an H3 integer to only include location data up to the base resolution
func shortenH3Integer(h3Int int64) int64 {
	// Cuts off the 12 left-most bits that don't code location
	out := (h3Int + BASE_CELL_SHIFT) % int64(math.Pow(2, 52))
	// Cuts off the rightmost bits corresponding to resolutions greater than the base resolution
	out = out >> (3 * (15 - BASE_RESOLUTION))
	return out
}

func unshortenH3Integer(shortInt int64) int64 {
	unshifted_int := shortInt << (3 * (15 - BASE_RESOLUTION))
	rebuilt_int := HEADER_INT + UNUSED_RESOLUTION_FILLER - BASE_CELL_SHIFT + unshifted_int
	return rebuilt_int
}

func encodeH3Int(h3Int int64) string {
	shortH3Int := shortenH3Integer(h3Int)
	encoedH3Int := encodeShortInt(shortH3Int)
	cleanEncoedShortH3 := cleanString(encoedH3Int)
	if len(cleanEncoedShortH3) <= CODE_LENGTH {
		cleanEncoedShortH3 = xstrings.RightJustify(cleanEncoedShortH3, CODE_LENGTH, PADDING_CHAR)
	}
	parts := []string{}
	for i := 0; i < len(cleanEncoedShortH3); i += TUPLE_LENGTH {
		parts = append(parts, cleanEncoedShortH3[i:i+TUPLE_LENGTH])
	}
	return "@" + strings.Join(parts, "-")
}

func encodeShortInt(shortInt int64) string {
	if shortInt == 0 {
		return string(ALPHABET[0])
	}
	res := ""
	for shortInt > 0 {
		remainder := shortInt % ALPHABET_LENGTH
		res = string(ALPHABET[remainder]) + res
		shortInt = int64(shortInt / ALPHABET_LENGTH)
	}
	return res
}

func cleanString(s string) string {
	for k, v := range REPLACEMENT_MAP {
		s = strings.Replace(s, k, v, -1)
	}
	return s
}

func GeoToPlacekey(lat, long float64) (string, error) {
	if lat < -90 || lat > 90 || long < -180 || long > 180 {
		return "", errors.New("invalid lat/long range")
	}
	return encodeH3Int(int64(h3.FromGeo(
		h3.GeoCoord{Latitude: lat, Longitude: long},
		RESOLUTION,
	))), nil
}

func parsePlacekey(placekey string) (string, string, error) {
	what, where := "", ""
	if strings.Contains(placekey, "@") {
		parts := strings.Split(placekey, "@")
		if len(parts) != 2 {
			return what, where, errors.New("bad placekey")
		}
		what = parts[0]
		where = parts[1]
	} else {
		where = placekey
	}
	return what, where, nil
}

func stripEncoding(s string) string {
	s = strings.Replace(s, "@", "", -1)
	s = strings.Replace(s, "-", "", -1)
	s = strings.Replace(s, PADDING_CHAR, "", -1)
	return s
}

func dirtyString(s string) string {
	for k, v := range REPLACEMENT_MAP {
		s = strings.Replace(s, k, v, -1)
	}
	return s
}

func decodeString(s string) int64 {
	var val int64
	reversedS := xstrings.Reverse(s)
	for i := 0; i < len(s); i++ {
		targetTogetIndex := string(reversedS[i])
		val += int64(math.Pow(float64(ALPHABET_LENGTH), float64(i))) *
			int64(strings.Index(ALPHABET, targetTogetIndex))
	}
	return val
}

func decodeToH3Int(wherePart string) int64 {
	code := stripEncoding(wherePart)
	dirtyEncoding := dirtyString(code)
	shortedInt := decodeString(dirtyEncoding)
	return unshortenH3Integer(shortedInt)
}

// PlacekeyToH3 convert placekey to H3 Index
func PlacekeyToH3(placekey string) (*h3.H3Index, error) {
	_, where, err := parsePlacekey(placekey)
	if err != nil {
		return nil, err
	}
	h3Int := decodeToH3Int(where)
	idx := h3.FromString(strconv.FormatUint(uint64(h3Int), 16))
	return &idx, nil
}

// PlacekeyToGeo convert placekey to latitude,longitude
func PlacekeyToGeo(placekey string) (float64, float64, error) {
	idx, err := PlacekeyToH3(placekey)
	if err != nil {
		return 0, 0, err
	}
	h3Idx := h3.ToGeo(*idx)
	return h3Idx.Latitude, h3Idx.Longitude, nil
}

func PlacekeyDistance(pk1 string, pk2 string) (float64, error) {
	lat1, long1, err := PlacekeyToGeo(pk1)
	if err != nil {
		return 0, err
	}

	lat2, long2, err := PlacekeyToGeo(pk2)
	if err != nil {
		return 0, err
	}
	return GeoDistance(lat1, long1, lat2, long2), nil
}
