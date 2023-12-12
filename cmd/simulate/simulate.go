/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package simulate

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	portDefault int16 = 8080
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

		cfgFile = cmd.Parent().Flags().Lookup("config-file").Value.String()

	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("simulate called")
		fmt.Println("Gui Bool", gui)
		fmt.Println("Gui port", port)
		fmt.Println(cmd.Flags().Lookup("set-gui-port").Changed)
	},
}

var (
	cfgFile string
	gui     bool
	port    int16
)

func init() {

	SimulateCmd.Flags().BoolVar(&gui, "set-gui", false, "Visualize the simulation inside browser")
	SimulateCmd.Flags().Int16Var(&port, "set-gui-port", portDefault, "set port of web visualization page")

	fmt.Println("simulate init call")

}
