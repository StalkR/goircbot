package spiritofmetal

import (
	"reflect"
	"testing"

	"github.com/StalkR/goircbot/lib/metal"
)

func TestSearch(t *testing.T) {
	for _, tt := range []struct {
		name string
		want []metal.Band
	}{
		{
			name: "converge",
			want: []metal.Band{
				metal.Band{Name: "Converge", Genre: "Hardcore", Country: "USA"},
				metal.Band{Name: "Convergence", Genre: "Death Dark", Country: "Austria"},
				metal.Band{Name: "Convergence", Genre: "Melodic Death", Country: "Italy"},
				metal.Band{Name: "Convergence", Genre: "Death Metal", Country: "USA"},
				metal.Band{Name: "Convergence From Within", Genre: "Death Metal", Country: "USA"},
				metal.Band{Name: "S.U.C.", Genre: "Death Grind", Country: "Brazil"},
			},
		},
		{
			name: "sdfgiousdfg",
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
