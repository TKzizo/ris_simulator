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

func (s *Simulation) D_channel(clusters []Cluster) cmat.Cmatrix {

	dnlos := DNLos(s, clusters)
	dlos := DLos(s)

	m := cmat.Add(dlos, dnlos)
	return m
}

func DNLos(s *Simulation, clusters []Cluster) cmat.Cmatrix {
	var nbr_scatterers int
	for _, cluster := range clusters {
		nbr_scatterers += len(cluster.Scatterers)
	}

	var AR_cs_rx, AR_cs_tx cmat.Cmatrix
	AR_cs_rx.Init(nbr_scatterers, s.Rx.N)
	AR_cs_tx.Init(nbr_scatterers, s.Tx.N)

	sf := distuv.Normal{Mu: 0, Sigma: math.Pow(s.Sigma_NLOS, 2), Src: rand.NewSource(uint64(time.Now().Unix()))}
	attenuation := make([]float64, nbr_scatterers, nbr_scatterers)
	beta := make([]complex128, nbr_scatterers, nbr_scatterers)
	// same complex gain ( small scale fading ) with and exess phase
	eta := make([]float64, nbr_scatterers, nbr_scatterers)

	// calculating the Array response
	index := 0
	for _, cluster := range clusters {
		for _, scatterer := range cluster.Scatterers {
			if s.Rx.Type == 0 {
				for x := 0; x < s.Rx.N; x++ {
					AR_cs_rx.Data[index][x] = cmplx.Exp(1i * complex(s.K*s.Rx.Dis*(float64(x)*math.Sin(scatterer.Phi_RX)*math.Cos(scatterer.Theta_RX)), 0))
				}
			}
			if s.Rx.Type == 1 {
				dx := int(math.Sqrt(float64(s.Rx.N)))
				dy := dx
				for x := 0; x < dx; x++ {
					for y := 0; y < dy; y++ {
						AR_cs_rx.Data[index][x*dx+y] = cmplx.Exp(1i * complex(s.K*s.Rx.Dis*(float64(x)*math.Sin(scatterer.Phi_RX)*math.Cos(scatterer.Theta_RX)+float64(y)*math.Sin(scatterer.Theta_RX)), 0))
					}
				}
			}
			if s.Tx.Type == 0 {
				for x := 0; x < s.Tx.N; x++ {
					AR_cs_tx.Data[index][x] = cmplx.Exp(1i * complex(s.K*s.Tx.Dis*(float64(x)*math.Sin(scatterer.Phi_RX)*math.Cos(scatterer.Theta_RX)), 0))
				}

			} else if s.Tx.Type == 1 {
				dx := int(math.Sqrt(float64(s.Tx.N)))
				dy := dx
				for x := 0; x < dx; x++ {
					for y := 0; y < dy; y++ {
						AR_cs_tx.Data[index][x*dx+y] = cmplx.Exp(1i * complex(s.K*s.Tx.Dis*(float64(x)*math.Sin(scatterer.Phi_TX)*math.Cos(scatterer.Theta_TX)+float64(y)*math.Sin(scatterer.Theta_TX)), 0))
					}
				}
			}

			attenuation[index] = L(s.Lambda, s.N_LOS, s.N_NLOS, s.B_LOS, s.B_NLOS, s.Frequency, s.F0, s.Sigma_LOS, s.Sigma_NLOS, sf, false, s.Rx.Xyz, s.Tx.Xyz, scatterer.Xyz)
			beta[index] = complex(rand.Float64()/math.Sqrt(2), rand.Float64()/math.Sqrt(2))
			eta[index] = s.K * (Distance(scatterer.Xyz, s.Ris.Xyz) - Distance(scatterer.Xyz, s.Rx.Xyz))
			index++
		}
	}

	c := cmat.Cmatrix{}
	c.Init(s.Rx.N, s.Tx.N)

	tmp := cmat.Cmatrix{}
	tmp.Init(s.Rx.N, s.Tx.N)

	for i := 0; i < nbr_scatterers; i++ {
		val := beta[i] * cmplx.Exp(1i*complex(eta[i], 0)) * complex(math.Sqrt(attenuation[i]), 0)
		//fmt.Println("beta: ", beta[i], " ge: ", ge[i], " att: ", attenuation[i])
		//fmt.Println("val: ", val)
		for x := 0; x < s.Tx.N; x++ {
			for y := 0; y < s.Rx.N; y++ {
				tmp.Data[y][x] = val * AR_cs_rx.Data[i][y] * AR_cs_tx.Data[i][x]
			}
		}

		c = cmat.Add(tmp, c)
	}
	//fmt.Println(cmat.Scale(c, complex(math.Sqrt(1.0/float64(nbr_scatterers)), 0)))
	//fmt.Scanln()

	return cmat.Scale(c, complex(math.Sqrt(1.0/float64(nbr_scatterers)), 0))

}

func DLos(s *Simulation) cmat.Cmatrix {

	random_src := rand.NewSource(uint64(time.Now().Unix()))
	//random_src := rand.NewSource(1)
	eta := distuv.Uniform{Min: 0, Max: 2 * math.Pi, Src: random_src} // shadow phase
	//fmt.Println(random_src)
	sf := distuv.Normal{Mu: 0, Sigma: math.Pow(s.Sigma_LOS, 2), Src: random_src}
	attenuation := math.Sqrt(L(s.Lambda, s.N_LOS, s.N_NLOS, s.B_LOS, s.B_NLOS, s.Frequency, s.F0, s.Sigma_LOS, s.Sigma_NLOS, sf, true, s.Rx.Xyz, s.Tx.Xyz))
	AR_tx_rx := cmat.Cmatrix{}
	AR_tx_rx.Init(s.Rx.N, s.Tx.N)

	AR_rx := make([]complex128, s.Rx.N, s.Rx.N)
	AR_tx := make([]complex128, s.Tx.N, s.Tx.N)

	//calculate angles between Tx-Rx (DLOS)
	s.Tx.Theta = float64(Sign(s.Rx.Xyz.Z, s.Tx.Xyz.Z)) * math.Atan2(math.Abs(s.Rx.Xyz.Z-s.Tx.Xyz.Z), Distance(s.Tx.Xyz, s.Rx.Xyz))
	s.Tx.Phi = float64(Sign(s.Tx.Xyz.Y, s.Rx.Xyz.Y)) * math.Atan2(math.Abs(s.Tx.Xyz.Y-s.Rx.Xyz.Y), math.Abs(s.Tx.Xyz.X-s.Rx.Xyz.X))
	s.Rx.Theta = math.Log(rand.Float64()/rand.Float64())*math.Sqrt(25/2) + rand.Float64()*180 - 90
	s.Rx.Phi = math.Log(rand.Float64()/rand.Float64())*math.Sqrt(25/2) + rand.Float64()*90 - 45

	if s.Tx.Type == 0 {
		for x := 0; x < s.Tx.N; x++ {
			AR_tx[x] = cmplx.Exp(1i * complex(s.K*s.Tx.Dis*(float64(x)*math.Sin(s.Tx.Phi)*math.Cos(s.Tx.Theta)), 0))
		}
	} else if s.Tx.Type == 1 {
		dx := int(math.Sqrt(float64(s.Tx.N)))
		dy := dx

		for x := 0; x < dx; x++ {
			for y := 0; y < dy; y++ {
				AR_tx[x*dx+y] = cmplx.Exp(1i * complex(s.K*s.Tx.Dis*(float64(x)*math.Sin(s.Tx.Phi)*math.Cos(s.Tx.Theta)+float64(y)*math.Sin(s.Tx.Theta)), 0))
			}
		}
	}

	if s.Rx.Type == 0 {
		for x := 0; x < s.Rx.N; x++ {
			AR_rx[x] = cmplx.Exp(1i * complex(s.K*s.Rx.Dis*(float64(x)*math.Sin(s.Rx.Phi)*math.Cos(s.Rx.Theta)), 0))
		}
	} else if s.Rx.Type == 1 {
		dx := int(math.Sqrt(float64(s.Rx.N)))
		dy := dx

		for x := 0; x < dx; x++ {
			for y := 0; y < dy; y++ {
				AR_rx[x*dx+y] = cmplx.Exp(1i * complex(s.K*s.Rx.Dis*(float64(x)*math.Sin(s.Rx.Phi)*math.Cos(s.Rx.Theta)+float64(y)*math.Sin(s.Rx.Theta)), 0))
			}
		}
	}

	for x := 0; x < len(AR_tx); x++ {
		for y := 0; y < len(AR_rx); y++ {
			AR_tx_rx.Data[y][x] = AR_tx[x] * AR_rx[y]
		}
	}
	//	fmt.Println(cmat.Scale(AR_tx_ris, complex(math.Sqrt(ge*attenuation), 0)*cmplx.Exp(1i*complex(eta.Rand(), 0))))
	//	fmt.Scanln()
	return cmat.Scale(AR_tx_rx, complex(attenuation, 0)*cmplx.Exp(1i*complex(eta.Rand(), 0)))
}
