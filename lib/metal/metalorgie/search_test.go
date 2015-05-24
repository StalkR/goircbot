package metalorgie

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
			name: "Conv",
			want: []band.Band{
				band.Band{Name: "Converge", Genre: "Hardcore Chaotique / Punk / Metal", Country: "USA"},
				band.Band{Name: "Convict", Genre: "Punk Rock / Pop Punk", Country: "Belgique"},
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
