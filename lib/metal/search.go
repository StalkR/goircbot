// Package metal is a library to search for metal bands on multiple sites.
package metal

import (
	"github.com/StalkR/goircbot/lib/metal/band"
	"github.com/StalkR/goircbot/lib/metal/metalarchives"
	"github.com/StalkR/goircbot/lib/metal/metalorgie"
	"github.com/StalkR/goircbot/lib/metal/spiritofmetal"
)

// Search finds bands by name.
func Search(name string) ([]band.Band, error) {
	bands, err := spiritofmetal.Search(name)
	if err != nil {
		return nil, err
	}
	if len(bands) > 0 {
		return bands, nil
	}
	bands, err = metalarchives.Search(name)
	if err != nil {
		return nil, err
	}
	if len(bands) > 0 {
		return bands, nil
	}
	bands, err = metalorgie.Search(name)
	if err != nil {
		return nil, err
	}
	return bands, nil
}
