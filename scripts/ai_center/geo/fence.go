package geo

import (
	"errors"
	"math"
	"sort"

	"github.com/randomSignal/mathematics"
)

type Position struct {
	X float64 `json:"x,omitempty"`
	Y float64 `json:"y,omitempty"`
}

func (p Position) GetAngle() float64 {
	angle := math.Atan2(p.Y, p.X)
	// 把三四象限的角度翻正, 四象限的角度会成为最小
	if angle <= 0 {
		angle = math.Pi*2 + angle
	}
	return angle
}

/*
非常重要，这里的fence并处在我们所说的笛卡尔坐标系，而是一个y轴向下的坐标轴，和屏幕的坐标系一致
*/

func (p Position) Distance(p1 Position) float64 {
	return math.Sqrt(math.Pow(p.X-p1.X, 2) + math.Pow(p.Y-p1.Y, 2))
}

type positionByClockWise []Position

func (p positionByClockWise) Len() int {
	return len(p)
}

func (p positionByClockWise) Less(i, j int) bool {
	angle1 := p[i].GetAngle()
	angle2 := p[j].GetAngle()

	if p[i].X < 0 && p[j].X < 0 { // 二三象限
		if p[i].X != p[j].X {
			return p[i].X < p[j].X
		} else {
			return angle1 < angle2
		}
	} else if p[i].X > 0 && p[j].X > 0 { //一四象限
		if (p[i].Y > 0 && p[j].Y > 0) || (p[i].Y < 0 && p[j].Y < 0) { // 同在一象限或者是四象限
			if angle1 < angle2 {
				return true
			}
			if angle1 == angle2 && p[i].Distance(Position{0, 0}) < p[j].Distance(Position{0, 0}) {
				return true
			}
			return false
		} else { // 分别在一四象限
			return p[i].Y > p[j].Y
		}
	} else if p[i].X == 0 && p[j].X == 0 {
		if angle1 > angle2 {
			return true
		} else if angle1 == angle2 && p[i].Distance(Position{0, 0}) < p[j].Distance(Position{0, 0}) {
			return true
		}
		return false
	}
	return p[i].X < p[j].X
}

func (p positionByClockWise) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// 围栏
type Fence struct {
	// 电子围栏名称
	Name string `json:"name"`
	// 围栏的的点阵
	PositionList []Position `json:"positionList"`
	// 围栏的与其他矩形的相交比例
	IntersectionRatio float64 `json:"intersectionRatio,omitempty" default:"0.7"`
}

type FenceConfig struct {
	FenceList []Fence `json:"fenceList,omitempty"`
}

func (fence Fence) MinPoint() Position {
	var (
		minX, minY float64
	)

	minX = fence.PositionList[0].X
	minY = fence.PositionList[0].Y

	for _, v := range fence.PositionList {
		if minX > v.X {
			minX = v.X
		}

		if minY > v.Y {
			minY = v.Y
		}
	}

	return Position{X: minY, Y: minY}
}

func (fence Fence) MaxPoint() Position {
	var (
		maxX, maxY float64
	)

	for _, v := range fence.PositionList {
		if maxX < v.X {
			maxX = v.X
		}

		if maxY < v.Y {
			maxY = v.Y
		}
	}

	return Position{X: maxX, Y: maxY}
}

// 矩阵交叠
func FenceIntrusion(fenceConfig FenceConfig, positionList [][]Position) ([][]Position, error) {
	var res [][]Position
	resOffset := map[int]int{}
	for _, fence := range fenceConfig.FenceList {
		if len(fence.PositionList) != 4 {
			return nil, errors.New("only support matrix")
		}
		fenceMatrix := mathematics.Matrix{
			PointList: [4]mathematics.Point{
				{X: fence.PositionList[0].X, Y: fence.PositionList[0].Y},
				{X: fence.PositionList[1].X, Y: fence.PositionList[1].Y},
				{X: fence.PositionList[2].X, Y: fence.PositionList[2].Y},
				{X: fence.PositionList[3].X, Y: fence.PositionList[3].Y},
			},
		}
		fenceMatrix = fenceMatrix.Correction()
		for i := 0; i < len(positionList); i++ {
			if _, ok := resOffset[i]; ok {
				continue
			}
			position := positionList[i]
			if len(position) != 4 {
				return nil, errors.New("only support matrix")
			}
			argsMatrix := mathematics.Matrix{
				PointList: [4]mathematics.Point{
					{X: position[0].X, Y: position[0].Y},
					{X: position[1].X, Y: position[1].Y},
					{X: position[2].X, Y: position[2].Y},
					{X: position[3].X, Y: position[3].Y},
				},
			}
			if fenceMatrix.IntersectionArea(argsMatrix)/argsMatrix.Area() > fence.IntersectionRatio {
				resOffset[i] = i
				res = append(res, position)
			}
		}
	}
	return res, nil
}

func (fence Fence) LeftPointIndex() int {
	// 最左边的点的坐标
	var (
		lp = Position{
			X: math.MaxFloat64,
			Y: math.MaxFloat64,
		}
		index = 0
	)

	for idx, pos := range fence.PositionList {
		if pos.X < lp.X {
			lp = pos
			index = idx
		}
		if pos.X == lp.X {
			if pos.Y < lp.Y {
				lp = pos
				index = idx
			}
		}
	}
	return index
}

func (fence *Fence) ClockWiseSpin() {

	// 使最左边的点作为起始点，顺时针旋转
	var center = Position{
		X: 0,
		Y: 0,
	}
	for _, pos := range fence.PositionList {
		center.X += pos.X
		center.Y += pos.Y
	}
	center.X, center.Y = center.X/float64(len(fence.PositionList)), center.Y/float64(len(fence.PositionList))
	// 将中心作为坐标轴的中心 建立新的坐标系
	for i := 0; i < len(fence.PositionList); i++ {
		fence.PositionList[i].X = fence.PositionList[i].X - center.X
		fence.PositionList[i].Y = fence.PositionList[i].Y - center.Y
	}
	sort.Sort(positionByClockWise(fence.PositionList))
	//回到原先的坐标系
	for i := 0; i < len(fence.PositionList); i++ {
		fence.PositionList[i].X = fence.PositionList[i].X + center.X
		fence.PositionList[i].Y = fence.PositionList[i].Y + center.Y
	}
}

func (fence Fence) IsClockWise() bool {
	// 是否顺时针旋转
	var center = Position{
		X: 0,
		Y: 0,
	}
	for _, pos := range fence.PositionList {
		center.X += pos.X
		center.Y += pos.Y
	}
	center.X, center.Y = center.X/float64(len(fence.PositionList)), center.Y/float64(len(fence.PositionList))
	// 将中心作为坐标轴的中心 建立新的坐标系
	for i := 0; i < len(fence.PositionList); i++ {
		fence.PositionList[i].X = fence.PositionList[i].X - center.X
		fence.PositionList[i].Y = fence.PositionList[i].Y - center.Y
	}
	res := sort.IsSorted(positionByClockWise(fence.PositionList))
	//回到原先的坐标系
	for i := 0; i < len(fence.PositionList); i++ {
		fence.PositionList[i].X = fence.PositionList[i].X + center.X
		fence.PositionList[i].Y = fence.PositionList[i].Y + center.Y
	}
	return res
}

func (fence Fence) IsRectangle() bool {
	// 是否矩形
	if len(fence.PositionList) != 4 {
		return false
	}

	var center = Position{
		X: 0,
		Y: 0,
	}
	for _, pos := range fence.PositionList {
		center.X += pos.X
		center.Y += pos.Y
	}
	center.X, center.Y = center.X/4, center.Y/4
	d1 := center.Distance(fence.PositionList[0])
	for i := 1; i < len(fence.PositionList); i++ {
		if d1 != center.Distance(fence.PositionList[i]) {
			return false
		}
	}
	return true
}

func (fence Fence) BoundingRectangle() Fence {
	// 不规则fence的外接矩形fence
	var minX, minY, maxX, maxY float64 = math.MaxFloat64, math.MaxFloat64, 0, 0
	for _, pos := range fence.PositionList {
		if pos.X < minX {
			minX = pos.X
		}
		if pos.Y < minY {
			minY = pos.Y
		}
		if pos.X > maxX {
			maxX = pos.X
		}
		if pos.Y > maxY {
			maxY = pos.Y
		}
	}
	return Fence{
		PositionList: []Position{
			{
				X: minX,
				Y: minY,
			},
			{
				X: maxX,
				Y: minY,
			},
			{
				X: maxX,
				Y: maxY,
			},
			{
				X: minX,
				Y: maxY,
			},
		},
	}
}

func (fence Fence) IsRectangleCW() bool {
	// 矩形是否顺时针
	var center = Position{
		X: 0,
		Y: 0,
	}
	for _, pos := range fence.PositionList {
		center.X += pos.X
		center.Y += pos.Y
	}
	center.X, center.Y = center.X/4, center.Y/4
	p1Value := fence.PositionList[0].X < center.X && fence.PositionList[0].Y < center.X
	p2Value := fence.PositionList[1].X > center.X && fence.PositionList[1].Y < center.X
	p3Value := fence.PositionList[2].X > center.X && fence.PositionList[2].Y > center.X
	p4Value := fence.PositionList[3].X < center.X && fence.PositionList[3].Y > center.X
	return p1Value && p2Value && p3Value && p4Value
}

func (fence *Fence) RectCWSpin() {
	// 矩形变为顺时针旋转
	var center = Position{
		X: 0,
		Y: 0,
	}
	for _, pos := range fence.PositionList {
		center.X += pos.X
		center.Y += pos.Y
	}
	center.X, center.Y = center.X/4, center.Y/4
	for j := 0; j < 2; j++ {
		for i := 0; i < len(fence.PositionList); i++ {
			if fence.PositionList[i].X < center.X && fence.PositionList[i].Y < center.Y && i != 0 {
				positionByClockWise(fence.PositionList).Swap(i, 0)
				continue
			}
			if fence.PositionList[i].X > center.X && fence.PositionList[i].Y < center.Y && i != 1 {
				positionByClockWise(fence.PositionList).Swap(i, 1)
				continue
			}
			if fence.PositionList[i].X > center.X && fence.PositionList[i].Y > center.Y && i != 2 {
				positionByClockWise(fence.PositionList).Swap(i, 2)
				continue
			}
			if fence.PositionList[i].X < center.X && fence.PositionList[i].Y > center.Y && i != 3 {
				positionByClockWise(fence.PositionList).Swap(i, 3)
				continue
			}
		}
	}
}

func (fence *Fence) RecTrans2NewCoorSys(fence1 Fence) (*Fence, error) {
	// 将fence1作为相对矩形位置进行坐标轴变换成一个新的fence
	if !fence.IsRectangle() {
		return nil, errors.New("error: this method only used in rectangle")
	}
	if !fence.IsRectangleCW() {
		fence.RectCWSpin()
	}
	if !fence1.IsRectangle() && !fence1.IsRectangleCW() {
		return nil, errors.New("error: fence1 must be a clockwise rectangle")
	}
	resFence := Fence{
		IntersectionRatio: fence.IntersectionRatio,
		PositionList:      []Position{},
	}
	basePointIndex := fence1.LeftPointIndex()
	for i := 0; i < len(fence1.PositionList); i++ {
		resFence.PositionList = append(resFence.PositionList, Position{
			X: fence1.PositionList[basePointIndex].X + fence.PositionList[i].X,
			Y: fence1.PositionList[basePointIndex].Y + fence.PositionList[i].Y,
		})
	}
	return &resFence, nil
}

func (fence *Fence) RecTrans2CoorSys(fence1 Fence) error {
	// 将一个fence1作为相对矩形位置进行坐标轴变换
	if !fence.IsRectangle() {
		return errors.New("error: this method only used in rectangle")
	}
	if !fence.IsRectangleCW() {
		fence.RectCWSpin()
	}
	if !fence1.IsRectangle() && !fence1.IsRectangleCW() {
		return errors.New("error: fence1 must be a clockwise rectangle")
	}
	basePointIndex := fence.LeftPointIndex()
	for i := 0; i < len(fence1.PositionList); i++ {
		fence.PositionList[i].X = fence.PositionList[i].X + fence1.PositionList[basePointIndex].X
		fence.PositionList[i].Y = fence.PositionList[i].Y + fence1.PositionList[basePointIndex].Y

	}
	return nil
}
