package spiritofmetal

import (
	"reflect"
	"testing"

	"github.com/StalkR/goircbot/lib/metal/band"
)

func TestSearch(t *testing.T) {
	for _, tt := range []struct {
		name string
		want []band.Band
	}{
		{
			name: "converge",
			want: []band.Band{
				band.Band{Name: "Converge", Genre: "Hardcore", Country: "USA"},
				band.Band{Name: "Convergence", Genre: "Melodic Death", Country: "Italy"},
				band.Band{Name: "Convergence", Genre: "Death Dark", Country: "Austria"},
				band.Band{Name: "Convergence From Within", Genre: "Death Metal", Country: "USA"},
			},
		},
		{
			name: "sdfgiousdfg",
			want: []band.Band{
				band.Band{Name: "Sidious", Genre: "Symphonic Death Black", Country: "United-Kingdom"},
			},
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
