package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ringsaturn/pk"
)

const (
	FromGeoFlag = "from-geo"
	ToGeoFlag   = "to-geo"
)

func showErr() {
	fmt.Printf("expected '%v' or '%v' subcommands\n", FromGeoFlag, ToGeoFlag)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		showErr()
		return
	}

	switch os.Args[1] {
	case FromGeoFlag:
		FromGeoCmd := flag.NewFlagSet(FromGeoFlag, flag.ExitOnError)
		lat := FromGeoCmd.Float64("lat", 0, "latitude")
		lng := FromGeoCmd.Float64("lng", 0, "longitude")
		if err := FromGeoCmd.Parse(os.Args[2:]); err != nil {
			panic(err)
		}
		placeKey, err := pk.GeoToPlacekey(*lat, *lng)
		if err != nil {
			panic(err)
		}
		fmt.Println(placeKey)
	case ToGeoFlag:
		ToGeoCmd := flag.NewFlagSet(ToGeoFlag, flag.ExitOnError)
		placekey := ToGeoCmd.String("pk", "", "the place key need to convert to geo")
		if err := ToGeoCmd.Parse(os.Args[2:]); err != nil {
			panic(err)
		}
		lat, long, err := pk.PlacekeyToGeo(*placekey)
		if err != nil {
			panic(err)
		}
		fmt.Println(lat, long)
	default:
		showErr()
	}
}
