# RIS SIMULTOR

## 1. Installation

### 1.1 Simulator

In order to run the simulator:
- [golang (at least version 1.18)](https://go.dev/dl/)

clone the simulator repo:
```bash
    git clone https://gitlab.eurecom.fr/ris-esi/ris-simulator.git

    cd ris-simultor
    # download dependencies 
    go mod tidy
    # build the binary
    go build .
    # run the test
    TO BE ADDED
```
### 1.2 Flexric

clone the flexric repo and download the dependecies it requires:
- [flexric](https://gitlab.eurecom.fr/ris-esi/flexric.git)

few compilation dependencies which are not mentioned in the ReadMe so execute the following:
>***Note***: the following commands may break some other dependencies in your system, so better run this in a VM
```sh
    sudo apt -y install gcc-9 g++-9 gcc-10 g++-10
    sudo apt update-alternatives --install /usr/bin/gcc gcc /usr/bin/gcc-9 9
    sudo apt update-alternatives --install /usr/bin/g++ g++ /usr/bin/g++-9 9
    sudo apt update-alternatives --install /usr/bin/gcc gcc /usr/bin/gcc-10 10
    sudo apt update-alternatives --install /usr/bin/g++ g++ /usr/bin/g++-10 10
    
    # to make sure that the commands ran successfully
    sudo update-alternatives --config gcc
    sudo update-alternatives --config gcc
```
after that follow the build instructions in flexric repo.

## 2.Simulation

To run the simulation we can use the default configuration in the **"init.json"** or we can make our own by making changes to the file (don't forget to backup the default config)

>**Note**: for this version of the simulator we can support up to 64 RIS elements (until we can fix the what seems to be a buffer overflow in the flexric)

```json
{
    //The frequency in which gnb shall run
	"Frequency":28.0,
    //Indoor environment of the simulation should abid to 3gpp constraints
    //units in meters
	"Environment": 
	{
		"length":75.0,
		"width":50.0,
		"height":3.5
	},
	"Equipements": 
	{
        // gNB
		"TX": 
		{
        //Number of Antennas
		 "Elements":1 ,
		 "Coordinates": 
		 {
			"x":0.0,
			"y":25.0,
			"z":2.0
		 },
         //Either 0:linear or 1:planer
		 "Type": 0
		},
		"RX": 
		{
		 "Elements":1 ,
		 "Coordinates": 
		 {
			"x":38.0,
			"y":48.0,
			"z":1.0
		 },
		 "Type": 0
		},
		"RIS": 
		{
         // number of RIS patches of now should be under 64
		 "Elements":16 ,
		 "Coordinates": 
		 {
			"x":40.0,
			"y":50.0,
			"z":2.0
		 }
		}
	}	
}

```

Now that everything is set-up we can run our simulation:
``` bash 
    # RUN Simulator
    ./RIS_SIMULATOR
```
Run E2 Agent:
```bash
    /path/to/flexric/build/examples/emulator/agent/emu_agent_gnb
```
Run RIC:
```bash
    /path/to/flexric/build/examples/nearRT-RIC
```
Run Xapp:
```bash
    python3 /path/to/flexric/build/examples/xApp/python3 xapp_ris_moni.py
```
