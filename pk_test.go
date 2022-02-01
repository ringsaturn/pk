package pk

import (
	"math"
	"testing"
)

func TestGeoToPlacekey(t *testing.T) {
	type args struct {
		lat  float64
		long float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "0,0",
			args: args{
				long: 0,
				lat:  0,
			},
			want: "@dvt-smp-tvz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GeoToPlacekey(tt.args.lat, tt.args.long); got != tt.want {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := PlacekeyToGeo(tt.args.placekey)
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
