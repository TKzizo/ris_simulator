package main

import (
	cmat "RIS_SIMULATOR/reducedComplex"
	"encoding/json"
	"fmt"
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

type AgentToRIC struct {
	Type string `json:"Type"`
	// The first three values represent the position x,y,z of the RAN
	// if the Type is RIS then the following values would be the H,G values
	Data []float64 `json:"Data"`
}

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
	RisChannl chan []float64
	TxChannl  chan []float64
	RxChannl  chan []float64
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
	s.RisChannl, s.TxChannl, s.RxChannl = s.setupSockets()
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

func (s *Simulation) setupSockets() (chan []float64, chan []float64, chan []float64) {
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
	risChannl := make(chan []float64)
	txChannl := make(chan []float64)
	rxChannl := make(chan []float64)

	go connHandler(socketRIS, "RIS", risChannl)
	go connHandler(socketTX, "TX", txChannl)
	go connHandler(socketRX, "RX", rxChannl)

	return risChannl, txChannl, rxChannl
}

func connHandler(socket net.Listener, agent string, channl chan []float64) {

	for {
		// Accept an incoming connection.
		conn, err := socket.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// Handle the connection in a separate goroutine.
		go func(conn net.Conn, channl chan []float64) {
			defer conn.Close()
			// Create a buffer for incoming data.
			// The max size of 1 iteration of MIMO with 256 patch, 16 tx antennas and 64 rx antennas is currently 500KB
			// For the tx and Rx we just gonna need 30B for positionf and type map
			var buf []byte
			switch agent {
			case "RIS":
				buf = make([]byte, 1024*500)
			default:
				buf = make([]byte, 100)
			}
			// Write Data to connection.
			encoder := json.NewEncoder(conn)

			for {
				select {
				case v := <-channl:
					msg := AgentToRIC{Type: agent, Data: v}
					encoder.Encode(msg)
				}
				if buf[0] == byte(0) {
					fmt.Println("nothing received")
				}
			}
			/*for {
				n, err := conn.Read(buf)
				if err != nil {
					log.Fatal(err)
				}
				// Echo the data back to the connection.
				_, err = conn.Write(buf[:n])
				if err != nil {
					log.Fatal(err)
				}
			}*/
		}(conn, channl)
	}
}
