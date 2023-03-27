package main

import (
	"math"
	"time"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

type Cluster struct {
	xyz        Coordinates
	mean_phi   float64 // mean azimuth
	mean_theta float64 // mean elevation
	Scatterers []Scatterer
}

type Scatterer struct {
	xyz       Coordinates
	Phi_RIS   float64
	Phi_TX    float64
	Theta_RIS float64
	Theta_TX  float64
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
				mean_phi:   rand.Float64()*180 - 90,
				mean_theta: rand.Float64()*90 - 45,
				Scatterers: make([]Scatterer, rand.Int()%29+1), // atleast one Sub-Ray per Cluster
			})
		GenerateAngles(&Clusters[i])
		GenerateCoordinates(&Clusters[i], s)
	}

	return []Cluster{}

}

func GenerateCoordinates(c *Cluster, s *Simulation) {
	position := 1 + rand.Float64()*(Distance(s.Ris.xyz, s.Tx.xyz)-1)

	// verification of coordinates over bound
	for {
		c.xyz.x = s.Tx.xyz.x + position*math.Cos(DegToRad(c.mean_theta))*math.Cos(DegToRad(c.mean_phi))
		c.xyz.y = s.Tx.xyz.y + position*math.Cos(DegToRad(c.mean_theta))*math.Sin(DegToRad(c.mean_phi))
		c.xyz.z = s.Tx.xyz.z + position*math.Sin(DegToRad(c.mean_theta))
		if c.xyz.z > s.Env.height || c.xyz.z < 0 || c.xyz.y > s.Env.width || c.xyz.y < 0 || c.xyz.x > s.Env.length || c.xyz.x < 0 {
			position = 0.8 * position
		} else {
			break
		}
	}

	// generating subRay Coodinates
	for i := 0; i < len(c.Scatterers); i++ {
		c.Scatterers[i].xyz.x = c.xyz.x + position*math.Cos(DegToRad(c.mean_theta))*math.Cos(DegToRad(c.mean_phi))

		c.Scatterers[i].xyz.y = c.xyz.y + position*math.Cos(DegToRad(c.mean_theta))*math.Sin(DegToRad(c.mean_phi))

		c.Scatterers[i].xyz.z = c.xyz.z + position*math.Sin(DegToRad(c.mean_theta))

		if c.Scatterers[i].xyz.z > s.Env.height || c.Scatterers[i].xyz.z < 0 || c.Scatterers[i].xyz.y > s.Env.width || c.Scatterers[i].xyz.y < 0 || c.Scatterers[i].xyz.x > s.Env.length || c.Scatterers[i].xyz.x < 0 {
			c.Scatterers = ignoreScatterer(c.Scatterers, i)
		}
	}

	// We need to have at least one Scatterer
	if len(c.Scatterers) == 1 {
		c.Scatterers = []Scatterer{
			{
				xyz:   Coordinates{x: c.xyz.x, y: c.xyz.y, z: c.xyz.z},
				phi:   c.mean_phi,
				theta: c.mean_theta,
			},
		}
	}
}

func GenerateAngles(c *Cluster) {
	for i := 0; i < len(c.Scatterers); i++ {
		c.Scatterers[i].phi = math.Log(rand.Float64()/rand.Float64())*math.Sqrt(25/2) + c.mean_phi
		c.Scatterers[i].theta = math.Log(rand.Float64()/rand.Float64())*math.Sqrt(25/2) + c.mean_theta

	}
}

func ignoreScatterer(scatterers []Scatterer, index int) []Scatterer {
	scatterers[index] = scatterers[len(scatterers)-1]
	return scatterers
}
