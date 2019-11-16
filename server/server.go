package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type CPUTempObj struct {
	TimeStamp   time.Time
	HostAddress string
	CPUTemp     float64
}

func ResponseServer(w http.ResponseWriter, r *http.Request) {
	w.Write(GetCPUTemp())
}

func randTemperature(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return math.Floor((min+rand.Float64()*(max-min))*100) / 100
}

// CPU temperature
func GetCPUTemp() []byte {

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	//Its a mockup CPU temperature
	cpuTempObj := new(CPUTempObj)
	cpuTempObj.TimeStamp = time.Now()
	cpuTempObj.HostAddress = hostname
	cpuTempObj.CPUTemp = randTemperature(3.0, 98.0)

	jsonObj, err := json.Marshal(cpuTempObj)
	if err != nil {
		log.Println(fmt.Sprintf("Could not marshal the response data: %v", err))
	}
	return jsonObj

}

func main() {

	// Register endpoint.
	http.HandleFunc("/redfish/v1/", ResponseServer)

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	http.ListenAndServe(hostname+":8000", nil)
}
