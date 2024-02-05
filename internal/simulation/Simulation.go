package simulation

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"unsafe"

	. "gitlab.eurecom.fr/ris-simulator/internal/controller"
	cmat "gitlab.eurecom.fr/ris-simulator/internal/reducedComplex"
	. "gitlab.eurecom.fr/ris-simulator/internal/structures"
	. "gitlab.eurecom.fr/ris-simulator/internal/utils"
)

type Simulation struct {
	Env        Environment
	Ris        RIS
	Tx         Tx_Rx
	Rx         Tx_Rx
	Frequency  float64
	F0         float64
	Lambda_p   float64
	Lambda     float64 // wave length
	K          float64
	N_LOS      float64 // Path Loss exponent
	B_LOS      float64 // systemc parameter
	Sigma_LOS  float64 // db
	N_NLOS     float64
	B_NLOS     float64
	Sigma_NLOS float64
	//channel    chan Updates
	Broadside int8 // 0: SideWall 1: OppositeWall
	Positions []Updates
	RisChannl chan SimAgentChannel
}

func (s *Simulation) Setup(cfg InitConfig, rxPositions string) {

	s.Ris.Setup(s.Lambda)
	s.Rx.Setup(s.Lambda)
	s.Tx.Setup(s.Lambda)

	s.Lambda = 3.0 / (10 * s.Frequency) // it's Simplified so it only supports GHz
	s.K = 2 * math.Pi / s.Lambda

	if s.Frequency == 28.0 {
		s.Lambda_p = 1.8
	} else if s.Frequency == 73.0 {
		s.Lambda_p = 1.9
	}

	if s.F0 == 0.0 {
		s.F0 = 24.2
	}

	if s.N_LOS == 0.0 { //Pathloss exponent
		s.N_LOS = 1.73
	}

	if s.N_NLOS == 0.0 {
		s.N_NLOS = 3.79
	}

	if s.B_NLOS == 0.0 {
		s.B_NLOS = 3.19
	}

	if s.Sigma_LOS == 0.0 {
		s.Sigma_LOS = 3.02
	}

	if s.Sigma_NLOS == 0.0 {
		s.Sigma_NLOS = 8.29
	}

	//s.RisChannl /* s.TxChannl, s.RxChannl*/ = s.setupSockets()
	s.SetupSockets()
	s.InputPositions(rxPositions)
	//s.CheckPositioning() // To apply the 3GPP standards
}
func (s *Simulation) Run() *[]cmat.Cmatrix {
	var h, g, d cmat.Cmatrix
	list := []cmat.Cmatrix{}
	// Re-run the calculation for every position of the user
	for _, update := range s.Positions {
		clusters := GenerateClusters(s)
		s.Ris.Xyz = update.Ris
		s.Tx.Xyz = update.Tx
		s.Rx.Xyz = update.Rx
		h = s.H_channel(clusters)
		list = append(list, h)
		g = s.G_channel()
		list = append(list, g)
		if update.Los {
			d = s.D_channel(clusters)
			//fmt.Println(d)
		} else {
			d.Init(s.Rx.N, s.Tx.N)
			//fmt.Println(d)
		}
		list = append(list, d)
	}
	return &list
}
func (s *Simulation) InputPositions(filePath string) {
	csvFile, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
	}
	log.Println("Position file opened successfully")

	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	log.Println("Number of User positions: ", len(csvLines))
	list_positions := []Updates{}

	for _, line := range csvLines[1:] {
		position := Updates{}
		if v, err := strconv.ParseFloat(line[0], 64); err == nil && v <= s.Env.Length && v >= 0 {
			position.Ris.X = v
		} else {
			fmt.Println("Position Over-Boundries Length Ris")
			continue
		}
		if v, err := strconv.ParseFloat(line[1], 64); err == nil && v <= s.Env.Width && v >= 0 {
			position.Ris.Y = v
		} else {
			fmt.Println("Position Over-Boundries width Ris")
			continue
		}

		if v, err := strconv.ParseFloat(line[2], 64); err == nil && v <= s.Env.Height && v >= 0 {
			position.Ris.Z = v
		} else {
			fmt.Println("Position Over-Boundries height Ris")
			continue
		}

		if v, err := strconv.ParseFloat(line[3], 64); err == nil && v <= s.Env.Length && v >= 0 {
			position.Tx.X = v
		} else {
			fmt.Println("Position Over-Boundries Length Tx")
			continue
		}

		if v, err := strconv.ParseFloat(line[4], 64); err == nil && v <= s.Env.Width && v >= 0 {
			position.Tx.Y = v
		} else {
			fmt.Println("Position Over-Boundries width Tx")
			continue
		}

		if v, err := strconv.ParseFloat(line[5], 64); err == nil && v <= s.Env.Height && v >= 0 {
			position.Tx.Z = v
		} else {
			fmt.Println("Position Over-Boundries height Tx")
			continue
		}

		if v, err := strconv.ParseFloat(line[6], 64); err == nil && v <= s.Env.Length && v >= 0 {
			position.Rx.X = v
		} else {
			fmt.Println("Position Over-Boundries Length Rx")
			continue
		}

		if v, err := strconv.ParseFloat(line[7], 64); err == nil && v <= s.Env.Width && v >= 0 {
			position.Rx.Y = v
		} else {
			fmt.Println("Position Over-Boundries width Rx")
			continue
		}

		if v, err := strconv.ParseFloat(line[8], 64); err == nil && v <= s.Env.Height && v >= 0 {
			position.Rx.Z = v
		} else {
			fmt.Println("Position Over-Boundries height Rx")
			continue
		}

		if v, err := strconv.ParseFloat(line[9], 64); err == nil {
			if v == 1 {
				position.Los = true
			} else {
				position.Los = false
			}
		}
		list_positions = append(list_positions, position)
	}
	/*for _, v := range list_positions {
		fmt.Println(v)
	}*/
	s.Positions = list_positions

}

func (s *Simulation) SetupSockets() /*chan RISCHANNL, chan []float64, chan []float64*/ {
	var risaddr string = "/tmp/ris.sock"

	// Remove socket if it already exists
	err := os.Remove(risaddr)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatal("Could not remove existing Socket: ", err)
		}
	}

	// Create Socket Listner
	socketRIS, err := net.Listen("unix", risaddr)
	if err != nil {
		log.Fatal("Could not create socket listner", err)
	}

	s.RisChannl = make(chan SimAgentChannel)
	bufferSize := s.Ris.N * int(unsafe.Sizeof(0.0)) * 2 //  nbr_patches * real x imag  * number of bytes
	if bufferSize <= 0 {
		log.Fatal("buffer size error.")
	}

	go ConnHandler(socketRIS, s.RisChannl, bufferSize)

	//return risChannl //, txChannl, rxChannl
	//return nil
}

func InitSimualtion(cfg InitConfig) *Simulation {
	ris := RIS{
		N: cfg.Equipements.Ris.Elements,
		Xyz: Coordinates{
			X: cfg.Equipements.Ris.Coord.X,
			Y: cfg.Equipements.Ris.Coord.Y,
			Z: cfg.Equipements.Ris.Coord.Z,
		},
	}

	tx := Tx_Rx{
		N: cfg.Equipements.Tx.Elements,
		Xyz: Coordinates{
			X: cfg.Equipements.Tx.Coord.X,
			Y: cfg.Equipements.Tx.Coord.Y,
			Z: cfg.Equipements.Tx.Coord.Z,
		},
		Type: cfg.Equipements.Tx.Type,
	}

	rx := Tx_Rx{
		N: cfg.Equipements.Rx.Elements,
		Xyz: Coordinates{
			X: cfg.Equipements.Rx.Coord.X,
			Y: cfg.Equipements.Rx.Coord.Y,
			Z: cfg.Equipements.Rx.Coord.Z,
		},
		Type: cfg.Equipements.Rx.Type,
	}

	simulation := Simulation{
		Ris:       ris,
		Tx:        tx,
		Rx:        rx,
		Frequency: cfg.Frequency,
		Env: Environment{
			Length: cfg.Env.Length,
			Width:  cfg.Env.Width,
			Height: cfg.Env.Height},

		Broadside: 0,
	}

	return &simulation
}
