package spiritofmetal

import (
	"reflect"
	"testing"
)

func TestSearch(t *testing.T) {
	for _, tt := range []struct {
		name string
		want []Band
	}{
		{
			name: "converge",
			want: []Band{
				Band{Name: "Converge", Genre: "Hardcore", Country: "USA"},
				Band{Name: "Convergence", Genre: "Melodic Death", Country: "Italy"},
				Band{Name: "Convergence", Genre: "Death Dark", Country: "Austria"},
				Band{Name: "Convergence From Within", Genre: "Death Metal", Country: "USA"},
			},
		},
		{
			name: "sdfgiousdfg",
			want: []Band{
				Band{Name: "Sidious", Genre: "Symphonic Death Black", Country: "United-Kingdom"},
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
