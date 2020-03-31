package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type CPUTempObj struct {
	TimeStamp   time.Time
	HostAddress string
	CPUTemp     float64
}

func lambdaStateDiscovery(v CPUTempObj) (string, float64, string, string) {
	cpu_temp := v.CPUTemp
	cpu_temp_state := "CPU_TEMP_NONDETERMINISTIC"
	host_address := v.HostAddress
	timestamp := v.TimeStamp.Format(time.StampNano)

	if cpu_temp <= 3 || cpu_temp >= 98 {
		cpu_temp_state = "CPU_TEMP_CRITICAL"
	} else if cpu_temp >= 93 && cpu_temp < 98 {
		cpu_temp_state = "CPU_TEMP_HIGH"
	} else if cpu_temp > 3 && cpu_temp < 93 {
		cpu_temp_state = "CPU_TEMP_OK"
	}
	return timestamp, cpu_temp, cpu_temp_state, host_address

}

func collectCPUTemperature(hostName string) {

	for {
		resp, err := http.Get("http://" + hostName + "/redfish/v1/Chassis/1/Thermal")
		if err != nil {
			return
		}
		defer resp.Body.Close()

		var result CPUTempObj
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal(body, &result)
		timestamp, cpu_temp, cpu_temp_state, host_address := lambdaStateDiscovery(result)
		fmt.Printf("%v %s %.2fC %s\n", timestamp, host_address, cpu_temp, cpu_temp_state)
		time.Sleep(1 * time.Second)
	}
}

func main() {

	// Poll 200 servers
	var nodeList [200]string

	//Fill array with server hostnames
	for i := range nodeList {
		var nodeNum = strconv.Itoa(i)
		nodeList[i] = "server" + nodeNum + ":8000"
	}

	var wg sync.WaitGroup

	for _, node := range nodeList {
		wg.Add(1)
		go func(nodeAddress string) {
			defer wg.Done()
			collectCPUTemperature(nodeAddress)
		}(node)
	}

	wg.Wait()
}
