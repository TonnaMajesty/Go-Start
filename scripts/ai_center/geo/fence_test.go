package geo

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestFenceIntrusion(t *testing.T) {
	var fenceConfig FenceConfig
	var positionList [][]Position

	/*
		dataStr := `[{"x":771,"y":272},{"x":932,"y":272},{"x":932,"y":588},{"x":771,"y":588}]`
		 fenceStr := `[{"x":478.1679389312977,"y":195.1145038167939},{"x":785.9541984732825,"y":195.1145038167939},{"x":785.9541984732825,"y":950.8396946564886},{"x":478.1679389312977,"y":950.8396946564886}]`
	*/
	fenceConfig.FenceList = append(fenceConfig.FenceList,
		Fence{
			IntersectionRatio: 0.1,
			PositionList: []Position{
				{771, 272},
				{932, 272},
				{932, 588},
				{771, 588},
			},
		},
		Fence{
			IntersectionRatio: 0.3,
			PositionList: []Position{
				{478, 195},
				{785, 195},
				{785, 950},
				{478, 950},
			},
		},
	)
	positionList = append(positionList, []Position{
		{15, 15},
		{25, 15},
		{25, 25},
		{15, 25},
	})

	res, err := FenceIntrusion(fenceConfig, positionList)
	require.NoError(t, err)
	spew.Dump(res)
}

func TestFenceRectangle(t *testing.T) {
	rectFence1 := Fence{
		PositionList: []Position{
			{680, 680},
			{870, 680},
			{870, 870},
			{680, 870},
		},
	}
	rectFence2 := Fence{
		PositionList: []Position{
			{20, 20},
			{30, 20},
			{30, 30},
			{20, 30},
		},
	}
	rectFence3 := Fence{
		PositionList: []Position{
			{700, 700},
			{710, 700},
			{710, 710},
			{700, 710},
		},
	}
	fence2 := Fence{
		PositionList: []Position{
			{1, 2},
			{3, 4},
			{5, 6},
			{7, 8},
		},
	}
	t.Run("是否是矩形", func(t *testing.T) {
		require.Equal(t, true, rectFence1.IsRectangle())
	})
	t.Run("矩形是否顺时针", func(t *testing.T) {
		require.Equal(t, true, rectFence1.IsRectangleCW())
	})
	t.Run("校验多边形的外接矩形是否正确", func(t *testing.T) {
		bindBox := fence2.BoundingRectangle()
		require.Equal(t, true, bindBox.IsRectangle() && bindBox.IsRectangleCW())
	})
	t.Run("矩形变化为顺时针", func(t *testing.T) {
		// 打乱顺序
		positionByClockWise(rectFence1.PositionList).Swap(0, 1)
		positionByClockWise(rectFence1.PositionList).Swap(0, 2)
		positionByClockWise(rectFence1.PositionList).Swap(3, 1)
		rectFence1.RectCWSpin() // 矩形顺时针变换
		require.Equal(t, true, rectFence1.IsRectangleCW())
	})
	t.Run("两个矩形相对变换", func(t *testing.T) {
		resFence, err := rectFence2.RecTrans2NewCoorSys(rectFence1)
		require.NoError(t, err)
		require.Equal(t, true, reflect.DeepEqual(resFence.PositionList, rectFence3.PositionList))
		err = rectFence2.RecTrans2CoorSys(rectFence1)
		require.NoError(t, err)
		require.Equal(t, true, reflect.DeepEqual(rectFence2.PositionList, rectFence3.PositionList))
	})
	t.Run("找最左边的点", func(t *testing.T) {
		positionByClockWise(rectFence1.PositionList).Swap(0, 1)
		require.Equal(t, 1, rectFence1.LeftPointIndex())
	})
}
