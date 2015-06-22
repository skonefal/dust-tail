package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"
)

const (
	SAMPLING_TIME             = 1 * time.Second        // interval between sampling
	HTTP_TIMEOUT              = 200 * time.Millisecond // endpoint timeout
	EXPERIMENT_TIME           = 5 * time.Second        //
	EXPERIMENT_RESULTS_FOLDER = "results"              // folder with results

	STATISTICS_ENDPOINT = "/monitor/serenity_statistics.json" // mesos worker statistics endpoint
)

var mesosAgents = [...]string{
	"http://localhost:5051"} //list of mesos workers that will be sampled

var (
	/// regexp for returning node name
	nodeRegexp          = regexp.MustCompile(`^(http:\/\/)?([^:]+):?(\d*)\/?`)
	experimentStartTime = time.Now()
	//	experimentStartTime, ee = time.Parse(time.RFC3339, time.Now().String())
)

type UsageStats struct {
	usage    string
	endpoint string
}

func getUsage(address string, responsec chan<- *UsageStats) {

	client := http.Client{
		Timeout: HTTP_TIMEOUT,
	}
	endpoint := address + STATISTICS_ENDPOINT
	response, err := client.Get(endpoint)
	if err != nil {
		fmt.Printf("Error while obtaining statistics: %s", err)
		return
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Error while reading statistics: %s", err)
			return
		}
		usageStats := UsageStats{
			usage:    string(contents),
			endpoint: endpoint,
		}

		responsec <- &usageStats
	}
}

func harvestUsage(usagec chan<- *UsageStats) {
	for _, address := range mesosAgents {
		go getUsage(address, usagec)
	}
}

func saveUsage(usage *UsageStats) {

	resultsFile, err := createResulsFilename(usage.endpoint)
	if err != nil {
		fmt.Printf("Error while making up file name for %s | %s", usage.endpoint, err)
		return
	}

	f, err := os.OpenFile(resultsFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString(usage.usage)
}

// createProperJsonArrays
// saveUsage function writes multiple arrays into single file
// this function creates a valid JSON array of arrays of resource usage
func createProperJsonArrays() {
	for _, endpoint := range mesosAgents {
		resultsFile, err := createResulsFilename(endpoint)
		if err != nil {
			fmt.Printf("Error while making up file name for %s | %s", endpoint, err)
			break
		}
		createJsonArrayOfArrays(resultsFile)
	}
}

func main() {

	err := createResultsPath(EXPERIMENT_RESULTS_FOLDER)
	if err != nil {
		return
	}

	usagec := make(chan *UsageStats, len(mesosAgents))
	tickSignal := time.After(SAMPLING_TIME)
	experimentTimeout := time.After(EXPERIMENT_TIME)
	for {
		select {
		case _ = <-tickSignal:
			harvestUsage(usagec)
			tickSignal = time.After(SAMPLING_TIME)
		case usage := <-usagec:
			saveUsage(usage)
		case _ = <-experimentTimeout:
			createProperJsonArrays()
			return
		}
	}
}
