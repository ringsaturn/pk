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
	"regexp"
	"strconv"
	"strings"

	"github.com/huandu/xstrings"
	"github.com/uber/h3-go/v4"
)

const (
	// below const are copy from Python SDK and most of them can be found from
	// Placekey white paper pdf.
	_RESOLUTION        = 10
	_BASE_RESOLUTION   = 12
	_EARTH_RADIUS      = 6371                           // km, for distance computing
	_ALPHABET_BASE     = "23456789BCDFGHJKMNPQRSTVWXYZ" // will be lower case when ini
	_CODE_LENGTH       = 9
	_TUPLE_LENGTH      = 3
	_PADDING_CHAR      = "a"
	_REPLACEMENT_CHARS = "eu"

	// _HIGH_RESOLUTION_OFFSET used when Placekey convert back to H3 ID.
	// Or will failed when `h3.IsValid(xxx)`.
	// TODO: check why Go SDK need this but Python SDK doesn't.
	_HIGH_RESOLUTION_SHIFT = 255
)

var (
	_BASE_CELL_SHIFT          = int64(math.Pow(2, 45)) // Adding this will increment the base cell value by 1
	_UNUSED_RESOLUTION_FILLER = int64(math.Pow(2, (3*(15-_BASE_RESOLUTION))-1))
	_ALPHABET                 string         // build in init
	_ALPHABET_LENGTH          int64          // build in init
	_HEADER_BITS              string         // build in init
	_HEADER_INT               int64          // build in init
	_FIRST_TUPLE_REGEX        string         // build in init
	_TUPLE_REGEX              string         // build in init
	_WHERE_REGEX              *regexp.Regexp // build in init
	_WHAT_REGEX               *regexp.Regexp // build in init
	_WHAT_V2_REGEX            *regexp.Regexp // build in init

	_REPLACEMENT_MAP = map[string]string{
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
	_ALPHABET = strings.ToLower(_ALPHABET_BASE)
	_ALPHABET_LENGTH = int64(len(_ALPHABET))

	_FIRST_TUPLE_REGEX = "[" + _ALPHABET + _REPLACEMENT_CHARS + _PADDING_CHAR + "]{3}"
	_TUPLE_REGEX = "[" + _ALPHABET + _REPLACEMENT_CHARS + "]{3}"
	_WHERE_REGEX = regexp.MustCompile("^" + strings.Join([]string{_FIRST_TUPLE_REGEX, _TUPLE_REGEX, _TUPLE_REGEX}, "-") + "$")
	_WHAT_REGEX = regexp.MustCompile("^[" + _ALPHABET + "]{3,}(-[" + _ALPHABET + "]{3,})?$")

	_WHAT_V2_REGEX = regexp.MustCompile("^[01][abcdefghijklmnopqrstuvwxyz234567]{9}$")
}

// zfill like Python's zfill
func zfill(rawString string, padString string, expectLength int) string {
	diff := expectLength - len(rawString)
	if diff > 0 {
		rawString = strings.Repeat(padString, diff) + rawString
	}
	return rawString
}

func init_H3_HEADER() {
	idx := strconv.FormatInt(int64(h3.NewLatLng(0, 0).Cell(_RESOLUTION)), 2)
	// Python bin(xxx) result has prefix "0bxxxxx"
	// Golang FormatInt does not
	filled := zfill(idx, "0", 64)
	bits := filled[:12]
	if bits != "000010001010" {
		panic(errors.New(bits))
	}
	_HEADER_BITS = bits

	for index, char := range xstrings.Reverse(_HEADER_BITS) {
		intV, err := strconv.ParseInt(string(char), 10, 0)
		if err != nil {
			panic(err)
		}
		_HEADER_INT = _HEADER_INT + intV*int64(math.Pow(2, float64(index)))
	}
	_HEADER_INT = _HEADER_INT * int64(math.Pow(2, 52))
	if _HEADER_INT != 621496748577128448 {
		panic(_HEADER_INT)
	}
}

func radians(degree float64) float64 {
	return degree * math.Pi / 180
}

func geoDistance(lat1, lng1, lat2, lng2 float64) float64 {
	lat1 = radians(lat1)
	lng1 = radians(lng1)
	lat2 = radians(lat2)
	lng2 = radians(lng2)

	hav_lat := 0.5 * (1 - math.Cos(lat1-lat2))
	hav_lng := 0.5 * (1 - math.Cos(lng1-lng2))
	radical := math.Sqrt(hav_lat + math.Cos(lat1)*math.Cos(lat2)*hav_lng)
	return 2 * _EARTH_RADIUS * math.Asin(radical) * 1000
}

// shortenH3Integer shorten an H3 integer to only include location data up to the base resolution
func shortenH3Integer(h3Int int64) int64 {
	// Cuts off the 12 left-most bits that don't code location
	out := (h3Int + _BASE_CELL_SHIFT) % int64(math.Pow(2, 52))
	// Cuts off the rightmost bits corresponding to resolutions greater than the base resolution
	out = out >> (3 * (15 - _BASE_RESOLUTION))
	return out
}

func unshortenH3Integer(shortInt int64) int64 {
	unshifted_int := shortInt << (3 * (15 - _BASE_RESOLUTION))
	rebuilt_int := _HEADER_INT + _UNUSED_RESOLUTION_FILLER - _BASE_CELL_SHIFT + unshifted_int
	return rebuilt_int + _HIGH_RESOLUTION_SHIFT
}

func encodeH3Int(h3Int int64) string {
	shortH3Int := shortenH3Integer(h3Int)
	encoedH3Int := encodeShortInt(shortH3Int)
	cleanEncoedShortH3 := cleanString(encoedH3Int)
	if len(cleanEncoedShortH3) <= _CODE_LENGTH {
		cleanEncoedShortH3 = xstrings.RightJustify(cleanEncoedShortH3, _CODE_LENGTH, _PADDING_CHAR)
	}
	parts := []string{}
	for i := 0; i < len(cleanEncoedShortH3); i += _TUPLE_LENGTH {
		parts = append(parts, cleanEncoedShortH3[i:i+_TUPLE_LENGTH])
	}
	return "@" + strings.Join(parts, "-")
}

func encodeShortInt(shortInt int64) string {
	if shortInt == 0 {
		return string(_ALPHABET[0])
	}
	res := ""
	for shortInt > 0 {
		remainder := shortInt % _ALPHABET_LENGTH
		res = string(_ALPHABET[remainder]) + res
		shortInt = int64(shortInt / _ALPHABET_LENGTH)
	}
	return res
}

func cleanString(s string) string {
	for k, v := range _REPLACEMENT_MAP {
		s = strings.Replace(s, k, v, -1)
	}
	return s
}

func GeoToPlacekey(lat, lng float64) (string, error) {
	if lat < -90 || lat > 90 || lng < -180 || lng > 180 {
		return "", errors.New("invalid lat/lng range")
	}
	return encodeH3Int(int64(h3.NewLatLng(lat, lng).Cell(_RESOLUTION))), nil
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
	s = strings.Replace(s, _PADDING_CHAR, "", -1)
	return s
}

func dirtyString(s string) string {
	for k, v := range _REPLACEMENT_MAP {
		s = strings.Replace(s, k, v, -1)
	}
	return s
}

func decodeString(s string) int64 {
	var val int64
	reversedS := xstrings.Reverse(s)
	for i := 0; i < len(s); i++ {
		targetTogetIndex := string(reversedS[i])
		val += int64(math.Pow(float64(_ALPHABET_LENGTH), float64(i))) *
			int64(strings.Index(_ALPHABET, targetTogetIndex))
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
func PlacekeyToH3(placekey string) (*h3.Cell, error) {
	_, where, err := parsePlacekey(placekey)
	if err != nil {
		return nil, err
	}
	idx := h3.Cell(decodeToH3Int(where))
	return &idx, nil
}

// PlacekeyToGeo convert placekey to latitude,longitude
func PlacekeyToGeo(placekey string) (float64, float64, error) {
	idx, err := PlacekeyToH3(placekey)
	if err != nil {
		return 0, 0, err
	}
	lat := idx.LatLng().Lat
	lng := idx.LatLng().Lng
	return lat, lng, nil
}

func PlacekeyDistance(pk1 string, pk2 string) (float64, error) {
	lat1, lng1, err := PlacekeyToGeo(pk1)
	if err != nil {
		return 0, err
	}

	lat2, lng2, err := PlacekeyToGeo(pk2)
	if err != nil {
		return 0, err
	}
	return geoDistance(lat1, lng1, lat2, lng2), nil
}

func validateWhat(what string) bool {
	return _WHAT_REGEX.MatchString(what) || _WHAT_V2_REGEX.MatchString(what)
}

func validateWhere(where string) bool {
	h, err := PlacekeyToH3(where)
	if err != nil {
		return false
	}
	return _WHERE_REGEX.MatchString(where) && h.IsValid()
}

// ValidatePlacekey will use Regex and H3 to
// validate Placekey's what(if provided) and where part.
func ValidatePlacekey(pk string) bool {
	what, where, err := parsePlacekey(pk)
	if err != nil {
		return false
	}
	if what == "" {
		return validateWhere(where)
	}
	return validateWhat(what) && validateWhere(where)
}
