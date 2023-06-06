package main

import (
	cmat "RIS_SIMULATOR/reducedComplex"
	"log"
	"math"
	"net"
	"os"
	"strconv"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distuv"
)

type Coordinates struct {
	x float64
	y float64
	z float64
}

func Distance(a, b Coordinates) float64 {
	return math.Sqrt(math.Pow(a.x-b.x, 2) + math.Pow(a.y-b.y, 2) + math.Pow(a.z-b.z, 2))
}
func RadToDeg(a float64) float64 {
	return a * 180 / math.Pi
}
func DegToRad(a float64) float64 {
	return a * math.Pi / 180
}
func sign(a, b float64) int8 {
	if a >= b {
		return 1
	}
	return -1
}
func Ge(theta float64) float64 {
	return Gain * math.Pow(math.Cos(theta), 2*q)
}
func L(s *Simulation, sf distuv.Normal, los bool, a ...Coordinates) float64 {
	if los == true {
		return math.Pow(10, (-20*math.Log10(4*math.Pi/s.Lambda)-10*s.n_LOS*(1+s.b_LOS*((s.Frequency-s.f0)/s.f0))*math.Log10(Distance(a[0], a[1]))-(rand.Float64()*s.sigma_LOS))/10)
	} else {
		return math.Pow(10, (-20*math.Log10(4*math.Pi/s.Lambda)-10*s.n_NLOS*(1+s.b_NLOS*((s.Frequency-s.f0)/s.f0))*math.Log10(Distance(a[0], a[1])+Distance(a[1], a[2]))-(rand.Float64()*s.sigma_NLOS))/10)
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

func complextoString(c complex128) string {
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

func setupSocket(addr string) {
	var SockAddr string = addr

	os.Remove(SockAddr)
	socket, err := net.Listen("unix", SockAddr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		// Accept an incoming connection.
		conn, err := socket.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// Handle the connection in a separate goroutine.
		go func(conn net.Conn) {
			defer conn.Close()
			// Create a buffer for incoming data.
			buf := make([]byte, 4096)

			// Read data from the connection.
			for {
				n, err := conn.Read(buf)
				if err != nil {
					log.Fatal(err)
				}
				// WE NEED TO GENERATE BIG CHUNKS OF DATA AND SEND IT
				// BEFORE THAT WE NEED TO CHANGE ITS TYPE FROM COMPLEX TO BYTE ??
				// AND CONFIRM RECEIVING IT
				// Echo the data back to the connection.
				_, err = conn.Write(buf[:n])
				if err != nil {
					log.Fatal(err)
				}
			}
		}(conn)
	}

}
