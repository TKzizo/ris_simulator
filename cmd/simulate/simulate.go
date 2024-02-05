/*
Copyright © 2023 Eurecom <adlen.ksentini@eurecom.fr>
*/
package simulate

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	cmat "gitlab.eurecom.fr/ris-simulator/internal/reducedComplex"
	. "gitlab.eurecom.fr/ris-simulator/internal/simulation"
	"gitlab.eurecom.fr/ris-simulator/internal/utils"
)

const (
	portDefault int = 9696
)

var (
	cfgFilePath   string
	cfg           utils.InitConfig
	userPositions string
	gui           bool
	//port          int
	simulation *Simulation
	//testMode  bool // This mode is for testing the xAPP
)

// simulateCmd represents the simulate command
var SimulateCmd = &cobra.Command{
	Use:   "simulate",
	Short: "Launch simulation",
	Long: `Description:
	launch the simulateur and connect to E2 agent to receive coefficients from Xapp`,
	PreRun: func(cmd *cobra.Command, args []string) {
		// This is a fix to having the port value set without setting gui to true
		/*{
			port := cmd.Flags().Lookup("set-viz-port")
			gui := cmd.Flags().Lookup("set-viz")
			if port.Changed && !gui.Changed {
				gui.Value.Set("true")
			}
		}*/

		cfgFilePath = cmd.Parent().Flags().Lookup("config-file").Value.String()
		userPositions = cmd.Parent().Flags().Lookup("rx-positions-file").Value.String()
		cfg = utils.InitCfg(cfgFilePath)

	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("simulate called")
		fmt.Println("Gui Bool: ", gui)
		//fmt.Println("Gui port: ", port)
		//fmt.Println(cmd.Flags().Lookup("set-viz-port").Changed)
		simulation = InitSimualtion(cfg)
		simulation.Setup(cfg, userPositions)
		if gui {
			fmt.Println("Visualization at: http://localhost:" + strconv.Itoa(portDefault))
			http.HandleFunc("/", mainHandler)
			http.HandleFunc("/pos", PostionHandler)
			http.HandleFunc("/init", initHandler)
			http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("web/assets"))))
			http.Handle("/loadtest/", http.StripPrefix("/loadtest/", http.FileServer(http.Dir("web/assets/loadtest"))))
			http.ListenAndServe(":"+strconv.Itoa(portDefault), nil)
		}
		simulate(*simulation)

	},
}

func init() {

	SimulateCmd.Flags().BoolVar(&gui, "set-viz", false, "Visualize the simulation inside browser")
	//SimulateCmd.Flags().IntVar(&port, "set-viz-port", portDefault, "set port of web visualization page") // Removed

	fmt.Println("simulate init call")

}

func simulate(simulation Simulation) {

	var SavedHG map[int64][]*cmat.Cmatrix = make(map[int64][]*cmat.Cmatrix, 10) // varilable was used to save the scheduling until: "the xAPP get fixed to send the coefficients inorder"

	//list := simulation.Run()
	//_, _ = os.Create("SNR.csv")

	for {
		list := simulation.Run() // channel may change ever slightly from (p1,t0) to (p1,t1)
		for i, v := range simulation.Positions {
			// List: [h0, g0, d0, ... ,h,g,d]
			h := &(*list)[i*3]
			g := &(*list)[i*3+1]
			d := &(*list)[i*3+2] // d is empty matrix if v.LOS == false

			//hd := utils.Destructure(h)
			//gd := utils.Destructure(g)
			//dd := utils.Destructure(d)
			//	simulation.RisChannl <- construct([]float64{simulation.Ris.xyz.x, simulation.Ris.xyz.y, simulation.Ris.xyz.z}, hd, gd)
			//	simulation.RisChannl <- []float64{simulation.Tx.xyz.x, simulation.Rx.xyz.y, simulation.Rx.xyz.z}
			ts := time.Now().Unix()
			//simulation.RisChannl <- [][]float64{[]float64{float64(ts)}, []float64{v.rx.x, v.rx.y, v.rx.z}, hd, gd, dd}
			simulation.RisChannl <- map[string][]float64{
				"RIS": []float64{float64(simulation.Ris.N), v.Ris.X, v.Ris.Y, v.Ris.Z},
				"TX":  []float64{float64(simulation.Tx.N), v.Tx.X, v.Tx.Y, v.Tx.Z},
				"RX":  []float64{float64(simulation.Rx.N), v.Rx.X, v.Rx.Y, v.Rx.Z},
				"TS":  []float64{float64(ts)},
			}
			SavedHG[ts] = []*cmat.Cmatrix{h, g, d}
			time.Sleep(20 * time.Millisecond)
			//generateData(simulation, 1)
		}
	}

	time.Sleep(10 * time.Millisecond)
}
