package metalarchives

import (
	"reflect"
	"testing"

	"github.com/StalkR/goircbot/lib/metal"
)

func TestSearch(t *testing.T) {
	t.Skip() // metalarchives enabled cloudflare :(
	for _, tt := range []struct {
		name string
		want []metal.Band
	}{
		{
			name: "Convergence",
			want: []metal.Band{
				metal.Band{Name: "Convergence", Genre: "Atmospheric/Industrial Death Metal", Country: "Austria"},
				metal.Band{Name: "Convergence", Genre: "Melodic Death Metal (early); Nu-Metal/Alternative Rock (later)", Country: "Italy"},
				metal.Band{Name: "Convergence", Genre: "Death Metal", Country: "United States"},
				metal.Band{Name: "Convergence from Within", Genre: "Death Metal", Country: "United States"},
				metal.Band{Name: "Theory of Convergence (a.k.a. TOC)", Genre: "Progressive Metal/Rock", Country: "China"},
				metal.Band{Name: "Converg3nce (a.k.a. Convergence, converg3nce.)", Genre: "Atmospheric Black Metal", Country: "Brazil"},
			},
		},
		{
			name: "Psycroptic",
			want: []metal.Band{
				metal.Band{Name: "Psycroptic", Genre: "Technical Death Metal", Country: "Australia"},
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
