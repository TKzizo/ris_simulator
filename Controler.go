package main

import (
	"math"
	"math/cmplx"
)

func GetCoefficients(H, G []complex128) []complex128 {

	Theta_ris := []complex128{}
	for i := 0; i < len(H); i++ {
		phi_n := cmplx.Phase(H[i])
		psi_n := cmplx.Phase(G[i])
		Theta_ris = append(Theta_ris, cmplx.Rect(1, math.Remainder(-(phi_n+psi_n), 2*math.Pi)))
	}

	//fmt.Println("RIS_Coeff: ", Theta_ris)
	return Theta_ris
}
