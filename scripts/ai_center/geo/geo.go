package geo

import (
	"errors"

	"github.com/sirupsen/logrus"
)

type Point2Df [2]float64
type MultiPointf []Point2D
type Polygonf []Point2D

type Point2D [2]int
type MultiPoint []Point2D
type Polygon []Point2D

func PositionsToPoints(p []Position) []Point2D {
	var ret []Point2D
	for _, v := range p {
		ret = append(ret, Point2D{int(v.X), int(v.Y)})
	}
	return ret
}

var errPolygonLength = errors.New("polygon length error")

func (p Polygon) ToPositionList() ([]Position, error) {
	if len(p) <= 3 { // 多边形至少3个点
		logrus.Errorf("Polygon ToPositionList failed: not a polygon, posList len[%v]", len(p))
		return nil, errPolygonLength
	}
	var posList []Position
	for _, v := range p { //坐标
		posList = append(posList, Position{
			X: float64(v[0]),
			Y: float64(v[1]),
		})
	}
	return posList, nil
}
