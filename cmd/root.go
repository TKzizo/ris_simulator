/*
Copyright © 2023  <adlen.ksentini@eurecom.fr>
*/
package cmd

import (
	"os"

	generate "gitlab.eurecom.fr/ris-simulator/cmd/generate"
	simulate "gitlab.eurecom.fr/ris-simulator/cmd/simulate"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "toolbox",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}

	}}

var (
	cfgFile       string
	logDir        string
	logPrefix     string
	userPositions string
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(generate.GenerateCmd)
	rootCmd.AddCommand(simulate.SimulateCmd)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config-file", ".init.json", "config file path.")
	rootCmd.PersistentFlags().StringVar(&logDir, "log-dir", "/tmp/", "logging directory path.")
	rootCmd.PersistentFlags().StringVar(&logPrefix, "log-prefix", "RIS_SIMULATION_", "log file prefix.")
	rootCmd.PersistentFlags().StringVar(&userPositions, "rx-positions-file", ".positions.csv", "path to file containing user positions through the simulation.")

}

/*
var SavedHG map[int64][]cmat.Cmatrix = make(map[int64][]cmat.Cmatrix)

func _main() {

	list := simulation.Run()
	_, _ = os.Create("SNR.csv")

	for {
		for i, v := range simulation.Positions {
			h := list[i*2]
			g := list[i*2+1]
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
*/
