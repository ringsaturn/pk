package pk

import (
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
