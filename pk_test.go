package pk_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/ringsaturn/pk"
)

func TestGeoToPlacekey(t *testing.T) {
	type args struct {
		lat  float64
		long float64
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "0,0",
			args: args{
				long: 0,
				lat:  0,
			},
			want:    "@dvt-smp-tvz",
			wantErr: false,
		},
		{
			name: "New York",
			args: args{
				long: -74.006058,
				lat:  40.712772,
			},
			want:    "@627-wbz-tjv",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pk.GeoToPlacekey(tt.args.lat, tt.args.long)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeoToPlacekey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GeoToPlacekey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= 0.001
}

func TestPlacekeyToGeo(t *testing.T) {
	type args struct {
		placekey string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		want1   float64
		wantErr bool
	}{
		{
			name: "0,0",
			args: args{
				placekey: "@dvt-smp-tvz",
			},
			want:    0,
			want1:   0,
			wantErr: false,
		},
		{
			name: "New York",
			args: args{
				placekey: "@627-wbz-tjv",
			},
			want:    40.712772,
			want1:   -74.006058,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := pk.PlacekeyToGeo(tt.args.placekey)
			if (err != nil) != tt.wantErr {
				t.Errorf("PlacekeyToGeo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !almostEqual(got, tt.want) {
				t.Errorf("PlacekeyToGeo() got = %v, want %v", got, tt.want)
			}
			if !almostEqual(got1, tt.want1) {
				t.Errorf("PlacekeyToGeo() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func BenchmarkGeoToPlacekey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = pk.GeoToPlacekey(39.9289, 116.3883)
	}
}

func BenchmarkPlacekeyToGeo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _ = pk.PlacekeyToGeo("@6qk-v3d-brk")
	}
}

func ExampleGeoToPlacekey() {
	k, _ := pk.GeoToPlacekey(39.9289, 116.3883)
	fmt.Println(k)
	// Output: @6qk-v3d-brk
}

func ExamplePlacekeyToGeo() {
	lat, long, _ := pk.PlacekeyToGeo("@6qk-v3d-brk")
	fmt.Printf("%.3f %.3f \n", lat, long)
	// Output: 39.929 116.388
}

func ExamplePlacekeyDistance() {
	dist, _ := pk.PlacekeyDistance("@qjk-m7r-whq", "@hvb-5d7-92k")
	fmt.Printf("%.1f\n", dist/1000)
	// Output: 13597.5
}

func TestValidatePlacekey(t *testing.T) {
	type args struct {
		pk string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "NY",
			args: args{
				pk: "@627-wbz-tjv",
			},
			want: true,
		},

		{
			name: "good what part",
			args: args{
				pk: "226@5vg-7gq-5mk",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pk.ValidatePlacekey(tt.args.pk); got != tt.want {
				t.Errorf("ValidatePlacekey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleValidatePlacekey() {
	fmt.Println(pk.ValidatePlacekey("@627-wbz-tjv"))
	// Output: true
}
