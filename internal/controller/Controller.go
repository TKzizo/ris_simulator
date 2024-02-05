package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	. "gitlab.eurecom.fr/ris-simulator/internal/structures"
	//. "gitlab.eurecom.fr/ris-simulator/internal/utils"
)

/*func GetCoefficients(H, G mat.CDense) mat.CDense { //SISO ie: H and G are column vector of the same size
	r, _ := H.Dims()
	Theta_ris := *mat.NewCDense(r, r, nil)
	for i := 0; i < r; i++ {
		phi_n := cmplx.Phase(H.At(i, 0))
		psi_n := cmplx.Phase(G.At(i, 0))
		Theta_ris.Set(i, i, cmplx.Rect(1, math.Remainder(-(phi_n+psi_n), 2*math.Pi)))
	}
	return Theta_ris
}

func GetCoefficients(H, G mat.CDense) mat.CDiagonal { //SISO ie: H and G are column vector of the same size
	r, _ := H.Dims()
	Theta_ris := mat.NewDiagCDense(r,nil)
	for i := 0; i < r; i++ {
		phi_n := cmplx.Phase(H.At(i, 0))
		psi_n := cmplx.Phase(G.At(i, 0))
		Theta_ris.SetDiag(i, cmplx.Rect(1, math.Remainder(-(phi_n+psi_n), 2*math.Pi)))
	}
	return Theta_ris
}*/

type AgentToRIC struct {
	Equipment string    `json:"Equipment"`
	TS        int       `json:"TS"`
	Positions []float64 `json:"Positions"`
	Elements  int       `json:"Elements"`
}

type RICToAgent struct {
	TS           int       `json:"TS"`
	Coefficients []float64 `json:"Coefficients"`
}

func ConnHandler(socket net.Listener, channl chan SimAgentChannel, bufferSize int) {
	log.Println("socket init")
	log.Println("waiting for connection at: ", socket.Addr())
	for {
		// Accept an incoming connection.
		conn, err := socket.Accept()
		if err != nil {
			log.Fatal("Failed to Establish socket Connection", err)
		}

		// Handle the connection in a separate goroutine.
		go func(conn net.Conn, channl chan SimAgentChannel) {
			defer conn.Close()

			buf := make([]byte, bufferSize)

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

				//go EvaluateCoeffs(int64(coef.TS), coef.Coefficients)
				//fmt.Println(coef.Coefficients)
			}

			var send string
			v := <-channl

			ris := AgentToRIC{
				Equipment: "Ris",
				Positions: v["RIS"][1:],
				Elements:  int(v["RIS"][0]),
				TS:        int(v["TS"][0]),
			}
			marshaled_ris, err := json.Marshal(ris)
			if err != nil {
				log.Print(err)
			}
			tx := AgentToRIC{
				Equipment: "Tx",
				Positions: v["TX"][1:],
				Elements:  int(v["TX"][0]),
				TS:        int(v["TS"][0]),
			}
			marshaled_tx, err := json.Marshal(tx)
			if err != nil {
				log.Print(err)
			}
			rx := AgentToRIC{
				Equipment: "Rx",
				Positions: v["RX"][1:],
				Elements:  int(v["RX"][0]),
				TS:        int(v["TS"][0]),
			}
			marshaled_rx, err := json.Marshal(rx)
			if err != nil {
				log.Print(err)
			}
			send = send + string(marshaled_ris) + "\n" + string(marshaled_tx) + "\n" + string(marshaled_rx) + "\n"

			/*for idx, f := range v[1:] {
				msg := AgentToRIC{
					Equipment: agent,
					Field:     Fields[idx],
					Data:      f, TS: int64(ts),
				}
				marshaled_msg, err := json.Marshal(msg)
				if err != nil {
					log.Print(err)
				}
				send = send + string(marshaled_msg) + "\n"
			}*/
			fmt.Println(send)
			//fmt.Println("Generated Channels : ")
			//fmt.Println(send[100:], "....")
			_, _ = conn.Write([]byte(send))
			//_, _ = conn.Write([]byte{1, 3, 5})
		}(conn, channl)
	}
}
