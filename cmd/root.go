/*
Copyright © 2023 Eurecom <adlen.ksentini@eurecom.fr>
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
