// TODO Replace all cmplx.Exp with cmplx.Rect
package simulation

import (
	"math"
	"math/cmplx"
	"time"

	. "gitlab.eurecom.fr/ris-simulator/internal/utils"

	cmat "gitlab.eurecom.fr/ris-simulator/internal/reducedComplex"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

func (s *Simulation) H_channel(clusters []Cluster) cmat.Cmatrix {
	// Calculate Phi and theta for RIS-TX and TX-RIS (H-LOS)
	s.Tx.Phi_RIS = float64(Sign(s.Tx.Xyz.Y, s.Ris.Xyz.Y)) * math.Atan2(math.Abs(s.Tx.Xyz.Y-s.Ris.Xyz.Y), math.Abs(s.Tx.Xyz.X-s.Ris.Xyz.X))
	s.Ris.Theta_Tx = float64(Sign(s.Tx.Xyz.Z, s.Ris.Xyz.Z)) * math.Asin(math.Abs(s.Ris.Xyz.Z-s.Tx.Xyz.Z)/Distance(s.Ris.Xyz, s.Tx.Xyz))
	s.Tx.Theta_RIS = float64(Sign(s.Tx.Xyz.Z, s.Ris.Xyz.Z)) * math.Asin(math.Abs(s.Ris.Xyz.Z-s.Tx.Xyz.Z)/Distance(s.Ris.Xyz, s.Tx.Xyz))
	if s.Broadside == 0 { // side wall
		s.Ris.Phi_Tx = float64(Sign(s.Ris.Xyz.X, s.Tx.Xyz.X)) * math.Atan2(math.Abs(s.Ris.Xyz.X-s.Tx.Xyz.X), math.Abs(s.Ris.Xyz.Y-s.Tx.Xyz.Y))
	} else if s.Broadside == 1 { // Opposite Wall
		s.Ris.Phi_Tx = float64(Sign(s.Ris.Xyz.Y, s.Ris.Xyz.Y)) * math.Atan2(math.Abs(s.Ris.Xyz.Y-s.Tx.Xyz.Y), math.Abs(s.Ris.Xyz.X-s.Tx.Xyz.X))
	}
	//Calculate Phi and theta for RIS-Cluster and TX-Cluster (H-NLOS)
	if s.Broadside == 0 { //SideWall
		for i := 0; i < len(clusters); i++ {
			for y := 0; y < len(clusters[i].Scatterers); y++ {
				clusters[i].Scatterers[y].Phi_RIS = float64(Sign(s.Ris.Xyz.X, clusters[i].Scatterers[y].Xyz.X)) * math.Atan2(math.Abs(s.Ris.Xyz.X-clusters[i].Scatterers[y].Xyz.X), math.Abs(s.Ris.Xyz.Y-clusters[i].Scatterers[y].Xyz.Y))
				clusters[i].Scatterers[y].Phi_TX = float64(Sign(s.Tx.Xyz.Y, clusters[i].Scatterers[y].Xyz.Y)) * math.Atan2(math.Abs(clusters[i].Scatterers[y].Xyz.Y-s.Tx.Xyz.Y), math.Abs(clusters[i].Scatterers[y].Xyz.X-s.Tx.Xyz.X))
				clusters[i].Scatterers[y].Phi_TX = DegToRad(clusters[i].Scatterers[y].Phi_TX)
				//clusters[i].Scatterers[y].Theta_RIS = float64(Sign(clusters[i].Scatterers[y].Xyz.Z, s.Ris.Xyz.Z)) * math.Asin(math.Abs(s.Ris.Xyz.Z-clusters[i].Scatterers[y].Xyz.Z)/Distance(s.Ris.Xyz, clusters[i].Scatterers[y].Xyz))
				clusters[i].Scatterers[y].Theta_TX = DegToRad(clusters[i].Scatterers[y].Theta_TX)
				//clusters[i].Scatterers[y].Theta_TX = float64(Sign(clusters[i].Scatterers[y].Xyz.Z, s.Tx.Xyz.Z)) * math.Asin(math.Abs(clusters[i].Scatterers[y].Xyz.Z-s.Tx.Xyz.Z)/Distance(s.Tx.Xyz, clusters[i].Scatterers[y].Xyz))

			}
		}
	} else if s.Broadside == 1 { //OppositeWall
		for i := 0; i < len(clusters); i++ {
			for y := 0; y < len(clusters[i].Scatterers); y++ {
				//fmt.Println("Generate: Theta_TX: ", DegToRad(clusters[i].Scatterers[y].Theta_TX), " Phi_TX: ", DegToRad(clusters[i].Scatterers[y].Phi_TX))
				clusters[i].Scatterers[y].Phi_RIS = float64(Sign(clusters[i].Scatterers[y].Xyz.Y, s.Ris.Xyz.Y)) * math.Atan2(math.Abs(s.Ris.Xyz.Y-clusters[i].Scatterers[y].Xyz.Y), math.Abs(s.Ris.Xyz.X-clusters[i].Scatterers[y].Xyz.X))
				//clusters[i].Scatterers[y].Phi_TX = float64(Sign(s.Tx.Xyz.Y, clusters[i].Scatterers[y].Xyz.Y)) * math.Atan2(math.Abs(clusters[i].Scatterers[y].Xyz.Y-s.Tx.Xyz.Y), math.Abs(clusters[i].Scatterers[y].Xyz.X-s.Tx.Xyz.X))
				clusters[i].Scatterers[y].Phi_TX = DegToRad(clusters[i].Scatterers[y].Phi_TX)
				clusters[i].Scatterers[y].Theta_RIS = float64(Sign(clusters[i].Scatterers[y].Xyz.Z, s.Ris.Xyz.Z)) * math.Asin(math.Abs(s.Ris.Xyz.Z-clusters[i].Scatterers[y].Xyz.Z)/Distance(s.Ris.Xyz, clusters[i].Scatterers[y].Xyz))
				//clusters[i].Scatterers[y].Theta_TX = float64(Sign(clusters[i].Scatterers[y].Xyz.Z, s.Tx.Xyz.Z)) * math.Asin(math.Abs(clusters[i].Scatterers[y].Xyz.Z-s.Tx.Xyz.Z)/Distance(s.Tx.Xyz, clusters[i].Scatterers[y].Xyz))
				clusters[i].Scatterers[y].Theta_TX = DegToRad(clusters[i].Scatterers[y].Theta_TX)
				//fmt.Println("Calculated: Theta_TX: ", clusters[i].Scatterers[y].Theta_TX, " Phi_TX: ", clusters[i].Scatterers[y].Phi_TX)
			}
		}
	}
	hnlos := HNLos(s, clusters)
	hlos := HLos(s)
	//fmt.Println(hnlos)
	//fmt.Scanln()
	m := cmat.Add(hlos, hnlos)
	//	fmt.Println(m.Data[0][0], hlos.Data[0][0], hnlos.Data[0][0])
	//	fmt.Scanln()
	return m
}

func HLos(s *Simulation) cmat.Cmatrix {
	//Ih_Ris_tx := distuv.Bernoulli{P: Determine_Pb(s.Ris.Xyz, s.Tx.Xyz), Src: rand.NewSource(rand.Uint64())} //bernoulli variable
	random_src := rand.NewSource(uint64(time.Now().Unix()))
	//random_src := rand.NewSource(1)
	eta := distuv.Uniform{Min: 0, Max: 2 * math.Pi, Src: random_src} // shadow phase
	//fmt.Println(random_src)
	sf := distuv.Normal{Mu: 0, Sigma: math.Pow(s.Sigma_LOS, 2), Src: random_src}
	ge := Ge(s.Ris.Theta_Tx)
	attenuation := math.Sqrt(L(s.Lambda, s.N_LOS, s.N_NLOS, s.B_LOS, s.B_NLOS, s.Frequency, s.F0, s.Sigma_LOS, s.Sigma_NLOS, sf, true, s.Ris.Xyz, s.Tx.Xyz) * ge)
	AR_tx_ris := cmat.Cmatrix{}
	AR_tx_ris.Init(s.Ris.N, s.Tx.N)

	dx := int(math.Sqrt(float64(s.Ris.N)))
	dy := dx

	AR_ris := make([]complex128, s.Ris.N, s.Ris.N)
	AR_tx := make([]complex128, s.Tx.N, s.Tx.N)

	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			AR_ris[x*dx+y] = cmplx.Exp(
				1i * complex(s.K*s.Ris.Dis*(float64(x)*math.Sin(s.Ris.Theta_Tx)+float64(y)*math.Sin(s.Ris.Phi_Tx)*math.Cos(s.Ris.Theta_Tx)), 0))
		}
	}
	if s.Tx.Type == 0 {
		for x := 0; x < s.Tx.N; x++ {
			AR_tx[x] = cmplx.Exp(1i * complex(s.K*s.Tx.Dis*(float64(x)*math.Sin(s.Tx.Phi_RIS)*math.Cos(s.Tx.Theta_RIS)), 0))
		}
		//		fmt.Println("The TX type: ", s.Tx.Type)
		//		fmt.Println(AR_tx)
		//		fmt.Scanln()
	} else if s.Tx.Type == 1 {
		dx := int(math.Sqrt(float64(s.Tx.N)))
		dy := dx

		for x := 0; x < dx; x++ {
			for y := 0; y < dy; y++ {
				AR_tx[x*dx+y] = cmplx.Exp(1i * complex(s.K*s.Tx.Dis*(float64(x)*math.Sin(s.Tx.Phi_RIS)*math.Cos(s.Tx.Theta_RIS)+float64(y)*math.Sin(s.Tx.Theta_RIS)), 0))
			}
		}
		//		fmt.Println("The TX type: ", s.Tx.Type)
		//		fmt.Println(AR_tx)
		//		fmt.Scanln()
	}
	//fmt.Println(AR_tx)
	for x := 0; x < len(AR_tx); x++ {
		for y := 0; y < len(AR_ris); y++ {
			AR_tx_ris.Data[y][x] = AR_tx[x] * AR_ris[y]
		}
	}
	//	fmt.Println(cmat.Scale(AR_tx_ris, complex(math.Sqrt(ge*attenuation), 0)*cmplx.Exp(1i*complex(eta.Rand(), 0))))
	//	fmt.Scanln()
	return cmat.Scale(AR_tx_ris, complex(attenuation, 0)*cmplx.Exp(1i*complex(eta.Rand(), 0)))
}

func HNLos(s *Simulation, clusters []Cluster) cmat.Cmatrix {

	var nbr_scatterers int
	for _, cluster := range clusters {
		nbr_scatterers += len(cluster.Scatterers)
	}

	var AR_cs_ris, AR_cs_tx cmat.Cmatrix
	AR_cs_ris.Init(nbr_scatterers, s.Ris.N)
	AR_cs_tx.Init(nbr_scatterers, s.Tx.N)

	sf := distuv.Normal{Mu: 0, Sigma: math.Pow(s.Sigma_NLOS, 2), Src: rand.NewSource(uint64(time.Now().Unix()))}
	ge := make([]float64, nbr_scatterers, nbr_scatterers)
	attenuation := make([]float64, nbr_scatterers, nbr_scatterers)
	beta := make([]complex128, nbr_scatterers, nbr_scatterers)

	dx := int(math.Sqrt(float64(s.Ris.N)))
	dy := dx
	index := 0
	for _, cluster := range clusters {
		for _, scatterer := range cluster.Scatterers {
			for x := 0; x < dx; x++ {
				for y := 0; y < dy; y++ {
					AR_cs_ris.Data[index][x*dx+y] = cmplx.Exp(1i * complex(s.K*s.Ris.Dis*(float64(x)*math.Sin(scatterer.Theta_RIS)+float64(y)*math.Sin(scatterer.Phi_RIS)*math.Cos(scatterer.Theta_RIS)), 0))
				}
			}
			if s.Tx.Type == 0 {

				for x := 0; x < s.Tx.N; x++ {

					AR_cs_tx.Data[index][x] = cmplx.Exp(1i * complex(s.K*s.Tx.Dis*(math.Sin(scatterer.Phi_TX)*math.Cos(scatterer.Theta_TX)), 0))
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
			ge[index] = Ge(scatterer.Theta_RIS)
			attenuation[index] = L(s.Lambda, s.N_LOS, s.N_NLOS, s.B_LOS, s.B_NLOS, s.Frequency, s.F0, s.Sigma_LOS, s.Sigma_NLOS, sf, false, s.Ris.Xyz, s.Tx.Xyz, scatterer.Xyz)
			beta[index] = complex(rand.Float64()/math.Sqrt(2), rand.Float64()/math.Sqrt(2))
			index++
		}
	}

	c := cmat.Cmatrix{}
	c.Init(s.Ris.N, s.Tx.N)

	tmp := cmat.Cmatrix{}
	tmp.Init(s.Ris.N, s.Tx.N)

	for i := 0; i < nbr_scatterers; i++ {
		val := beta[i] * complex(math.Sqrt(ge[i]*attenuation[i]), 0)
		//fmt.Println("beta: ", beta[i], " ge: ", ge[i], " att: ", attenuation[i])
		//fmt.Println("val: ", val)
		for x := 0; x < s.Tx.N; x++ {
			for y := 0; y < s.Ris.N; y++ {
				tmp.Data[y][x] = val * AR_cs_ris.Data[i][y] * AR_cs_tx.Data[i][x]
			}
		}

		c = cmat.Add(tmp, c)
	}
	//fmt.Println(cmat.Scale(c, complex(math.Sqrt(1.0/float64(nbr_scatterers)), 0)))
	//fmt.Scanln()

	return cmat.Scale(c, complex(math.Sqrt(1.0/float64(nbr_scatterers)), 0))
}
