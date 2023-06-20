package main

import (
	"time"
)

func main() {

	ris := RIS{N: 4, xyz: Coordinates{x: 40, y: 50, z: 2}}
	tx := Tx_Rx{N: 1, Type: 0, xyz: Coordinates{x: 0, y: 25, z: 2}}
	rx := Tx_Rx{N: 1, Type: 0, xyz: Coordinates{x: 38, y: 48, z: 1}}

	simulation := Simulation{
		Ris:       ris,
		Tx:        tx,
		Rx:        rx,
		Frequency: 28.0,
		Env:       Environment{75.0, 50.0, 3.5}}

	simulation.Setup()
	for {
		list := simulation.Run()
		for i, v := range simulation.Positions {
			h := list[i*2]
			g := list[i*2+1]
			hd := destructure(h)
			gd := destructure(g)

			//	simulation.RisChannl <- construct([]float64{simulation.Ris.xyz.x, simulation.Ris.xyz.y, simulation.Ris.xyz.z}, hd, gd)
			//	simulation.RisChannl <- []float64{simulation.Tx.xyz.x, simulation.Rx.xyz.y, simulation.Rx.xyz.z}
			simulation.RisChannl <- []float64{v.rx.x, v.rx.y, v.rx.z}
			simulation.RisChannl <- hd
			simulation.RisChannl <- gd
			time.Sleep(2 * time.Second)
			//generateData(simulation, 1)
		}
	}
}

/*func generateData(simulation Simulation, nbr_itr int) {

	h := make([][]string, simulation.Ris.N*simulation.Tx.N)
	for i := 0; i < len(h); i++ {
		h[i] = make([]string, nbr_itr)
	}
	g := make([][]string, simulation.Ris.N*simulation.Rx.N)
	for i := 0; i < len(g); i++ {
		g[i] = make([]string, nbr_itr)
	}
	//var wg sync.WaitGroup
	for i := 0; i < nbr_itr; i++ {

		H, G := simulation.Run()
		for hi, line := range H.Data {
			for hii, ele := range line {
				h[hi*len(line)+hii][i] = complextoString(ele)
			}
		}
		for gi, line := range G.Data {
			for gii, ele := range line {
				g[gi*len(line)+gii][i] = complextoString(ele)
			}
		}

		fmt.Println(i, " generation")

	}
	hcsv, err := os.Create("data/h_sim.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(hcsv)
	for _, row := range h {
		_ = csvwriter.Write(row)
	}
	csvwriter.Flush()
	defer hcsv.Close()

	gcsv, err := os.Create("data/g_sim.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter = csv.NewWriter(gcsv)
	for _, row := range g {
		_ = csvwriter.Write(row)
	}
	csvwriter.Flush()
	defer gcsv.Close()
}
*/
