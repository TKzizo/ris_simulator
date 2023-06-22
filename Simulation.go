package main

import (
	cmat "RIS_SIMULATOR/reducedComplex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"sync"
	"time"
)

const (
	q    float64 = 0.283 // related to the gain
	Gain float64 = math.Pi
	Pt   float64 = 0.9         // Power of transmitter
	P_n  float64 = 0.000000001 // variance of noise at the receiver
)

type AgentToRIC struct {
	Equipment string    `json:"Equipment"`
	Field     string    `json:"Field`
	Data      []float64 `json:"Data"`
}

type RICToAgent struct {
	Coefficients []float64 `json:"Coefficients`
}

type RISCHANNL [][]float64

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
	RisChannl chan RISCHANNL
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
	s.RisChannl /* s.TxChannl, s.RxChannl*/ = s.setupSockets()
	s.InputPositions()
	//s.CheckPositioning() // To apply the 3GPP standards
}

func (s *Simulation) Run() []cmat.Cmatrix {
	var h, g cmat.Cmatrix
	list := []cmat.Cmatrix{}
	// Re-run the calculation for every position of the user
	for _, update := range s.Positions {
		clusters := GenerateClusters(s)
		s.Ris.xyz = update.ris
		s.Tx.xyz = update.tx
		s.Rx.xyz = update.rx
		h = s.H_channel(clusters)
		list = append(list, h)
		g = s.G_channel()
		list = append(list, g)
	}
	return list

}

func (s *Simulation) setupSockets() chan RISCHANNL /*, chan []float64, chan []float64*/ {
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
	/*socketTX, err := net.Listen("unix", txaddr)
	if err != nil {
		log.Fatal(err)
	}
	socketRX, err := net.Listen("unix", rxaddr)
	if err != nil {
		log.Fatal(err)
	}*/
	risChannl := make(chan RISCHANNL)
	//txChannl := make(chan []float64)
	//rxChannl := make(chan []float64)

	go connHandler(socketRIS, "RIS", risChannl)
	//go connHandler(socketTX, "TX", txChannl)
	//go connHandler(socketRX, "RX", rxChannl)

	return risChannl //, txChannl, rxChannl
}

func connHandler(socket net.Listener, agent string, channl chan RISCHANNL) {
	var NumberOfReads int = 0
	var mutex sync.Mutex
	for {
		// Accept an incoming connection.
		conn, err := socket.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// Handle the connection in a separate goroutine.
		go func(conn net.Conn, channl chan RISCHANNL) {
			defer conn.Close()

			buf := make([]byte, 256*2*20) // 256 patch * real x imag * number of bytes

			err := conn.SetReadDeadline(time.Now().Add(40 * time.Millisecond))
			if err != nil {
				log.Fatal(err)
			}
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println(err, "Reading Ris Coefficients")
			}
			if n != 0 {
				//fmt.Print(string(buf))
				coef := RICToAgent{}
				if err := json.Unmarshal(buf[:n], &coef); err != nil {
					log.Print(err, "received: ", n)
				}

				go EvaluateCoeffs(NumberOfReads, coef.Coefficients)
				mutex.Lock()
				NumberOfReads++
				mutex.Unlock()
				//fmt.Println(coef.Coefficients)
			}

			var send string
			Fields := []string{"Position", "H", "G"}

			v := <-channl
			for idx, f := range v {
				msg := AgentToRIC{Equipment: agent, Field: Fields[idx], Data: f}
				marshaled_msg, err := json.Marshal(msg)
				if err != nil {
					log.Print(err)
				}
				send = send + string(marshaled_msg) + "\n"
			}
			//fmt.Println("Generated Channels : ")
			//fmt.Println(send[100:], "....")
			_, _ = conn.Write([]byte(send))
			//_, _ = conn.Write([]byte{1, 3, 5})
		}(conn, channl)
	}
}
