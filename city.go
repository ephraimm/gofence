package geofence

import (
	"errors"
	"os"

	"github.com/buckhx/diglet/geo"
)

type CityFence struct {
	features map[string][]*geo.Feature
	boros    []*geo.Feature
}

// Only for demonstrative purposes
// Checks the containing city first for inclusion, then features. Fully inspects each geometry in containing city
// This requires the NYC_BOROS_PATH envvar to be set to the Borrough Boundaries geojson file
// It can be found here http://www1.nyc.gov/site/planning/data-maps/open-data/districts-download-metadata.page
func NewCityFence() (fence *CityFence, err error) {
	path := os.Getenv("NYC_BOROS_PATH")
	if path == "" {
		err = errors.New("Missing NYC_BOROS_PATH envvar for CityFence")
		return
	}
	bfeatures, err := geo.NewGeojsonSource(path, nil).Publish()
	if err != nil {
		return
	}
	var boros []*geo.Feature
	for b := range bfeatures {
		boros = append(boros, b)
	}
	fence = &CityFence{
		boros:    boros,
		features: make(map[string][]*geo.Feature, 5),
	}
	return
}

// Features must contain a tag BoroName to match to a burrough
func (u *CityFence) Add(f *geo.Feature) {
	u.features[f.Tags("BoroName")] = append(u.features[f.Tags("BoroName")], f)
}

func (u *CityFence) Get(c geo.Coordinate) []*geo.Feature {
	var bn string
	for _, boro := range u.boros {
		if boro.Contains(c) {
			bn = boro.Tags("BoroName")
			break
		}
	}
	if bn == "" {
		return nil
	}
	var ins []*geo.Feature
	for _, f := range u.features[bn] {
		for _, shp := range f.Geometry {
			if shp.Contains(c) {
				ins = append(ins, f)
			}
		}
	}
	return ins
}
