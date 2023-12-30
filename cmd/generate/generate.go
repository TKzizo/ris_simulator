/*
Copyright © 2023 Eurecom <adlen.ksentini@eurecom.fr>
*/
package generate

import (
	"fmt"

	"github.com/spf13/cobra"
	//. "gitlab.eurecom.fr/ris-simulator/internal/logger"
	. "gitlab.eurecom.fr/ris-simulator/internal/simulation"
	"gitlab.eurecom.fr/ris-simulator/internal/utils"
)

// generateCmd represents the generate command
var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate channels coefficients",
	Long: `Description:
	Generate the coefficients of channels H,G,D through a number of iterations
	and saves the result in a csv file.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if c := cmd.Flags().Lookup("channels"); c.Changed {
			channelsList = uniqListChannels(channelsList)
		}
		cfgFilePath = cmd.Parent().Flags().Lookup("config-file").Value.String()
		userPositions = cmd.Parent().Flags().Lookup("rx-positions-file").Value.String()
		cfg = utils.InitCfg(cfgFilePath)
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Generate command called")
		simulation := InitSimualtion(cfg)
		simulation.Setup(cfg, userPositions)
		generateData(simulation)
	},
}

var (
	cfgFilePath   string
	cfg           utils.InitConfig
	userPositions string
	channelsList  []string
	outputDir     string
	iterations    int
)

func init() {
	GenerateCmd.Flags().StringVar(&outputDir, "ouput-dir", "data/", "path to directory to save the results")
	GenerateCmd.Flags().IntVar(&iterations, "iterations", 1000, "path to directory to save the results")
	GenerateCmd.Flags().StringSliceVar(&channelsList, "channels", []string{"H", "G", "D"}, "Set of channel to be saved")
}

// This function makes sure that the list contains a set of H,G or D
func uniqListChannels(list []string) []string {
	//some bit juggling
	//1 H
	//2 G
	//4 D
	var bits int8 = 0
	var output []string
	for _, v := range channelsList {
		switch v {
		case "H":
			bits |= 1
		case "G":
			bits |= 2
		case "D":
			bits |= 4
		}
	}
	if (bits & 1) > 0 {
		output = append(output, "H")
	}
	if (bits & 2) > 0 {
		output = append(output, "G")
	}
	if (bits & 4) > 0 {
		output = append(output, "D")
	}

	return output
}
