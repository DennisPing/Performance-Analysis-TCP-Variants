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

	RenoReno := []string{agents[1], agents[1]}
	NewrenoReno := []string{agents[2], agents[1]}
	VegasVegas := []string{agents[3], agents[3]}
	NewrenoVegas := []string{agents[2], agents[3]}
	combos := [][]string{RenoReno, NewrenoReno, VegasVegas, NewrenoVegas}

	pwd, _ := os.Getwd()
	basedir := filepath.Dir(pwd)

	// Check if the output directory exists
	if _, err := os.Stat(basedir + "/results"); os.IsNotExist(err) {
		os.Mkdir(basedir+"/results", 0777)
	}
	if _, err := os.Stat(basedir + "/results/exp02"); os.IsNotExist(err) {
		os.Mkdir(basedir+"/results/exp02", 0777)
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(combos))
	for _, combo := range combos {
		go Experiment02(wg, combo)
	}
	wg.Wait()
	fmt.Println("Finished!")

}

func Experiment02(wg *sync.WaitGroup, combo []string) {
	defer wg.Done()

	agent1 := combo[0]
	agent2 := combo[1]

	// Split the agent string by '/' and get the last element
	split1 := strings.Split(agent1, "/")
	suffix1 := split1[len(split1)-1]
	if suffix1 == "TCP" {
		suffix1 = "Tahoe"
	}

	split2 := strings.Split(agent2, "/")
	suffix2 := split2[len(split2)-1]
	if suffix2 == "TCP" {
		suffix2 = "Tahoe"
	}

	pwd, _ := os.Getwd()
	basedir := filepath.Dir(pwd)
	filename := basedir + "/results/exp02/exp02_" + suffix1 + "_" + suffix2 + ".csv"
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString("cbr_rate,avg_throughput1,std_throughput1,avg_latency1,std_latency1,avg_drops1,std_drops1,avg_throughput2,std_throughput2,avg_latency2,std_latency2,avg_drops2,std_drops2\n")
	file.Close()

	var results [][]float64

	// The main simulation loop of 50 trails for each cbr_rate from 1 to 9 Mbps
	for rate := 1; rate < 10; rate++ {
		start := time.Now()
		fmt.Printf("Starting %s/%s with rate %d\n", agent1, agent2, rate)
		cumul_throughputs1 := make([]float64, 0)
		cumul_latencies1 := make([]float64, 0)
		cumul_drops1 := make([]float64, 0)

		cumul_throughputs2 := make([]float64, 0)
		cumul_latencies2 := make([]float64, 0)
		cumul_drops2 := make([]float64, 0)

		// Simulation variables
		fid := 1
		from_node := 1 // ns2 counts from 0, so this is node #2 in the diagram
		to_node := 2

		for tcp2_start := 0.0; tcp2_start <= 5.0; tcp2_start += 0.05 {
			traces := Simulation02(agent1, agent2, fid, from_node, to_node, tcp2_start, float64(rate))

			// Prepare the trace data
			traces = pkg.FilterByType(traces, "tcp")
			traces1 := pkg.FilterByFid(traces, 1)
			traces2 := pkg.FilterByFid(traces, 2)

			// Calculate throughput, latency, and dropped packets
			window_size := 0.2
			_, _, throughput1 := pkg.CalculateThroughput(traces1, from_node, to_node, tcp2_start, window_size)
			_, _, latency1 := pkg.CalculateLatency(traces1, from_node, to_node, tcp2_start)
			drops1 := pkg.CountDrops(traces1)

			_, _, throughput2 := pkg.CalculateThroughput(traces2, from_node, to_node, tcp2_start, window_size)
			_, _, latency2 := pkg.CalculateLatency(traces1, from_node, to_node, tcp2_start)
			drops2 := pkg.CountDrops(traces1)

			// Add the results to the cumulative results
			cumul_throughputs1 = append(cumul_throughputs1, throughput1)
			cumul_latencies1 = append(cumul_latencies1, latency1)
			cumul_drops1 = append(cumul_drops1, float64(drops1))

			cumul_throughputs2 = append(cumul_throughputs2, throughput2)
			cumul_latencies2 = append(cumul_latencies2, latency2)
			cumul_drops2 = append(cumul_drops2, float64(drops2))
		}

		avg_throughput1 := pkg.Mean(cumul_throughputs1)
		avg_latency1 := pkg.Mean(cumul_latencies1)
		avg_drops1 := pkg.Mean(cumul_drops1)
		std_throughput1 := pkg.StdDev(cumul_throughputs1)
		std_latency1 := pkg.StdDev(cumul_latencies1)
		std_drops1 := pkg.StdDev(cumul_drops1)

		avg_throughput2 := pkg.Mean(cumul_throughputs2)
		avg_latency2 := pkg.Mean(cumul_latencies2)
		avg_drops2 := pkg.Mean(cumul_drops2)
		std_throughput2 := pkg.StdDev(cumul_throughputs2)
		std_latency2 := pkg.StdDev(cumul_latencies2)
		std_drops2 := pkg.StdDev(cumul_drops2)

		results = append(results,
			[]float64{float64(rate), avg_throughput1, std_throughput1, avg_latency1, std_latency1, avg_drops1, std_drops1,
				avg_throughput2, std_throughput2, avg_latency2, std_latency2, avg_drops2, std_drops2})

		end := time.Since(start).Round(time.Second)
		fmt.Printf("Finished %s/%s with rate %d in %s\n", agent1, agent2, rate, end)
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
func Simulation02(agent1 string, agent2 string, fid int, from_node int, to_node int, tcp2_start float64, cbr_rate float64) []*pkg.Trace {
	split1 := strings.Split(agent1, "/")
	suffix1 := split1[len(split1)-1]
	split2 := strings.Split(agent2, "/")
	suffix2 := split2[len(split2)-1]
	filename := "outfile_" + suffix1 + "_" + suffix2 + ".tr"

	// Run a command from the shell
	tcp2Start := strconv.FormatFloat(tcp2_start, 'f', -1, 64)
	cmd := exec.Command("ns", "../ns2/simulation02.tcl", agent1, agent2, tcp2Start, strconv.FormatFloat(cbr_start, 'f', -1, 64), strconv.FormatFloat(cbr_rate, 'f', -1, 64), filename, "False")
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
