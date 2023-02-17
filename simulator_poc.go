package main

import (
	"fmt"
	"math"
	"math/cmplx"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/stat/distuv"
)

type Coordinates struct {
	x float64
	y float64
	z float64
}

const (
	Nt                = 1    // number of transmiter antennas
	Nr                = 1    // number of receiver antennas
	N                 = 4    // number of RIS elements
	b                 = 0    // system parameter
	n                 = 1.73 //Path Loss Exponent
	Frequency         = 28
	lambda            = 3.0 / Frequency * 10 // wavelength
	Gain              = math.Pi
	q                 = 0.285      //
	dis               = lambda / 2 // distance between antennas and elements of transimters and RIS
	k         float64 = 2 * math.Pi / lambda
	sigma             = 3.02 // dB
)

var (
	xyz_RIS  Coordinates
	xyz_Tx   Coordinates
	xyz_Rx   Coordinates
	d_Tx_RIS float64
	d_RIS_Rx float64
)

func distance(a, b Coordinates) float64 {
	return math.Sqrt(math.Pow(a.x-b.x, 2) + math.Pow(a.y-b.y, 2) + math.Pow(a.z-b.z, 2))
}

func degToRad(a float64) float64 {
	return a * math.Pi / 180
}

func radToDeg(a float64) float64 {
	return a * 180 / math.Pi
}

func sign(a, b float64) int8 {
	if a >= b {
		return 1
	}
	return -1
}
func _Ge(theta float64) float64 {
	return Gain * math.Pow(math.Cos(theta), 2*q)
}

func _L(sf distuv.Normal, a, b Coordinates) float64 {
	return math.Pow(10, (-20*math.Log10(4*math.Pi/lambda)-10*n*math.Log10(distance(a, b))-sf.Rand())/10)
}

func array_response(phi_RIS, theta_RIS float64) []complex128 {
	var vec []complex128
	for x := 0; x < int(math.Sqrt(N)); x++ {
		for y := 0; y < int(math.Sqrt(N)); y++ {
			argument := k * dis * (float64(x)*math.Sin(theta_RIS) + float64(y)*math.Sin(phi_RIS)*math.Cos(theta_RIS))
			//fmt.Println("argument[", x, y, "]", argument)
			expo := cmplx.Exp(1i * complex(argument, 0))
			//fmt.Println("Expo array_response: ", expo)
			vec = append(vec, expo)
		}
	}
	//	fmt.Println(vec)
	return vec
}

func H_channel(xyz_Tx, xyz_RIS Coordinates) []complex128 {

	phi_RIS := float64(sign(xyz_RIS.x, xyz_Tx.x)) * math.Atan2(math.Abs(xyz_RIS.x-xyz_Tx.x), math.Abs(xyz_RIS.y-xyz_Tx.y))
	theta_RIS := float64(sign(xyz_Tx.z, xyz_RIS.z)) * math.Asin(math.Abs(xyz_RIS.z-xyz_Tx.z)/distance(xyz_RIS, xyz_Tx))
	//phi_Tx := float64(sign(xyz_Tx.y, xyz_RIS.y)) * math.Atan2(math.Abs(xyz_Tx.y-xyz_RIS.y), math.Abs(xyz_Tx.x-xyz_RIS.x))
	//theta_Tx := float64(sign(xyz_Tx.z, xyz_RIS.z)) * math.Asin(math.Abs(xyz_RIS.z-xyz_Tx.z)/distance(xyz_RIS, xyz_Tx))
	var pB float64
	if d := distance(xyz_RIS, xyz_Tx); d <= 1.2 {
		pB = 1
	} else if (1.2 < d) && (d <= 6.5) {
		pB = math.Exp((-d + 1.2) / 4.7)
	} else {
		pB = 0.32 * math.Exp((-d+6.5)/32.6)
	}

	fmt.Println("Pb bernoulli: ", pB)
	Ih_Ris_tx := distuv.Bernoulli{P: pB, Src: rand.NewSource(rand.Uint64())} //bernoulli variable

	eta := distuv.Uniform{Min: 0, Max: 2 * math.Pi, Src: rand.NewSource(rand.Uint64())} // Uniforma variable

	sf := distuv.Normal{Mu: 0, Sigma: sigma * sigma, Src: rand.NewSource(rand.Uint64())} // variable loi normale for shadow fading

	RIS_array_response := array_response(phi_RIS, theta_RIS)

	Ge_RIS := _Ge(theta_RIS)

	for i := 0; i < len(RIS_array_response); i++ {
		attenuation := _L(sf, xyz_Tx, xyz_RIS) * Ge_RIS
		fmt.Println("attenuation :", attenuation)
		bernoli := Ih_Ris_tx.Rand()
		fmt.Println("Bernoulli :", bernoli)
		tmp1 := math.Sqrt(attenuation * bernoli)
		fmt.Println("tmp1: ", tmp1)
		val1 := complex(tmp1, 0)
		fmt.Println("val1 ", val1)
		val2 := (RIS_array_response[i] * cmplx.Exp(1i*complex(eta.Rand(), 0)))
		fmt.Println("val2 ", val2)
		val12 := val1 * val2
		fmt.Println("val12 ", val12)
		RIS_array_response[i] = val12
	}

	fmt.Println("H_Channel: ", RIS_array_response)
	return RIS_array_response
}

func G_channel(xyz_Rx, xyz_RIS Coordinates) []complex128 {

	phi_RIS := float64(sign(xyz_RIS.x, xyz_Rx.x)) * math.Atan2(math.Abs(xyz_RIS.x-xyz_Rx.x), math.Abs(xyz_RIS.y-xyz_Rx.y))
	theta_RIS := float64(sign(xyz_Tx.x, xyz_RIS.x)) * math.Asin(math.Abs(xyz_RIS.z-xyz_Rx.z)/distance(xyz_RIS, xyz_Tx))

	eta := distuv.Uniform{Min: 0, Max: 2 * math.Pi, Src: rand.NewSource(rand.Uint64())} // Uniforma variable

	sf := distuv.Normal{Mu: 0, Sigma: sigma * sigma, Src: rand.NewSource(rand.Uint64())} // variable loi normale for shadow fading

	RIS_array_response := array_response(phi_RIS, theta_RIS)

	Ge_RIS := _Ge(theta_RIS)

	for i := 0; i < len(RIS_array_response); i++ {
		RIS_array_response[i] = complex(math.Sqrt(_L(sf, xyz_Rx, xyz_RIS)*Ge_RIS), 0) * (RIS_array_response[i] * cmplx.Exp(1i*complex(eta.Rand(), 0)))
	}

	fmt.Println("G_channel: ", RIS_array_response)
	return RIS_array_response
}

func RIS_Coefficients(H_channel, G_channel []complex128) []complex128 { //Optimal RIS coefficients for SISO setup

	Theta_ris := []complex128{}
	for i := 0; i < len(H_channel); i++ {
		phi_n := cmplx.Phase(H_channel[i])
		psi_n := cmplx.Phase(G_channel[i])
		Theta_ris = append(Theta_ris, cmplx.Rect(1, math.Remainder(-(phi_n+psi_n), 2*math.Pi)))
	}

	fmt.Println("RIS_Coeff: ", Theta_ris)
	return Theta_ris
}

func SNR(H_channel []complex128, G_channel []complex128, RIS_Coefficients []complex128, P_t, P_n float64) (float64, float64) {
	var temp []complex128
	var res complex128
	for i, v := range G_channel {
		temp = append(temp, v*RIS_Coefficients[i])
	}

	for i, v := range H_channel {
		res += temp[i] * v
	}
	return math.Log2(1 + math.Pow(cmplx.Abs(res), 2)*P_t/P_n), math.Pow(cmplx.Abs(res), 2) * P_t / P_n

}

func main() {

	xyz_Tx := Coordinates{x: 0, y: 25, z: 2}
	xyz_RIS := Coordinates{x: 40, y: 50, z: 2}
	xyz_Rx := Coordinates{x: 38, y: 48, z: 1}

	H := H_channel(xyz_Tx, xyz_RIS)
	G := G_channel(xyz_Rx, xyz_RIS)
	Theta_ris := RIS_Coefficients(H, G)

	r, snr := SNR(H, G, Theta_ris, 0.05, math.Pow(10, -9))

	fmt.Printf("Rate: %f SNR: %f", r, snr)
}
