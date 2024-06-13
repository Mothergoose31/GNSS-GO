package helpers

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"
)

//  todo Azimuth-Elevation-Range Coordinates

func rotate(axis []float64, angle float64) *mat.Dense {
	cosAngle := math.Cos(angle)
	sinAngle := math.Sin(angle)
	oneMinusCos := 1 - cosAngle

	ret1 := mat.NewDense(3, 3, []float64{
		oneMinusCos * axis[0] * axis[0], oneMinusCos * axis[0] * axis[1], oneMinusCos * axis[0] * axis[2],
		oneMinusCos * axis[1] * axis[0], oneMinusCos * axis[1] * axis[1], oneMinusCos * axis[1] * axis[2],
		oneMinusCos * axis[2] * axis[0], oneMinusCos * axis[2] * axis[1], oneMinusCos * axis[2] * axis[2],
	})

	ret2 := mat.NewDense(3, 3, []float64{
		cosAngle, 0, 0,
		0, cosAngle, 0,
		0, 0, cosAngle,
	})

	ret3 := mat.NewDense(3, 3, []float64{
		0, -sinAngle * axis[2], sinAngle * axis[1],
		sinAngle * axis[2], 0, -sinAngle * axis[0],
		-sinAngle * axis[1], sinAngle * axis[0], 0,
	})

	var result mat.Dense
	result.Add(ret1, ret2)
	result.Add(&result, ret3)
	fmt.Println(result.Dims())

	return &result
}
