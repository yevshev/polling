package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	//"sync"
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

func GetNodeCPUTemp(nodeIP string, timeOut int) (CPUTempObj, error) {

	var result CPUTempObj
	url := "http://" + nodeIP + ":8000/redfish/v1/Chassis/1/Thermal"

	client := http.Client{
		Timeout: (time.Duration(timeOut) * time.Millisecond),
	}

	resp, err := client.Get(url)

	if err != nil {
		fmt.Printf("\nThe HTTP request failed with error %s\n", err)
		return result, err
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return result, err
		}

		json.Unmarshal(body, &result)
		return result, nil
	}

}

func main() {

	// Poll 50 servers
	var nodeList [50]string
	for i := range nodeList {
		nodeNum := strconv.Itoa(i)
		nodeList[i] = "server" + nodeNum
	}

	var nodeCPUTempList []CPUTempObj
	var errorList []string

	respc, errc := make(chan CPUTempObj), make(chan error)

	//ticker := time.NewTicker(60 * time.Second)

	// main loop
	// for {
	// 	select {
	// 	case <-ticker.C:
	for _, node := range nodeList {

		go func(nodeAddress string, timeout int) {

			//println(nodeAddress)
			resp, err := GetNodeCPUTemp(nodeAddress, timeout)
			fmt.Println(resp)
			if err != nil {
				errc <- err
				return
			}
			respc <- resp
		}(node, 500)

	}
	for i := 0; i < len(nodeList); i++ {
		select {
		case res := <-respc:
			nodeCPUTempList = append(nodeCPUTempList, res)
			timestamp, cpu_temp, cpu_temp_state, host_address := lambdaStateDiscovery(res)
			fmt.Printf("%v %s %.2fC %s\n", timestamp, host_address, cpu_temp, cpu_temp_state)
		case e := <-errc:
			errorList = append(errorList, e.Error())

		}
	}
	fmt.Printf("\n Total Successful Responses From all Nodes: %d\n", len(nodeCPUTempList))
	fmt.Printf("\n Total Errors: %d\n", len(errorList))

	// 	}
	// }
}
