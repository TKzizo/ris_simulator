package main

import (
	"math"
	"math/cmplx"

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
func L(s *Simulation, sf distuv.Normal, a, b Coordinates) float64 {
	return math.Pow(10, (-20*math.Log10(4*math.Pi/s.Lambda)-10*s.n*math.Log10(Distance(a, b))-sf.Rand())/10)
}
func Array_Response_RIS_Tx(s *Simulation, r *RIS, t *Tx_Rx) []complex128 {
	var vec []complex128
	for x := 0; x < int(math.Sqrt(float64(r.N))); x++ {
		for y := 0; y < int(math.Sqrt(float64(r.N))); y++ {
			argument := s.k * r.dis * (float64(x)*math.Sin(r.Theta_Tx) + float64(y)*math.Sin(r.Phi_Tx)*math.Cos(r.Theta_Tx))
			//fmt.Println("argument[", x, y, "]", argument)
			expo := cmplx.Exp(1i * complex(argument, 0))
			//fmt.Println("Expo array_response: ", expo)
			vec = append(vec, expo)
		}
	}
	//	fmt.Println(vec)
	return vec
}
func Array_Response_RIS_Rx(s *Simulation, r *RIS, u *Tx_Rx) []complex128 {
	var vec []complex128
	for x := 0; x < int(math.Sqrt(float64(r.N))); x++ {
		for y := 0; y < int(math.Sqrt(float64(r.N))); y++ {
			argument := s.k * r.dis * (float64(x)*math.Sin(r.Theta_Rx) + float64(y)*math.Sin(r.Phi_Rx)*math.Cos(r.Theta_Tx))
			//fmt.Println("argument[", x, y, "]", argument)
			expo := cmplx.Exp(1i * complex(argument, 0))
			//fmt.Println("Expo array_response: ", expo)
			vec = append(vec, expo)
		}
	}
	//	fmt.Println(vec)
	return vec
}
func Determine_Pb(a, b Coordinates) float64 {
	d := Distance(a, b)
	if d <= 1.2 {
		return 1
	} else if (1.2 < d) && (d <= 6.5) {
		return math.Exp((-d + 1.2) / 4.7)
	} else {
		return 0.32 * math.Exp((-d+6.5)/32.6)
	}
}
