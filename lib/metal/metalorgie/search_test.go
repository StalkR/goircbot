package metalorgie

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
			name: "Conv",
			want: []metal.Band{
				metal.Band{Name: "Converge", Genre: "Hardcore Chaotique / Punk / Metal", Country: "USA"},
				metal.Band{Name: "Convict", Genre: "Punk Rock / Pop Punk", Country: "Belgique"},
				metal.Band{Name: "Convulsing", Genre: "Death Metal / Black Metal", Country: "Australie"},
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
