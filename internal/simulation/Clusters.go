package simulation

import (
	"math"
	"time"

	. "gitlab.eurecom.fr/ris-simulator/internal/structures"
	. "gitlab.eurecom.fr/ris-simulator/internal/utils"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

type Cluster struct {
	Xyz        Coordinates
	Mean_phi   float64 // mean azimuth
	Mean_theta float64 // mean elevation
	Scatterers []Scatterer
}

type Scatterer struct {
	Xyz       Coordinates
	Phi_RIS   float64
	Phi_TX    float64
	Phi_RX    float64
	Theta_RIS float64
	Theta_TX  float64
	Theta_RX  float64
}

func GenerateClusters(s *Simulation) []Cluster {
	nbrClusters := int(distuv.Poisson{Lambda: s.Lambda_p, Src: rand.NewSource(uint64(time.Now().Unix()))}.Rand())
	if nbrClusters < 1 {
		nbrClusters = 1
	}
	var Clusters []Cluster
	for i := 0; i < nbrClusters; i++ {
		Clusters = append(Clusters,
			Cluster{
				Mean_phi:   rand.Float64()*180 - 90,
				Mean_theta: rand.Float64()*90 - 45,
				Scatterers: make([]Scatterer, rand.Int()%29+1), // atleast one Sub-Ray per Cluster
			})
		GenerateAngles(&Clusters[i])
		GenerateCoordinates(&Clusters[i], s)
	}

	return Clusters

}

func GenerateCoordinates(c *Cluster, s *Simulation) {
	position := 1 + rand.Float64()*(Distance(s.Ris.Xyz, s.Tx.Xyz)-1)

	// verification of coordinates over bound
	for {
		c.Xyz.X = s.Tx.Xyz.X + position*math.Cos(DegToRad(c.Mean_theta))*math.Cos(DegToRad(c.Mean_phi))
		c.Xyz.Y = s.Tx.Xyz.Y - position*math.Cos(DegToRad(c.Mean_theta))*math.Sin(DegToRad(c.Mean_phi))
		c.Xyz.Z = s.Tx.Xyz.Z + position*math.Sin(DegToRad(c.Mean_theta))
		if c.Xyz.Z > s.Env.Height || c.Xyz.Z < 0 || c.Xyz.Y > s.Env.Width || c.Xyz.Y < 0 || c.Xyz.X > s.Env.Length || c.Xyz.X < 0 {
			position = 0.8 * position
		} else {
			break
		}
	}

	// generating subRay Coodinates
	i := 0
	for i < len(c.Scatterers) {
		c.Scatterers[i].Xyz.X = s.Tx.Xyz.X + position*math.Cos(DegToRad(c.Scatterers[i].Theta_TX))*math.Cos(DegToRad(c.Scatterers[i].Phi_TX))

		c.Scatterers[i].Xyz.Y = s.Tx.Xyz.Y - position*math.Cos(DegToRad(c.Scatterers[i].Theta_TX))*math.Sin(DegToRad(c.Scatterers[i].Phi_TX))

		c.Scatterers[i].Xyz.Z = s.Tx.Xyz.Z + position*math.Sin(DegToRad(c.Scatterers[i].Theta_TX))

		if c.Scatterers[i].Xyz.Z > s.Env.Height || c.Scatterers[i].Xyz.Z < 0 || c.Scatterers[i].Xyz.Y > s.Env.Width || c.Scatterers[i].Xyz.Y < 0 || c.Scatterers[i].Xyz.X > s.Env.Length || c.Scatterers[i].Xyz.X < 0 {
			c.Scatterers = ignoreScatterer(c.Scatterers, i)
			continue
		}
		i++
	}

	// We need to have at least one Scatterer
	if len(c.Scatterers) == 1 {
		c.Scatterers = []Scatterer{
			{
				Xyz:      Coordinates{X: c.Xyz.X, Y: c.Xyz.Y, Z: c.Xyz.Z},
				Phi_TX:   c.Mean_phi,
				Theta_TX: c.Mean_theta,
			},
		}
	}
}

func GenerateAngles(c *Cluster) {

	for i := 0; i < len(c.Scatterers); i++ {
		c.Scatterers[i].Phi_TX = math.Log(rand.Float64()/rand.Float64())*math.Sqrt(25/2) + c.Mean_phi
		c.Scatterers[i].Phi_RX = math.Log(rand.Float64()/rand.Float64())*math.Sqrt(25/2) + c.Mean_phi
		c.Scatterers[i].Theta_TX = math.Log(rand.Float64()/rand.Float64())*math.Sqrt(25/2) + c.Mean_theta
		c.Scatterers[i].Theta_RX = math.Log(rand.Float64()/rand.Float64())*math.Sqrt(25/2) + c.Mean_theta
	}
}

func ignoreScatterer(scatterers []Scatterer, index int) []Scatterer {
	l := scatterers[:len(scatterers)-1]
	if index < len(l) {
		l[index] = scatterers[len(scatterers)-1]
	}
	return l
}
