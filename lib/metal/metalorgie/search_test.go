package metalorgie

import (
	"reflect"
	"testing"

	"github.com/StalkR/goircbot/lib/metal"
)

func TestSearch(t *testing.T) {
	t.Skip() // cloudflare denies tests from github actions, works otherwise
	for _, tt := range []struct {
		name string
		want []metal.Band
	}{
		{
			name: "Conv",
			want: []metal.Band{
				metal.Band{Name: "Converge", Genre: "Hardcore / Mathcore / Metal / Noise / Punk", Country: "US"},
				metal.Band{Name: "Convict", Genre: "Pop / Punk / Punk Rock", Country: "BE"},
				metal.Band{Name: "Convulsing", Genre: "Black Metal / Death Metal", Country: "AU"},
			},
		},
		{
			name: "sdfgsdfg",
			want: nil,
		},
	} {
		got, err := Search(tt.name)
		if err != nil {
			t.Errorf("Search(%s): err: %v", tt.name, err)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Search(%s): got %s; want %s", tt.name, got, tt.want)
		}
	}
}
