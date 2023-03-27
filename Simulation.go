package main

import (
	"fmt"
	"math"
	"math/cmplx"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distuv"
)

const (
	q    float64 = 0.283 // related to the gain
	Gain float64 = math.Pi
	Pt   float64 = 0.05        // Power of transmitter
	P_n  float64 = 0.000000001 // variance of noise at the receiver
)

type Updates struct {
	Rx Coordinates
	Tx Coordinates
}

type Simulation struct {
	Env       Environment
	Ris       RIS
	Tx        Tx_Rx
	Rx        Tx_Rx
	Frequency float64
	Lambda_p  float64
	Lambda    float64 // wave length
	k         float64
	PLE       float64 // Path Loss exponent
	b         float64 // systemc parameter
	sigma     float64 // db
	channel   chan Updates
	Broadside int8 // 0: SideWall 1: OppositeWall
	//arrayType int // ULA - PA
}

func (s *Simulation) Setup() {
	s.Lambda = 3.0 / 10 * s.Frequency // it's Simplified so it only supports GHz
	s.k = 2 * math.Pi / s.Lambda

	if s.Frequency == 28.0 {
		s.Lambda_p = 1.8
	} else if s.Frequency == 73.0 {
		s.Lambda_p = 1.9
	}

	if s.PLE == 0.0 { //Pathloss exponent
		s.PLE = 1.73
	}
	// we didn't check for s.b because as default value
	// we would have given it 0.0 which already its current value

	if s.sigma == 0.0 {
		s.sigma = 3.02
	}

	s.Ris.Setup(s.Lambda)
}

func (s *Simulation) G_channel() mat.CDense {

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

func (s *Simulation) Update() {
	for {
		select {
		case u := <-s.channel:
			s.Rx.xyz = u.Rx
			s.Tx.xyz = u.Tx
		}
	}
}

func (s *Simulation) Rate_SNR(H, G, Theta mat.CDense, SNR int) (float64, float64) {

	var temp1 mat.CDense
	var temp2 mat.CDense

	temp1.Mul(G.T(), &Theta)
	temp2.Mul(&temp1, &H)
	fmt.Println(temp2.RawCMatrix().Data)

	snr := math.Pow(cmplx.Abs(temp2.RawCMatrix().Data[0]), 2) * Pt / P_n

	return math.Log2(1 + snr), float64(SNR) * snr
}
