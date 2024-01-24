package generate

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	. "gitlab.eurecom.fr/ris-simulator/internal/simulation"
	"gitlab.eurecom.fr/ris-simulator/internal/structures"
	. "gitlab.eurecom.fr/ris-simulator/internal/utils"
)

const configFormat string = `
{
"NumberRxPositions":%d,
"NumberIterations":%d,
"Positions":
	{
		"Ris":%v,
		"Tx":%v,
		"Rx":%v
	}
}
`

var (
	h, g, d    [][]string
	hb, gb, db bool
)

func generateData(simulation *Simulation) {

	if err := DefaultDir(outputDir); err != nil {
		panic(err)
	}
	timeNow := time.Now().UTC().String()
	directoryNameFormat := outputDir + "Simulation_" + timeNow

	if err := CreateDir(directoryNameFormat); err != nil {
		panic(err)
	}

	nbr_positions := len(simulation.Positions)

	saveRunningConfig(directoryNameFormat, nbr_positions, simulation)
	//setting up Data strctures and directories to save results
	for _, v := range channelsList {
		switch v {
		case "H":
			hb = true
			h = make([][]string, simulation.Ris.N*simulation.Tx.N)
			for i := 0; i < len(h); i++ {
				h[i] = make([]string, iterations)
			}
			if err := CreateDir(directoryNameFormat + "/H_Channel"); err != nil {
				panic(err)
			}
		case "G":
			gb = true
			g = make([][]string, simulation.Ris.N*simulation.Rx.N)
			for i := 0; i < len(g); i++ {
				g[i] = make([]string, iterations)
			}
			if err := CreateDir(directoryNameFormat + "/G_Channel"); err != nil {
				panic(err)
			}
		case "D":
			db = true
			d = make([][]string, simulation.Tx.N*simulation.Rx.N)
			for i := 0; i < len(d); i++ {
				d[i] = make([]string, iterations)
			}
			if err := CreateDir(directoryNameFormat + "/D_Channel"); err != nil {
				panic(err)
			}
		}
	}

	for y := 0; y < iterations; y++ {
		list := simulation.Run()
		index_channels := 0
		if hb {
			for i := 0; i < nbr_positions; i++ {
				H := &(*list)[index_channels+i*3]
				for hi, line := range H.Data {
					for hii, ele := range line {
						h[hi*len(line)+hii][i] = ComplextoString(ele)
					}
				}
			}
			index_channels++
		}
		if gb {
			for i := 0; i < nbr_positions; i++ {
				G := &(*list)[index_channels+i*3]
				for gi, line := range G.Data {
					for gii, ele := range line {
						g[gi*len(line)+gii][i] = ComplextoString(ele)
					}
				}
			}
			index_channels++
		}
		if db {
			for i := 0; i < nbr_positions; i++ {

				D := &(*list)[index_channels+i*3]
				for di, line := range D.Data {
					for dii, ele := range line {
						d[di*len(line)+dii][i] = ComplextoString(ele)
					}
				}
			}
			index_channels++
		}

		fmt.Println(y, " iteration")

		if hb {
			path := directoryNameFormat + "/H_Channel/iteration_" + strconv.Itoa(y) + ".csv"
			saveToCsv(path, h)
		}
		if gb {
			path := directoryNameFormat + "/G_Channel/iteration_" + strconv.Itoa(y) + ".csv"
			saveToCsv(path, g)
		}
		if db {
			path := directoryNameFormat + "/D_Channel/iteration_" + strconv.Itoa(y) + ".csv"
			saveToCsv(path, d)
		}
	}
}

func saveToCsv(path string, channel [][]string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(file)
	for _, row := range channel {
		_ = csvwriter.Write(row)
	}
	csvwriter.Flush()
	file.Close()
}

func saveRunningConfig(path string, nbrPositions int, sim *Simulation) {

	file, err := os.Create(path + "/simulation_info.json")
	if err != nil {
		log.Print(err)
	}

	defer file.Close()

	pRis := fmt.Sprintf("%v", sim.Ris.Xyz)
	pTx := fmt.Sprintf("%v", sim.Tx.Xyz)
	pRx := fmt.Sprintf("%v", func(list []structures.Updates) []structures.Coordinates {
		result := []structures.Coordinates{}
		for _, v := range list {
			result = append(result, v.Rx)
		}
		return result
	}(sim.Positions))

	//	cfg = strings.ReplaceAll(cfg, " ", ",", -1)
	replacer := strings.NewReplacer(" ", ",", "{", "[", "}", "]")
	pRis = replacer.Replace(pRis)
	pTx = replacer.Replace(pTx)
	pRx = replacer.Replace(pRx)

	cfg := fmt.Sprintf(configFormat,
		nbrPositions,
		iterations,
		pRis,
		pTx,
		pRx,
	)

	file.Write([]byte(cfg))
}
