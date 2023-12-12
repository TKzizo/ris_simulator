package utils

import (
	"fmt"
	"math"
	"math/cmplx"
	"strconv"

	cmat "gitlab.eurecom.fr/ris-simulator/internal/reducedComplex"
	. "gitlab.eurecom.fr/ris-simulator/internal/structures"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distuv"
)

func Distance(a, b Coordinates) float64 {
	return math.Sqrt(math.Pow(a.X-b.X, 2) + math.Pow(a.Y-b.Y, 2) + math.Pow(a.Z-b.Z, 2))
}
func RadToDeg(a float64) float64 {
	return a * 180 / math.Pi
}
func DegToRad(a float64) float64 {
	return a * math.Pi / 180
}
func Sign(a, b float64) int8 {
	if a > b {
		return 1
	} else if a == b {
		return 0
	}
	return -1
}
func Ge(theta float64) float64 {
	return Gain * math.Pow(math.Cos(theta), 2*Q)
}
func L(Lambda float64, N_LOS float64, N_NLOS float64, B_LOS float64, B_NLOS float64, Frequency float64, F0 float64, Sigma_LOS float64, Sigma_NLOS float64, sf distuv.Normal, los bool, a ...Coordinates) float64 {
	if los == true {
		return math.Pow(10, (-20*math.Log10(4*math.Pi/Lambda)-10*N_LOS*(1+B_LOS*((Frequency-F0)/F0))*math.Log10(Distance(a[0], a[1]))-(rand.Float64()*Sigma_LOS))/10)
	} else {
		return math.Pow(10, (-20*math.Log10(4*math.Pi/Lambda)-10*N_NLOS*(1+B_NLOS*((Frequency-F0)/F0))*math.Log10(Distance(a[0], a[1])+Distance(a[1], a[2]))-(rand.Float64()*Sigma_NLOS))/10)
	}
}

func Determine_Pb(a, b Coordinates) float64 {
	d := Distance(a, b)
	//LOS probability in Indoor office
	if d <= 1.2 {
		return 1
	} else if (1.2 < d) && (d <= 6.5) {
		return math.Exp((-d + 1.2) / 4.7)
	} else {
		return 0.32 * math.Exp((-d+6.5)/32.6)
	}
}

func cmat2cdense(matr cmat.Cmatrix) *mat.CDense {
	size := matr.Row * matr.Col
	data := make([]complex128, size)
	for x := 0; x < matr.Row; x++ {
		for y := 0; y < matr.Col; y++ {
			data[x*matr.Col+y] = matr.Data[x][y]
		}

	}
	return mat.NewCDense(matr.Row, matr.Col, data)
}

func ComplextoString(c complex128) string {
	bitSize := 128
	if bitSize != 64 && bitSize != 128 {
		panic("invalid bitSize")
	}
	bitSize >>= 1 // complex64 uses float32 internally

	// Check if imaginary part has a sign. If not, add one.
	im := strconv.FormatFloat(imag(c), 'e', 4, bitSize)
	if im[0] != '+' && im[0] != '-' {
		im = "+" + im
	}

	return strconv.FormatFloat(real(c), 'e', 4, bitSize) + im + "i"

}

/*func rate(H, G mat.CDense, Theta mat.CDiagonal) float64 {

	var temp1 mat.CDense
	var temp2 mat.CDense
	rate := 0.0

	temp1.Mul(G.T(), Theta)
	temp2.Mul(&temp1, &H)
	rate = math.Log2(math.Pow(cmplx.Abs(temp2.At(0, 0)), 2) * Pt / P_n)

	return rate
}
*/
/*func (s *Simulation) MIMO_Rate(H, G mat.CDense, Theta mat.CDiagonal) float64 {
	var temp1 mat.CDense
	var temp2 mat.CDense
	rate := 0.0

	temp1.Mul(G.T(), Theta)
	temp2.Mul(&temp1, &H)
	temp1.Mul(temp2.H(), &temp2)
	for i := 0; i < temp1.RawCMatrix().Rows; i++ {
		temp1.Set(i, i, temp1.At(i, i)+complex(1, 0))
	}
	temp1.Scale(complex(Pt/P_n, 0), &temp1)
	var lu cLU
	lu.Factorize(temp1)
	rate = lu.Det()

	return rate

}*/

func destructure(m cmat.Cmatrix) []float64 {

	v := make([]float64, m.Row*m.Col*2)

	for x := 0; x < m.Row; x++ {
		for y := 0; y < m.Col; y++ {
			index := (m.Col*x + y) * 2
			v[index] = real(m.Data[x][y])
			v[index+1] = imag(m.Data[x][y])
		}
	}

	return v
}
func construct(bytes ...[]float64) []float64 {
	var ret []float64
	for _, b := range bytes {
		ret = append(ret, b...)
	}

	return ret
}

/*
func EvaluateCoeffs(nbr int64, coef []float64) {

	fd, err := os.Open("SNR.csv")
	if err != nil {
		log.Print(err)
	}
	//csvwriter := csv.NewWriter(fd)

	var simCoef []complex128
	var CalculatedCoef []complex128
	var RandomCoef []complex128
	for _, value := range coef {
		CalculatedCoef = append(CalculatedCoef, cmplx.Rect(1, value))
		RandomCoef = append(RandomCoef, cmplx.Rect(1, rand.Float64()*2*math.Pi))

	}

	h := cmat.Transpose(SavedHG[nbr][0]).Data[0]

	g := SavedHG[nbr][1].Data[0]
	fmt.Println("Xapp Phases: ", coef)
	simCoef = GetCoefficients(h, g)
	//coefWOR := GetCoefficients_without_remainder(h, g)

	rate := RateSISO(h, g, CalculatedCoef)
	rate2 := RateSISO(h, g, RandomCoef)
	//rate3 := RateSISO(h, g, simCoef)
	//rate4 := RateSISO(h, g, coefWOR)

	fmt.Println("Coefficients From xApp : ")
	fmt.Println(CalculatedCoef[4:], "......")
	fmt.Println("Coefficients From SIM : ")
	fmt.Println(simCoef[4:], "......")

	//row := []string{strconv.FormatFloat(rate, 'f', -1, 64), strconv.FormatFloat(rate2, 'f', -1, 64), strconv.FormatFloat(rate3, 'f', -1, 64)}
	fmt.Println("Rate with Optimal Coefficients :", rate)
	//fmt.Println("Rate without remainder Coefficients :", rate4)
	fmt.Println("Rate with Random Coefficients :", rate2)
	//fmt.Println("Rate with Sim Coefficients :", rate3)
	//csvwriter.Write(row)
	//csvwriter.Flush()
	fd.Close()

}
*/
func RateSISO(H, G, Theta []complex128) float64 {

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

func SNRSISO(H_channel []complex128, G_channel []complex128, RIS_Coefficients []complex128, P_t, P_n float64) (float64, float64) {
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

func GetCoefficients(H, G []complex128) []complex128 {

	Theta_ris := []complex128{}
	fmt.Println("Sim Phases: ")
	for i := 0; i < len(H); i++ {
		phi_n := cmplx.Phase(H[i])
		psi_n := cmplx.Phase(G[i])
		fmt.Print(math.Remainder(-(phi_n+psi_n), 2*math.Pi), ", ")
		Theta_ris = append(Theta_ris, cmplx.Rect(1, math.Remainder(-(phi_n+psi_n), 2*math.Pi)))
	}

	fmt.Println()
	return Theta_ris
}

func GetCoefficients_without_remainder(H, G []complex128) []complex128 {

	Theta_ris := []complex128{}
	fmt.Println("Sim Phases: ")
	for i := 0; i < len(H); i++ {
		phi_n := cmplx.Phase(H[i])
		psi_n := cmplx.Phase(G[i])
		fmt.Print(math.Remainder(-(phi_n+psi_n), 2*math.Pi), ", ")
		Theta_ris = append(Theta_ris, cmplx.Rect(1, -(phi_n+psi_n)))
	}

	fmt.Println()
	return Theta_ris
}
