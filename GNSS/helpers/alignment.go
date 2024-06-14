package helpers

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"gonum.org/v1/gonum/mat"
)

//  todo Azimuth-Elevation-Range Coordinates

type NedConverter struct {
	Ecef2NedMatrix *mat.Dense
}

// =====================================================

// =====================================================

func matMult(a, b [3][3]float64) [3][3]float64 {
	var result [3][3]float64

	result[0][0] = a[0][0]*b[0][0] + a[0][1]*b[1][0] + a[0][2]*b[2][0]
	result[0][1] = a[0][0]*b[0][1] + a[0][1]*b[1][1] + a[0][2]*b[2][1]
	result[0][2] = a[0][0]*b[0][2] + a[0][1]*b[1][2] + a[0][2]*b[2][2]

	result[1][0] = a[1][0]*b[0][0] + a[1][1]*b[1][0] + a[1][2]*b[2][0]
	result[1][1] = a[1][0]*b[0][1] + a[1][1]*b[1][1] + a[1][2]*b[2][1]
	result[1][2] = a[1][0]*b[0][2] + a[1][1]*b[1][2] + a[1][2]*b[2][2]

	result[2][0] = a[2][0]*b[0][0] + a[2][1]*b[1][0] + a[2][2]*b[2][0]
	result[2][1] = a[2][0]*b[0][1] + a[2][1]*b[1][1] + a[2][2]*b[2][1]
	result[2][2] = a[2][0]*b[0][2] + a[2][1]*b[1][2] + a[2][2]*b[2][2]

	return result
}

// =====================================================

// =====================================================

func rotMatrix(roll, pitch, yaw float64) [3][3]float64 {

	cr, sr := math.Cos(roll), math.Sin(roll)
	cp, sp := math.Cos(pitch), math.Sin(pitch)
	cy, sy := math.Cos(yaw), math.Sin(yaw)

	rr := [3][3]float64{
		{1, 0, 0},
		{0, cr, -sr},
		{0, sr, cr},
	}

	rp := [3][3]float64{
		{cp, 0, sp},
		{0, 1, 0},
		{-sp, 0, cp},
	}

	ry := [3][3]float64{
		{cy, -sy, 0},
		{sy, cy, 0},
		{0, 0, 1},
	}

	return matMult(ry, matMult(rp, rr))
}

// =====================================================

// =====================================================

func Rotate(axis []float64, angle float64) *mat.Dense {
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

// =====================================================

// =====================================================

func QuaterionProduct(q, r [4]float64) [4]float64 {
	var t [4]float64
	t[0] = r[0]*q[0] - r[1]*q[1] - r[2]*q[2] - r[3]*q[3]
	t[1] = r[0]*q[1] + r[1]*q[0] - r[2]*q[3] + r[3]*q[2]
	t[2] = r[0]*q[2] + r[1]*q[3] + r[2]*q[0] - r[3]*q[1]
	t[3] = r[0]*q[3] - r[1]*q[2] + r[2]*q[1] + r[3]*q[0]
	return t
}

func Quaternion2Euler(quats [][]float64) [][]float64 {
	eulers := make([][]float64, len(quats))

	for i, quat := range quats {
		q0, q1, q2, q3 := quat[0], quat[1], quat[2], quat[3]

		q0q1 := 2 * q0 * q1
		q0q2 := 2 * q0 * q2
		q0q3 := 2 * q0 * q3
		q1q1 := 2 * q1 * q1
		q1q2 := 2 * q1 * q2
		q2q2 := 2 * q2 * q2
		q2q3 := 2 * q2 * q3
		q3q3 := 2 * q3 * q3

		gamma := math.Atan2(q0q1+q2q3, 1-q1q1-q2q2)
		theta := math.Asin(q0q2 - q3*q1)
		psi := math.Atan2(q0q3+q1q2, 1-q2q2-q3q3)

		eulers[i] = []float64{gamma, theta, psi}
	}

	return eulers
}

func euler2Quaternion(eulerAngles [][]float64) [][]float64 {
	quaternions := make([][]float64, len(eulerAngles))
	for i := range quaternions {
		quaternions[i] = make([]float64, 4)
	}
	for i, angles := range eulerAngles {
		roll, pitch, yaw := angles[0], angles[1], angles[2]
		cosHalfRoll := math.Cos(roll / 2)
		cosHalfPitch := math.Cos(pitch / 2)
		cosHalfYaw := math.Cos(yaw / 2)
		sinHalfRoll := math.Sin(roll / 2)
		sinHalfPitch := math.Sin(pitch / 2)
		sinHalfYaw := math.Sin(yaw / 2)
		q0 := cosHalfRoll*cosHalfPitch*cosHalfYaw + sinHalfRoll*sinHalfPitch*sinHalfYaw
		q1 := sinHalfRoll*cosHalfPitch*cosHalfYaw - cosHalfRoll*sinHalfPitch*sinHalfYaw
		q2 := cosHalfRoll*sinHalfPitch*cosHalfYaw + sinHalfRoll*cosHalfPitch*sinHalfYaw
		q3 := cosHalfRoll*cosHalfPitch*sinHalfYaw - sinHalfRoll*sinHalfPitch*cosHalfYaw
		quaternions[i] = []float64{q0, q1, q2, q3}
		if quaternions[i][0] < 0 {
			for j := range quaternions[i] {
				quaternions[i][j] = -quaternions[i][j]
			}
		}
	}
	return quaternions
}

// TODO Euler2Quaternion
// Quaternion2Rot
// ROT2Quaternion
// Euler2Rot
// Rot2Euler
// ===================== Benchmarking =====================
func generateRandomMatrix() [3][3]float64 {
	var matrix [3][3]float64
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			matrix[i][j] = r.Float64()
		}
	}

	return matrix
}

func TestMatrixMultiplication1(iterations int) {
	start := time.Now()

	for i := 0; i < iterations; i++ {
		a := generateRandomMatrix()
		b := generateRandomMatrix()
		matMult(a, b)
	}

	duration := time.Since(start)
	fmt.Printf("Time taken for %d iterations: %v\n", iterations, duration)

}

func randomAngle() float64 {
	return (rand.Float64() * 2 * math.Pi) - math.Pi
}

func generateRandomAngles() (float64, float64, float64) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	roll := randomAngle()
	pitch := randomAngle()
	yaw := randomAngle()
	return roll, pitch, yaw
}

func TestMatrixRotation(iterations int) {
	start := time.Now()

	for i := 0; i < iterations; i++ {
		roll, pitch, yaw := generateRandomAngles()
		rotMatrix(roll, pitch, yaw)
	}

	duration := time.Since(start)
	fmt.Printf("Time taken for %d iterations: %v\n", iterations, duration)

}

func generateRandomQuaternion() [][]float64 {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	w := rand.Float64()*2 - 1
	x := rand.Float64()*2 - 1
	y := rand.Float64()*2 - 1
	z := rand.Float64()*2 - 1

	length := math.Sqrt(w*w + x*x + y*y + z*z)
	w /= length
	x /= length
	y /= length
	z /= length

	quaternion := [][]float64{
		{w, x, y, z},
	}

	return quaternion
}

func generateEuler() [][]float64 {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	roll := randomAngle()
	pitch := randomAngle()
	yaw := randomAngle()

	euler := [][]float64{
		{roll, pitch, yaw},
	}

	return euler

}

func BenchEuler2Quarterion(iterations int) {
	start := time.Now()

	for i := 0; i < iterations; i++ {
		euler := generateEuler()

		euler2Quaternion(euler)
	}

	duration := time.Since(start)
	fmt.Printf("Time taken for %d iterations: %v\n", iterations, duration)
}

func BenchQuaternion2Euler(iterations int) {
	start := time.Now()

	for i := 0; i < iterations; i++ {
		quat := generateRandomQuaternion()

		Quaternion2Euler(quat)
	}

	duration := time.Since(start)
	fmt.Printf("Time taken for %d iterations: %v\n", iterations, duration)
}

// ===================== Benchmarking =====================
