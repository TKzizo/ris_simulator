package main

import (
	cmat "RIS_SIMULATOR/reducedComplex"
	"log"
	"math"
	"net"
	"os"
)

const (
	q    float64 = 0.283 // related to the gain
	Gain float64 = math.Pi
	Pt   float64 = 0.05        // Power of transmitter
	P_n  float64 = 0.000000001 // variance of noise at the receiver
)

type Updates struct {
	ris Coordinates
	rx  Coordinates
	tx  Coordinates
	los bool
}

type Simulation struct {
	Env        Environment
	Ris        RIS
	Tx         Tx_Rx
	Rx         Tx_Rx
	Frequency  float64
	f0         float64
	Lambda_p   float64
	Lambda     float64 // wave length
	k          float64
	n_LOS      float64 // Path Loss exponent
	b_LOS      float64 // systemc parameter
	sigma_LOS  float64 // db
	n_NLOS     float64
	b_NLOS     float64
	sigma_NLOS float64
	//channel    chan Updates
	Broadside int8 // 0: SideWall 1: OppositeWall
	Positions []Updates
}

func (s *Simulation) Setup() {
	s.Lambda = 3.0 / (10 * s.Frequency) // it's Simplified so it only supports GHz
	s.k = 2 * math.Pi / s.Lambda

	if s.Frequency == 28.0 {
		s.Lambda_p = 1.8
	} else if s.Frequency == 73.0 {
		s.Lambda_p = 1.9
	}

	if s.f0 == 0.0 {
		s.f0 = 24.2
	}

	if s.n_LOS == 0.0 { //Pathloss exponent
		s.n_LOS = 1.73
	}

	if s.n_NLOS == 0.0 {
		s.n_NLOS = 3.79
	}

	if s.b_NLOS == 0.0 {
		s.b_NLOS = 3.19
	}

	if s.sigma_LOS == 0.0 {
		s.sigma_LOS = 3.02
	}

	if s.sigma_NLOS == 0.0 {
		s.sigma_NLOS = 8.29
	}

	s.Ris.Setup(s.Lambda)
	s.Rx.Setup(s.Lambda)
	s.Tx.Setup(s.Lambda)

	s.InputPositions()
	//s.CheckPositioning() // To apply the 3GPP standards
}

func (s *Simulation) Run() (*cmat.Cmatrix, *cmat.Cmatrix) {
	var h, g cmat.Cmatrix
	//list := []cmat.Cmatrix{h, g}
	// Re-run the calculation for every position of the user
	for _, update := range s.Positions {
		clusters := GenerateClusters(s)
		s.Ris.xyz = update.ris
		s.Tx.xyz = update.tx
		s.Rx.xyz = update.rx
		h = s.H_channel(clusters)
		//list = append(list, h)
		g = s.G_channel()
		//list = append(list, g)
	}
	return &h, &g

}

func (s *Simulation) setupSockets() {
	var risaddr string = "/tmp/ris.sock"
	var txaddr string = "/tmp/tx.sock"
	var rxaddr string = "/tmp/rx.sock"

	// Remove socket if it already exists
	os.Remove(risaddr)
	os.Remove(txaddr)
	os.Remove(rxaddr)

	// Create Listner for each socket
	socketRIS, err := net.Listen("unix", risaddr)
	if err != nil {
		log.Fatal(err)
	}
	socketTX, err := net.Listen("unix", txaddr)
	if err != nil {
		log.Fatal(err)
	}
	socketRX, err := net.Listen("unix", rxaddr)
	if err != nil {
		log.Fatal(err)
	}

	go connHandler(socketRIS)
	go connHandler(socketTX)
	go connHandler(socketRX)
}

func connHandler(socket net.Listener) {

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
				// Echo the data back to the connection.
				_, err = conn.Write(buf[:n])
				if err != nil {
					log.Fatal(err)
				}
			}
		}(conn)
	}
}
