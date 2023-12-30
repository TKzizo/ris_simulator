/*
Copyright © 2023 Eurecom <adlen.ksentini@eurecom.fr>
*/
package simulate

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	. "gitlab.eurecom.fr/ris-simulator/internal/simulation"
	"gitlab.eurecom.fr/ris-simulator/internal/utils"
)

const (
	portDefault int16 = 8080
)

var (
	cfgFilePath   string
	cfg           utils.InitConfig
	userPositions string
	gui           bool
	port          int16
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
		{
			port := cmd.Flags().Lookup("set-gui-port")
			gui := cmd.Flags().Lookup("set-gui")
			if port.Changed && !gui.Changed {
				gui.Value.Set("true")
			}
		}

		cfgFilePath = cmd.Parent().Flags().Lookup("config-file").Value.String()
		userPositions = cmd.Parent().Flags().Lookup("rx-positions-file").Value.String()
		cfg = utils.InitCfg(cfgFilePath)

	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("simulate called")
		fmt.Println("Gui Bool: ", gui)
		fmt.Println("Gui port: ", port)
		fmt.Println(cmd.Flags().Lookup("set-gui-port").Changed)
		simulation := InitSimualtion(cfg)
		simulation.Setup(cfg, userPositions)
		simulate(*simulation)

	},
}

func init() {

	SimulateCmd.Flags().BoolVar(&gui, "set-gui", false, "Visualize the simulation inside browser")
	SimulateCmd.Flags().Int16Var(&port, "set-gui-port", portDefault, "set port of web visualization page")

	fmt.Println("simulate init call")

}

var SavedHG map[int64][]cmat.Cmatrix = make(map[int64][]cmat.Cmatrix)

func simulate(simulation Simulation) {

	list := simulation.Run()
	//_, _ = os.Create("SNR.csv")

	for {
		for i, v := range simulation.Positions {
			// List: [h0, g0, d0, ... ,h,g,d]
			h := list[i*3]
			g := list[i*3+1]
			hd := destructure(h)
			gd := destructure(g)
			//	simulation.RisChannl <- construct([]float64{simulation.Ris.xyz.x, simulation.Ris.xyz.y, simulation.Ris.xyz.z}, hd, gd)
			//	simulation.RisChannl <- []float64{simulation.Tx.xyz.x, simulation.Rx.xyz.y, simulation.Rx.xyz.z}
			ts := time.Now().Unix()
			simulation.RisChannl <- [][]float64{[]float64{float64(ts)}, []float64{v.rx.x, v.rx.y, v.rx.z}, hd, gd}
			SavedHG[ts] = []cmat.Cmatrix{h, g}
			time.Sleep(3 * time.Second)
			//generateData(simulation, 1)
		}
	}

	time.Sleep(10 * time.Second)
}
