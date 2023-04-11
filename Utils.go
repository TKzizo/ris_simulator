package main

import (
	cmat "RIS_SIMULATOR/reducedComplex"
	"math"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distuv"
)

type Coordinates struct {
	x float64
	y float64
	z float64
}

func Distance(a, b Coordinates) float64 {
	return math.Sqrt(math.Pow(a.x-b.x, 2) + math.Pow(a.y-b.y, 2) + math.Pow(a.z-b.z, 2))
}
func RadToDeg(a float64) float64 {
	return a * 180 / math.Pi
}
func DegToRad(a float64) float64 {
	return a * math.Pi / 180
}
func sign(a, b float64) int8 {
	if a >= b {
		return 1
	}
	return -1
}
func Ge(theta float64) float64 {
	return Gain * math.Pow(math.Cos(theta), 2*q)
}
func L(s *Simulation, sf distuv.Normal, los bool, a ...Coordinates) float64 {
	if los == true {
		return math.Pow(10, (-20*math.Log10(4*math.Pi/s.Lambda)-10*s.n_LOS*(1+s.b_LOS*((s.Frequency-s.f0)/s.f0))*math.Log10(Distance(a[0], a[1]))-sf.Rand())/10)
	} else {
		return math.Pow(10, (-20*math.Log10(4*math.Pi/s.Lambda)-10*s.n_LOS*(1+s.b_LOS*((s.Frequency-s.f0)/s.f0))*math.Log10(Distance(a[0], a[1])+Distance(a[1], a[2]))-sf.Rand())/10)
	}
}

func Determine_Pb(a, b Coordinates) float64 {
	d := Distance(a, b)
	//LOS probability in Indoor office
	if d <= 1.2 {
		return 1
	} else if (1.2 < d) && (d <= 6.5) {
		return math.Exp((-d + 1.2) / 4.7)
	} else {
		return 0.32 * math.Exp((-d+6.5)/32.6)
	}
}

func cmat2cdense(matr cmat.Cmatrix) *mat.CDense {
	size := matr.Row * matr.Col
	data := make([]complex128, size)
	for x := 0; x < matr.Row; x++ {
		for y := 0; y < matr.Col; y++ {
			data[x*matr.Col+y] = matr.Data[x][y]
		}

	}
	return mat.NewCDense(matr.Row, matr.Col, data)
}
