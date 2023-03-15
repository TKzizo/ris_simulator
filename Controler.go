package main

import (
	"math"
	"math/cmplx"

	"gonum.org/v1/gonum/mat"
)

func GetCoefficients(H, G mat.CDense) mat.CDense { //SISO ie: H and G are column vector of the same size
	r, _ := H.Dims()
	Theta_ris := *mat.NewCDense(r, r, nil)
	for i := 0; i < r; i++ {
		phi_n := cmplx.Phase(H.At(i, 0))
		psi_n := cmplx.Phase(G.At(i, 0))
		Theta_ris.Set(i, i, cmplx.Rect(1, math.Remainder(-(phi_n+psi_n), 2*math.Pi)))
	}
	return Theta_ris
}

/*func GetCoefficients(H, G mat.CDense) mat.CDiagonal { //SISO ie: H and G are column vector of the same size
	r, _ := H.Dims()
	Theta_ris := mat.NewDiagCDense(r,nil)
	for i := 0; i < r; i++ {
		phi_n := cmplx.Phase(H.At(i, 0))
		psi_n := cmplx.Phase(G.At(i, 0))
		Theta_ris.SetDiag(i, cmplx.Rect(1, math.Remainder(-(phi_n+psi_n), 2*math.Pi)))
	}
	return Theta_ris
}*/
