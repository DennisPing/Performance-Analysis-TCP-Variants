package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/DennisPing/Performance-Analysis-TCP-Variants/common"
)

func main() {

	agents := []string{"Agent/TCP", "Agent/TCP/Reno", "Agent/TCP/Newreno", "Agent/TCP/Vegas"}

	pwd, _ := os.Getwd()
	basedir := filepath.Dir(pwd)

	// Check if the output directory exists
	if _, err := os.Stat(basedir + "/results"); os.IsNotExist(err) {
		os.Mkdir(basedir+"/results", 0777)
	}
	if _, err := os.Stat(basedir + "/results/exp01"); os.IsNotExist(err) {
		os.Mkdir(basedir+"/results/exp01", 0777)
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(agents))
	for _, agent := range agents {
		go ConcurrentSimulation(wg, agent)
	}
	wg.Wait()
	fmt.Println("Finished!")

}

func ConcurrentSimulation(wg *sync.WaitGroup, agent string) {
	defer wg.Done()

	// Split the agent string by '/' and get the last element
	split := strings.Split(agent, "/")
	suffix := split[len(split)-1]
	if suffix == "TCP" {
		suffix = "Tahoe"
	}

	pwd, _ := os.Getwd()
	basedir := filepath.Dir(pwd)
	filename := basedir + "/test/exp01/exp01_" + suffix + ".csv"
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString("cbr_rate,avg_throughput,std_throughput\n")
	file.Close()

	// A 3D slice to store cbr_rate, avg_throughput, std_throughput
	var results [][][]float64

	// The main simulation loop of 50 trails for each cbr_rate from 1 to 9 Mbps
	for rate := 1; rate < 10; rate++ {
		start := time.Now()
		fmt.Printf("Starting %s with rate %d\n", agent, rate)
		cumul_throughputs := make([]float64, 0)

		// Simulation variables
		fid := 1
		from_node := 1
		to_node := 2
		cbr_start := 0.0

		for tcp_start := 0.5; tcp_start <= 5.5; tcp_start += 0.1 {
			traces := Simulation01(agent, fid, from_node, to_node, tcp_start, cbr_start, float64(rate))
			window_size := 0.2
			_, _, throughput := common.CalculateThroughput(traces, from_node, to_node, tcp_start, window_size)
			cumul_throughputs = append(cumul_throughputs, throughput)
		}

		avg_throughput := common.Mean(cumul_throughputs)
		std_throughput := common.StdDev(cumul_throughputs)
		results = append(results, [][]float64{
			{float64(rate), avg_throughput, std_throughput},
		})
		end := time.Since(start).Round(time.Second)
		fmt.Printf("Finished %s with rate %d in %s\n", agent, rate, end)
	}

	// Write results to CSV file
	file2, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file2.Close()
	w := csv.NewWriter(file2)
	defer w.Flush()
	for _, result := range results {
		w.Write([]string{strconv.Itoa(int(result[0][0])), strconv.FormatFloat(result[0][1], 'f', -1, 64), strconv.FormatFloat(result[0][2], 'f', -1, 64)})
	}
}

// Run Simulation 1 using ns2 return a slice of traces
func Simulation01(agent string, fid int, from_node int, to_node int, tcp_start float64, cbr_start float64, cbr_rate float64) []*common.Trace {
	split := strings.Split(agent, "/")
	suffix := split[len(split)-1]
	filename := "outfile_" + suffix + ".tr"

	// Run a command from the shell
	cmd := exec.Command("ns", "../ns2/simulation01.tcl", agent, strconv.FormatFloat(tcp_start, 'f', -1, 64), strconv.FormatFloat(cbr_start, 'f', -1, 64), strconv.FormatFloat(cbr_rate, 'f', -1, 64), filename, "False")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	traces, err := common.ParseTraceFile(filename)
	if err != nil {
		panic(err)
	}
	os.Remove(filename)

	// Prepare the trace data
	traces = common.FilterByType(traces, "tcp")
	traces = common.FilterByFid(traces, fid)
	return traces
}
