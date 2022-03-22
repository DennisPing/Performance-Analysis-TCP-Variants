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

	"github.com/DennisPing/Performance-Analysis-TCP-Variants/pkg"
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
		go Experiment01(wg, agent)
	}
	wg.Wait()
	fmt.Println("Finished!")

}

func Experiment01(wg *sync.WaitGroup, agent string) {
	defer wg.Done()

	// Split the agent string by '/' and get the last element
	split := strings.Split(agent, "/")
	suffix := split[len(split)-1]
	if suffix == "TCP" {
		suffix = "Tahoe"
	}

	pwd, _ := os.Getwd()
	basedir := filepath.Dir(pwd)
	filename := basedir + "/results/exp01/exp01_" + suffix + ".csv"
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString("cbr_rate,avg_throughput,std_throughput,avg_latency,std_latency,avg_drops,std_drops\n")
	file.Close()

	var results [][]float64

	// The main simulation loop for cbr_rate of from 1 to 9 Mbps
	for rate := 1; rate < 10; rate++ {
		start := time.Now()
		fmt.Printf("Starting %s with rate %d\n", suffix, rate)
		cumul_throughputs := make([]float64, 0)
		cumul_latencies := make([]float64, 0)
		cumul_drops := make([]float64, 0)

		// Simulation variables
		fid := 1
		from_node := 1 // ns2 counts from 0, so this is N2 -> N3
		to_node := 2
		cbr_start := 0.0

		for tcp_start := 0.5; tcp_start <= 5.5; tcp_start += 0.1 {
			traces := Simulation01(agent, tcp_start, cbr_start, float64(rate))
			// Prepare the trace data
			traces = pkg.FilterByType(traces, "tcp")
			traces = pkg.FilterByFid(traces, fid)

			// Calculate throughput, latency, and dropped packets
			window_size := 0.2
			_, _, throughput := pkg.CalculateThroughput(traces, from_node, to_node, tcp_start, window_size)
			_, _, latency := pkg.CalculateLatency(traces, from_node, to_node, tcp_start)
			drops := pkg.CountDrops(traces)

			cumul_throughputs = append(cumul_throughputs, throughput)
			cumul_latencies = append(cumul_latencies, latency)
			cumul_drops = append(cumul_drops, float64(drops))
		}

		avg_throughput := pkg.Mean(cumul_throughputs)
		avg_latency := pkg.Mean(cumul_latencies)
		avg_drops := pkg.Mean(cumul_drops)
		std_throughput := pkg.StdDev(cumul_throughputs)
		std_latency := pkg.StdDev(cumul_latencies)
		std_drops := pkg.StdDev(cumul_drops)

		results = append(results, []float64{float64(rate), avg_throughput, std_throughput, avg_latency,
			std_latency, avg_drops, std_drops})

		end := time.Since(start).Round(time.Second)
		fmt.Printf("Finished %s with rate %d in %s\n", suffix, rate, end)
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
		line := make([]string, len(result))
		line[0] = strconv.Itoa(int(result[0])) // cbr_rate is an int
		for i := 1; i < len(result); i++ {
			line[i] = strconv.FormatFloat(result[i], 'f', 10, 64) // everything else is a float
		}
		w.Write(line)
	}
}

// Run Simulation 1 using ns2 and return a slice of traces
func Simulation01(agent string, tcp_start float64, cbr_start float64, cbr_rate float64) []*pkg.Trace {

	split := strings.Split(agent, "/")
	suffix := split[len(split)-1]
	filename := "outfile_" + suffix + ".tr"

	// Run a command from the shell
	cmd := exec.Command("ns", "../ns2/simulation01.tcl", agent, strconv.FormatFloat(tcp_start, 'f', -1, 64),
		strconv.FormatFloat(cbr_start, 'f', -1, 64), strconv.FormatFloat(cbr_rate, 'f', -1, 64), filename, "False")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	traces, err := pkg.ParseTraceFile(filename)
	if err != nil {
		panic(err)
	}
	os.Remove(filename)
	return traces
}
