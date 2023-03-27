package main

import (
	cmat "RIS_SIMULATOR/reducedComplex"
	"math"
	"math/cmplx"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

func (s *Simulation) H_channel(clusters []Cluster) cmat.Cmatrix {
	// Calculate Phi and theta for RIS-TX and TX-RIS
	s.Ris.Phi_Tx = float64(sign(s.Ris.xyz.x, s.Tx.xyz.x)) * math.Atan2(math.Abs(s.Ris.xyz.x-s.Tx.xyz.x), math.Abs(s.Ris.xyz.y-s.Tx.xyz.y))
	s.Ris.Theta_Tx = float64(sign(s.Tx.xyz.z, s.Ris.xyz.z)) * math.Asin(math.Abs(s.Ris.xyz.z-s.Tx.xyz.z)/Distance(s.Ris.xyz, s.Tx.xyz))
	s.Tx.Phi_RIS = float64(sign(s.Tx.xyz.y, s.Ris.xyz.y)) * math.Atan2(math.Abs(s.Tx.xyz.y-s.Ris.xyz.y), math.Abs(s.Tx.xyz.x-s.Ris.xyz.x))
	s.Tx.Theta_RIS = float64(sign(s.Tx.xyz.z, s.Ris.xyz.z)) * math.Asin(math.Abs(s.Ris.xyz.z-s.Tx.xyz.z)/Distance(s.Ris.xyz, s.Tx.xyz))

	//Calculate Phi and theta for RIS-Cluster and TX-Cluster
	//SideWall
	if s.Broadside == 0 {
		for i := 0; i < len(clusters); i++ {
			for y := 0; y < len(clusters[i].Scatterers); y++ {
				clusters[i].Scatterers[y].Phi_RIS = float64(sign(s.Ris.xyz.x, clusters[i].Scatterers[y].xyz.x)) * math.Atan2(math.Abs(s.Ris.xyz.x-clusters[i].Scatterers[y].xyz.x), math.Abs(s.Ris.xyz.y-clusters[i].Scatterers[y].xyz.y))
				clusters[i].Scatterers[y].Phi_TX = float64(sign(s.Tx.xyz.y, clusters[i].Scatterers[y].xyz.y)) * math.Atan2(math.Abs(clusters[i].Scatterers[y].xyz.y-s.Tx.xyz.y), math.Abs(clusters[i].Scatterers[y].xyz.x-s.Tx.xyz.x))
				clusters[i].Scatterers[y].Theta_RIS = float64(sign(clusters[i].Scatterers[y].xyz.z, s.Ris.xyz.z)) * math.Asin(math.Abs(s.Ris.xyz.z-clusters[i].Scatterers[y].xyz.z)/Distance(s.Ris.xyz, clusters[i].Scatterers[y].xyz))
				clusters[i].Scatterers[y].Theta_TX = float64(sign(clusters[i].Scatterers[y].xyz.z, s.Tx.xyz.z)) * math.Asin(math.Abs(clusters[i].Scatterers[y].xyz.z-s.Tx.xyz.z)/Distance(s.Tx.xyz, clusters[i].Scatterers[y].xyz))
			}
		}
		//OppositeWall
	} else if s.Broadside == 1 {
		for i := 0; i < len(clusters); i++ {
			for y := 0; y < len(clusters[i].Scatterers); y++ {
				clusters[i].Scatterers[y].Phi_RIS = float64(sign(clusters[i].Scatterers[y].xyz.y, s.Ris.xyz.y)) * math.Atan2(math.Abs(s.Ris.xyz.y-clusters[i].Scatterers[y].xyz.y), math.Abs(s.Ris.xyz.x-clusters[i].Scatterers[y].xyz.x))
				clusters[i].Scatterers[y].Phi_TX = float64(sign(s.Tx.xyz.y, clusters[i].Scatterers[y].xyz.y)) * math.Atan2(math.Abs(clusters[i].Scatterers[y].xyz.y-s.Tx.xyz.y), math.Abs(clusters[i].Scatterers[y].xyz.x-s.Tx.xyz.x))
				clusters[i].Scatterers[y].Theta_RIS = float64(sign(clusters[i].Scatterers[y].xyz.z, s.Ris.xyz.z)) * math.Asin(math.Abs(s.Ris.xyz.z-clusters[i].Scatterers[y].xyz.z)/Distance(s.Ris.xyz, clusters[i].Scatterers[y].xyz))
				clusters[i].Scatterers[y].Theta_TX = float64(sign(clusters[i].Scatterers[y].xyz.z, s.Tx.xyz.z)) * math.Asin(math.Abs(clusters[i].Scatterers[y].xyz.z-s.Tx.xyz.z)/Distance(s.Tx.xyz, clusters[i].Scatterers[y].xyz))
			}
		}
	}

	//Ih_Ris_tx := distuv.Bernoulli{P: Determine_Pb(s.Ris.xyz, s.Tx.xyz), Src: rand.NewSource(rand.Uint64())} //bernoulli variable

	eta := distuv.Uniform{Min: 0, Max: 2 * math.Pi, Src: rand.NewSource(1)} // Uniforma variable

	sf := distuv.Normal{Mu: 0, Sigma: math.Pow(s.sigma, 2), Src: rand.NewSource(1)} // variable loi normale for shadow fading

	Ge_RIS := Ge(s.Ris.Theta_Tx)
	attenuation := math.Sqrt(L(s, sf, s.Ris.xyz, s.Tx.xyz) * Ge_RIS)
	scalar := cmplx.Rect(attenuation, eta.Rand())

	return cmat.Scale(Array_Response_H(s, clusters), scalar)
}

func Array_Response_H(s *Simulation, clusters []Cluster) cmat.Cmatrix {
	var nbr_scatterers int
	for _, cluster := range clusters {
		nbr_scatterers += len(cluster.Scatterers)
	}

	var AR_cs_ris, AR_cs_tx cmat.Cmatrix
	AR_cs_ris.Init(nbr_scatterers, s.Ris.N)
	AR_cs_tx.Init(nbr_scatterers, s.Tx.N)
	dx := int(math.Sqrt(float64(s.Ris.N)))
	dy := dx
	index := 0
	for _, cluster := range clusters {
		for _, scatterer := range cluster.Scatterers {
			index++
			for x := 0; x < dx; x++ {
				for y := 0; y < dy; y++ {
					AR_cs_ris.Data[index][x*dx+y] = cmplx.Exp(1i * complex(s.k*s.Ris.dis*(float64(x)*math.Sin(scatterer.Theta_RIS)+float64(y)*math.Sin(scatterer.Phi_RIS)*math.Cos(scatterer.Theta_RIS)), 0))
				}
			}
			if s.Tx.Type == 0 {
				for x := 0; x < s.Tx.N; x++ {
					AR_cs_tx.Data[index][x] = cmplx.Exp(1i * complex(s.k*s.Ris.dis*(math.Sin(scatterer.Phi_TX)*math.Cos(scatterer.Theta_TX)), 0))
				}

			} else if s.Tx.Type == 1 {
				dx := int(math.Sqrt(float64(s.Tx.N)))
				dy := dx
				for x := 0; x < dx; x++ {
					for y := 0; y < dy; y++ {
						AR_cs_tx.Data[index][x*dx+y] = cmplx.Exp(1i * complex(s.k*s.Ris.dis*(float64(x)*math.Sin(scatterer.Phi_TX)*math.Cos(scatterer.Theta_TX)+float64(y)*math.Sin(scatterer.Theta_TX)), 0))
					}
				}
			}
		}
	}

	return cmat.Mul(cmat.Transpose(AR_cs_ris), AR_cs_tx)
}
