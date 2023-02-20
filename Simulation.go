package main

import (
	"fmt"
	"math"
	"math/cmplx"

	"golang.org/x/exp/rand"
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
	Ris       RIS
	Tx        Tx_Rx
	Rx        Tx_Rx
	Frequency float64
	Lambda    float64 // wave length
	k         float64
	n         float64 // Path Loss exponent
	b         float64 // systemc parameter
	sigma     float64 // db
	channel   chan Updates

	//Scenario int // sideWall - oppositeWall
	//arrayType int // ULA - PA
}

func (s *Simulation) Setup() {
	s.Lambda = 3.0 / 10 * s.Frequency
	s.k = 2 * math.Pi / s.Lambda

	if s.n == 0.0 {
		s.n = 1.73
	}
	// we didn't check for s.b because as default value
	// we would have given it 0.0 which already its current value

	if s.sigma == 0.0 {
		s.sigma = 3.02
	}

	s.Ris.Setup(s.Lambda)
	/*
	   fmt.Println("RIS: ", s.Ris)
	   fmt.Println("Tx: ", s.Tx)
	   fmt.Println("Rx: ", s.Rx)
	   fmt.Println("Frequency: ", s.Frequency)
	   fmt.Println("Lambda: ", s.Lambda)
	   fmt.Println("k: ", s.k)
	   fmt.Println("n: ", s.n)
	   fmt.Println("sigma: ", s.sigma)
	*/
}
func (s *Simulation) H_channel() []complex128 {

	s.Ris.Phi_Tx = float64(sign(s.Ris.xyz.x, s.Tx.xyz.x)) * math.Atan2(math.Abs(s.Ris.xyz.x-s.Tx.xyz.x), math.Abs(s.Ris.xyz.y-s.Tx.xyz.y))
	s.Ris.Theta_Tx = float64(sign(s.Tx.xyz.z, s.Ris.xyz.z)) * math.Asin(math.Abs(s.Ris.xyz.z-s.Tx.xyz.z)/Distance(s.Ris.xyz, s.Tx.xyz))
	s.Tx.Phi_RIS = float64(sign(s.Tx.xyz.y, s.Ris.xyz.y)) * math.Atan2(math.Abs(s.Tx.xyz.y-s.Ris.xyz.y), math.Abs(s.Tx.xyz.x-s.Ris.xyz.x))
	s.Tx.Theta_RIS = float64(sign(s.Tx.xyz.z, s.Ris.xyz.z)) * math.Asin(math.Abs(s.Ris.xyz.z-s.Tx.xyz.z)/Distance(s.Ris.xyz, s.Tx.xyz))

	Ih_Ris_tx := distuv.Bernoulli{P: Determine_Pb(s.Ris.xyz, s.Tx.xyz), Src: rand.NewSource(rand.Uint64())} //bernoulli variable

	eta := distuv.Uniform{Min: 0, Max: 2 * math.Pi, Src: rand.NewSource(rand.Uint64())} // Uniforma variable

	sf := distuv.Normal{Mu: 0, Sigma: math.Pow(s.sigma, 2), Src: rand.NewSource(rand.Uint64())} // variable loi normale for shadow fading

	RIS_array_response := Array_Response_RIS_Tx(s, &s.Ris, &s.Tx)

	Ge_RIS := Ge(s.Ris.Theta_Tx)

	for i := 0; i < len(RIS_array_response); i++ {
		attenuation := L(s, sf, s.Ris.xyz, s.Tx.xyz) * Ge_RIS
		fmt.Println("H attenuation :", attenuation)
		bernoli := Ih_Ris_tx.Rand()
		fmt.Println("H Bernoulli :", bernoli)
		tmp1 := math.Sqrt(attenuation * bernoli)
		fmt.Println("H tmp1: ", tmp1)
		val1 := complex(tmp1, 0)
		fmt.Println("H val1 ", val1)
		val2 := (RIS_array_response[i] * cmplx.Exp(1i*complex(eta.Rand(), 0)))
		fmt.Println("H val2 ", val2)
		val12 := val1 * val2
		fmt.Println("H val12 ", val12)
		RIS_array_response[i] = val12
	}

	fmt.Println("H_Channel: ", RIS_array_response)
	return RIS_array_response
}
func (s *Simulation) G_channel() []complex128 {

	s.Ris.Phi_Rx = float64(sign(s.Ris.xyz.x, s.Rx.xyz.x)) * math.Atan2(math.Abs(s.Ris.xyz.x-s.Rx.xyz.x), math.Abs(s.Ris.xyz.y-s.Rx.xyz.y))
	s.Ris.Theta_Rx = float64(sign(s.Rx.xyz.z, s.Ris.xyz.z)) * math.Asin(math.Abs(s.Ris.xyz.z-s.Rx.xyz.z)/Distance(s.Ris.xyz, s.Rx.xyz))
	s.Rx.Phi_RIS = float64(sign(s.Rx.xyz.y, s.Ris.xyz.y)) * math.Atan2(math.Abs(s.Rx.xyz.y-s.Ris.xyz.y), math.Abs(s.Rx.xyz.x-s.Ris.xyz.x))
	s.Rx.Theta_RIS = float64(sign(s.Rx.xyz.z, s.Ris.xyz.z)) * math.Asin(math.Abs(s.Ris.xyz.z-s.Rx.xyz.z)/Distance(s.Ris.xyz, s.Rx.xyz))

	eta := distuv.Uniform{Min: 0, Max: 2 * math.Pi, Src: rand.NewSource(rand.Uint64())} // Uniforma variable

	sf := distuv.Normal{Mu: 0, Sigma: math.Pow(s.sigma, 2), Src: rand.NewSource(rand.Uint64())} // variable loi normale for shadow fading

	RIS_array_response := Array_Response_RIS_Rx(s, &s.Ris, &s.Rx)

	Ge_RIS := Ge(s.Ris.Theta_Rx)

	for i := 0; i < len(RIS_array_response); i++ {
		attenuation := L(s, sf, s.Ris.xyz, s.Tx.xyz) * Ge_RIS
		fmt.Println("H attenuation :", attenuation)
		tmp1 := math.Sqrt(attenuation)
		fmt.Println("H tmp1: ", tmp1)
		val1 := complex(tmp1, 0)
		fmt.Println("H val1 ", val1)
		val2 := (RIS_array_response[i] * cmplx.Exp(1i*complex(eta.Rand(), 0)))
		fmt.Println("H val2 ", val2)
		val12 := val1 * val2
		fmt.Println("H val12 ", val12)
		RIS_array_response[i] = val12
	}

	fmt.Println("H_Channel: ", RIS_array_response)
	return RIS_array_response
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
func (s *Simulation) Rate(H, G, Theta []complex128) float64 {

	var temp []complex128
	var res complex128
	for i, v := range G {
		temp = append(temp, v*Theta[i])
	}

	for i, v := range H {
		res += temp[i] * v
	}
	return math.Log2(1 + math.Pow(cmplx.Abs(res), 2)*Pt/P_n)
}
