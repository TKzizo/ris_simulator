package main

import (
	cmat "RIS_SIMULATOR/reducedComplex"
	"math"
	"math/cmplx"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distuv"
)

func (s *Simulation) G_channel() cmat.Cmatrix {

	var G mat.CDense

	s.Ris.Phi_Rx = float64(sign(s.Ris.xyz.x, s.Rx.xyz.x)) * math.Atan2(math.Abs(s.Ris.xyz.x-s.Rx.xyz.x), math.Abs(s.Ris.xyz.y-s.Rx.xyz.y))
	s.Ris.Theta_Rx = float64(sign(s.Rx.xyz.z, s.Ris.xyz.z)) * math.Asin(math.Abs(s.Ris.xyz.z-s.Rx.xyz.z)/Distance(s.Ris.xyz, s.Rx.xyz))
	s.Rx.Phi_RIS = float64(sign(s.Rx.xyz.y, s.Ris.xyz.y)) * math.Atan2(math.Abs(s.Rx.xyz.y-s.Ris.xyz.y), math.Abs(s.Rx.xyz.x-s.Ris.xyz.x))
	s.Rx.Theta_RIS = float64(sign(s.Rx.xyz.z, s.Ris.xyz.z)) * math.Asin(math.Abs(s.Ris.xyz.z-s.Rx.xyz.z)/Distance(s.Ris.xyz, s.Rx.xyz))

	eta := distuv.Uniform{Min: 0, Max: 2 * math.Pi, Src: rand.NewSource(1)} // Uniforma variable

	sf := distuv.Normal{Mu: 0, Sigma: math.Pow(s.sigma, 2), Src: rand.NewSource(1)} // variable loi normale for shadow fading

	RX_array_response := Array_Response(s.k, int(math.Sqrt(float64(s.Rx.N))), int(math.Sqrt(float64(s.Rx.N))), s.Ris.dis, s.Rx.Phi_RIS, s.Rx.Theta_RIS)
	RIS_array_response := Array_Response(s.k, int(math.Sqrt(float64(s.Ris.N))), int(math.Sqrt(float64(s.Ris.N))), s.Ris.dis, s.Ris.Phi_Rx, s.Ris.Theta_Rx)
	G.Mul(RIS_array_response, RX_array_response.T())
	Ge_RIS := Ge(s.Ris.Theta_Rx)
	attenuation := math.Sqrt(L(s, sf, s.Ris.xyz, s.Tx.xyz) * Ge_RIS)
	scalar := cmplx.Rect(attenuation, eta.Rand())
	G.Scale(scalar, &G)

	return G
}
