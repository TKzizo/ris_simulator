package simulation

import (
	"math"
	"math/cmplx"
	"time"

	cmat "gitlab.eurecom.fr/ris-simulator/internal/reducedComplex"
	. "gitlab.eurecom.fr/ris-simulator/internal/utils"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

func (s *Simulation) G_channel() cmat.Cmatrix {

	s.Ris.Theta_Rx = float64(Sign(s.Rx.Xyz.Z, s.Ris.Xyz.Z)) * math.Asin(math.Abs(s.Rx.Xyz.Z-s.Ris.Xyz.Z)/Distance(s.Ris.Xyz, s.Rx.Xyz))
	//s.Rx.Theta_RIS = float64(Sign(s.Ris.Xyz.Z, s.Rx.Xyz.Z)) * math.Asin(math.Abs(s.Rx.Xyz.Z-s.Ris.Xyz.Z)/Distance(s.Ris.Xyz, s.Rx.Xyz))

	s.Rx.Phi_RIS = DegToRad(math.Log(rand.Float64()/rand.Float64())*math.Sqrt(25/2) + rand.Float64()*180 - 90)
	s.Rx.Theta_RIS = DegToRad(math.Log(rand.Float64()/rand.Float64())*math.Sqrt(25/2) + rand.Float64()*180 - 90)
	if s.Broadside == 0 { //Side Wall
		s.Ris.Phi_Rx = float64(Sign(s.Ris.Xyz.X, s.Rx.Xyz.X)) * math.Atan2(math.Abs(s.Rx.Xyz.X-s.Ris.Xyz.X), math.Abs(s.Rx.Xyz.Y-s.Ris.Xyz.Y))

		//s.Rx.Phi_RIS = float64(Sign(s.Ris.Xyz.Y, s.Rx.Xyz.Y)) * math.Atan2(math.Abs(s.Rx.Xyz.Y-s.Ris.Xyz.Y), math.Abs(s.Rx.Xyz.X-s.Ris.Xyz.X))
	} else if s.Broadside == 1 { // Opposite wall
		s.Ris.Phi_Rx = float64(Sign(s.Rx.Xyz.Y, s.Ris.Xyz.Y)) * math.Atan2(math.Abs(s.Rx.Xyz.Y), math.Abs(s.Rx.Xyz.Y))
		//s.Rx.Phi_RIS = float64(Sign(s.Rx.Xyz.X, s.Ris.Xyz.Y)) * math.Atan2(math.Abs(s.Rx.Xyz.Y-s.Ris.Xyz.Z), math.Abs(s.Rx.Xyz.X-s.Ris.Xyz.X))

	}

	AR_rx_ris := cmat.Cmatrix{}
	AR_rx_ris.Init(s.Ris.N, s.Rx.N)

	dx := int(math.Sqrt(float64(s.Ris.N)))
	dy := dx

	AR_ris := make([]complex128, s.Ris.N, s.Ris.N)
	AR_rx := make([]complex128, s.Rx.N, s.Rx.N)

	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			AR_ris[x*dx+y] = cmplx.Exp(
				1i * complex(s.K*s.Ris.Dis*(float64(x)*math.Sin(s.Ris.Theta_Rx)+float64(y)*math.Sin(s.Ris.Phi_Rx)*math.Cos(s.Ris.Theta_Rx)), 0))
		}
	}

	if s.Rx.Type == 0 {
		for x := 0; x < s.Rx.N; x++ {
			AR_rx[x] = cmplx.Exp(1i * complex(s.K*s.Rx.Dis*(float64(x)*math.Sin(s.Rx.Phi_RIS)*math.Cos(s.Rx.Theta_RIS)), 0))
		}
	} else if s.Rx.Type == 1 {
		dx := int(math.Sqrt(float64(s.Rx.N)))
		dy := dx

		for x := 0; x < dx; x++ {
			for y := 0; y < dy; y++ {
				AR_rx[x*dx+y] = cmplx.Exp(1i * complex(s.K*s.Rx.Dis*(float64(x)*math.Sin(s.Rx.Phi_RIS)*math.Cos(s.Rx.Theta_RIS)+float64(y)*math.Sin(s.Rx.Theta_RIS)), 0))
			}
		}
	}
	for x := 0; x < len(AR_rx); x++ {
		for y := 0; y < len(AR_ris); y++ {
			AR_rx_ris.Data[y][x] = AR_rx[x] * AR_ris[y]
		}
	}
	random_src := rand.NewSource(uint64(time.Now().Unix()))
	random_src = rand.NewSource(1)
	eta := distuv.Uniform{Min: 0, Max: 2 * math.Pi, Src: random_src} // Uniforma variable

	sf := distuv.Normal{Mu: 0, Sigma: math.Pow(s.Sigma_LOS, 2), Src: random_src} // variable loi normale for shadow fading
	ge := Ge(s.Ris.Theta_Rx)
	attenuation := L(s.Lambda, s.N_LOS, s.N_NLOS, s.B_LOS, s.B_NLOS, s.Frequency, s.F0, s.Sigma_LOS, s.Sigma_NLOS, sf, true, s.Ris.Xyz, s.Rx.Xyz)

	return cmat.Transpose(cmat.Scale(AR_rx_ris, complex(math.Sqrt(ge*attenuation), 0)*cmplx.Exp(1i*complex(eta.Rand(), 0))))
}
