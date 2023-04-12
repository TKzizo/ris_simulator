package main

import (
	cmat "RIS_SIMULATOR/reducedComplex"
	"math"
	"math/cmplx"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

func (s *Simulation) G_channel() cmat.Cmatrix {

	s.Ris.Theta_Rx = float64(sign(s.Rx.xyz.z, s.Ris.xyz.z)) * math.Asin(math.Abs(s.Rx.xyz.z-s.Ris.xyz.z)/Distance(s.Ris.xyz, s.Rx.xyz))
	//s.Rx.Theta_RIS = float64(sign(s.Ris.xyz.z, s.Rx.xyz.z)) * math.Asin(math.Abs(s.Rx.xyz.z-s.Ris.xyz.z)/Distance(s.Ris.xyz, s.Rx.xyz))

	if s.Broadside == 0 { //Side Wall
		s.Ris.Phi_Rx = float64(sign(s.Ris.xyz.x, s.Rx.xyz.x)) * math.Atan2(math.Abs(s.Rx.xyz.x-s.Ris.xyz.x), math.Abs(s.Rx.xyz.y-s.Ris.xyz.y))

		s.Rx.Phi_RIS = math.Log(rand.Float64()/rand.Float64())*math.Sqrt(25/2) + rand.Float64()*180 - 90
		//s.Rx.Phi_RIS = float64(sign(s.Ris.xyz.y, s.Rx.xyz.y)) * math.Atan2(math.Abs(s.Rx.xyz.y-s.Ris.xyz.y), math.Abs(s.Rx.xyz.x-s.Ris.xyz.x))
	} else if s.Broadside == 1 { // Opposite wall
		s.Ris.Phi_Rx = float64(sign(s.Rx.xyz.y, s.Ris.xyz.y)) * math.Atan2(math.Abs(s.Rx.xyz.y), math.Abs(s.Rx.xyz.y))
		//s.Rx.Phi_RIS = float64(sign(s.Rx.xyz.x, s.Ris.xyz.y)) * math.Atan2(math.Abs(s.Rx.xyz.y-s.Ris.xyz.z), math.Abs(s.Rx.xyz.x-s.Ris.xyz.x))
		s.Rx.Phi_RIS = math.Log(rand.Float64()/rand.Float64())*math.Sqrt(25/2) + rand.Float64()*180 - 90

	}

	AR_rx_ris := cmat.Cmatrix{}
	AR_rx_ris.Init(s.Ris.N, s.Tx.N)

	dx := int(math.Sqrt(float64(s.Ris.N)))
	dy := dx

	AR_ris := make([]complex128, s.Ris.N, s.Ris.N)
	AR_rx := make([]complex128, s.Rx.N, s.Rx.N)

	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			AR_ris[x*dx+y] = cmplx.Exp(
				1i * complex(s.k*s.Ris.dis*(float64(x)*math.Sin(s.Ris.Theta_Rx)+float64(y)*math.Sin(s.Ris.Phi_Rx)*math.Cos(s.Ris.Theta_Rx)), 0))
		}
	}

	if s.Tx.Type == 0 {
		for x := 0; x < s.Tx.N; x++ {
			AR_rx[x] = cmplx.Exp(1i * complex(s.k*s.Rx.dis*(float64(x)*math.Sin(s.Rx.Phi_RIS)*math.Cos(s.Rx.Theta_RIS)), 0))
		}
	} else if s.Tx.Type == 1 {
		dx := int(math.Sqrt(float64(s.Ris.N)))
		dy := dx

		for x := 0; x < dx; x++ {
			for y := 0; y < dy; y++ {
				AR_rx[x*dx+y] = cmplx.Exp(1i * complex(s.k*s.Rx.dis*(float64(x)*math.Sin(s.Rx.Phi_RIS)*math.Cos(s.Rx.Theta_RIS)+float64(y)*math.Sin(s.Rx.Theta_RIS)), 0))
			}
		}
	}

	for x := 0; x < len(AR_rx); x++ {
		for y := 0; y < len(AR_ris); y++ {
			AR_rx_ris.Data[y][x] = AR_rx[x] * AR_ris[y]
		}
	}

	eta := distuv.Uniform{Min: 0, Max: 2 * math.Pi, Src: rand.NewSource(1)} // Uniforma variable

	sf := distuv.Normal{Mu: 0, Sigma: math.Pow(s.sigma_LOS, 2), Src: rand.NewSource(1)} // variable loi normale for shadow fading
	ge := Ge(s.Ris.Theta_Rx)
	attenuation := L(s, sf, true, s.Ris.xyz, s.Rx.xyz)

	return cmat.Scale(AR_rx_ris, complex(math.Sqrt(ge*attenuation), 0)*cmplx.Exp(1i*complex(eta.Rand(), 0)))
}
