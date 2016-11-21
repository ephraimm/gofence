package geofence

import (
	"github.com/buckhx/diglet/geo"
)

// Use an quadtree for segmention. Only inspect geometries whose quad tree leaf contain the query
type Qfence struct {
	zoom  int
	qtree map[string][]*geo.Feature
}

func NewQfence(zoom int) *Qfence {
	return &Qfence{
		zoom:  zoom,
		qtree: make(map[string][]*geo.Feature), //TODO hint size
	}
}

func (q *Qfence) Add(f *geo.Feature) {
	for _, shp := range f.Geometry {
		qkeys := shapeQkeys(shp, q.zoom)
		for _, key := range qkeys {
			q.qtree[key] = append(q.qtree[key], f)
		}
	}
}

func (q *Qfence) Get(c geo.Coordinate) (matchs []*geo.Feature) {
	key := c.ToTile(q.zoom).QuadKey()
	for _, f := range q.qtree[key] {
		if f.Contains(c) {
			matchs = append(matchs, f)
		}
	}
	return
}

func shapeQkeys(shp *geo.Shape, zoom int) (keys []string) {
	//TODO move this to geo package shp.QuadKeys
	// bbox could be trimmed
	bbox := shp.BoundingBox()
	ne := bbox.NorthEast().ToTile(zoom)
	sw := bbox.SouthWest().ToTile(zoom)
	cur := sw
	for x := sw.X; x <= ne.X; x++ {
		for y := sw.Y; y >= ne.Y; y-- { //origin is NW
			cur.X, cur.Y = x, y
			keys = append(keys, cur.QuadKey())
		}
	}
	return
}
