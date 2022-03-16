package main

import (
	"encoding/csv"
	"fmt"
	"math"
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

	RenoDropTail := []string{"Agent/TCP/Reno", "DropTail"}
	RenoRED := []string{"Agent/TCP/Reno", "RED"}
	Sack1DropTail := []string{"Agent/TCP/Sack1", "DropTail"}
	Sack1RED := []string{"Agent/TCP/Sack1", "RED"}
	combos := [][]string{RenoDropTail, RenoRED, Sack1DropTail, Sack1RED}

	pwd, _ := os.Getwd()
	basedir := filepath.Dir(pwd)

	// Check if the output directory exists
	if _, err := os.Stat(basedir + "/results"); os.IsNotExist(err) {
		os.Mkdir(basedir+"/results", 0777)
	}
	if _, err := os.Stat(basedir + "/results/exp03"); os.IsNotExist(err) {
		os.Mkdir(basedir+"/results/exp03", 0777)
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(combos))
	for _, combo := range combos {
		go Experiment03(wg, combo)
	}
	wg.Wait()
	fmt.Println("Finished!")

}

func Experiment03(wg *sync.WaitGroup, combo []string) {
	defer wg.Done()

	agent := combo[0]
	queue := combo[1]

	// Split the agent string by '/' and get the last element
	split := strings.Split(agent, "/")
	suffix := split[len(split)-1]
	if suffix == "TCP" {
		suffix = "Tahoe"
	}

	pwd, _ := os.Getwd()
	basedir := filepath.Dir(pwd)
	filename := basedir + "/results/exp03/exp03_" + suffix + "_" + queue + ".csv"
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	header := "avg_throughput1,std_throughput1,avg_latency1,std_latency1,avg_drops1,std_drops1," +
		"avg_throughput2,std_throughput2,avg_latency2,std_latency2,avg_drops2,std_drops2\n"
	file.WriteString(header)
	file.Close()

	var results [][]float64

	start := time.Now()
	fmt.Printf("Starting %s with queue %s\n", suffix, queue)
	cumul_throughputs1 := make([]float64, 0)
	cumul_latencies1 := make([]float64, 0)
	cumul_drops1 := make([]float64, 0)

	cumul_throughputs2 := make([]float64, 0)
	cumul_latencies2 := make([]float64, 0)
	cumul_drops2 := make([]float64, 0)

	// sample_times := make([]float64, 0)
	// sample_throughputs := make([]float64, 0)

	// Simulation variables
	fid := 1
	from_node := 1 // ns2 counts from 0, so this is N2 -> N3
	to_node := 2

	// TCP starts at t=0, let it stabilize, then start CBR at t=5
	for cbr_start := 5.0; cbr_start <= 10.0; cbr_start += 0.1 {
		traces := Simulation03(agent, queue, fid, from_node, to_node, cbr_start)

		// Prepare the trace data
		tcpTraces := pkg.FilterByType(traces, "tcp")
		cbrTraces := pkg.FilterByType(traces, "cbr")
		tcpTraces = pkg.FilterByFid(tcpTraces, 1)
		cbrTraces = pkg.FilterByFid(cbrTraces, 2)

		// Calculate throughput, latency, and dropped packets
		window_size := 0.2
		time_ticks1, throughput_ticks1, throughput1 := pkg.CalculateThroughput(tcpTraces, from_node, to_node, 0.0, window_size)
		_, _, latency1 := pkg.CalculateLatency(tcpTraces, from_node, to_node, 0.0)
		drops1 := pkg.CountDrops(tcpTraces)

		time_ticks2, throughput_ticks2, throughput2 := pkg.CalculateThroughput(cbrTraces, from_node, to_node, cbr_start, window_size)
		_, _, latency2 := pkg.CalculateLatency(cbrTraces, from_node, to_node, cbr_start)
		drops2 := pkg.CountDrops(cbrTraces)

		// Add the results to the cumulative results
		cumul_throughputs1 = append(cumul_throughputs1, throughput1)
		cumul_latencies1 = append(cumul_latencies1, latency1)
		cumul_drops1 = append(cumul_drops1, float64(drops1))

		cumul_throughputs2 = append(cumul_throughputs2, throughput2)
		cumul_latencies2 = append(cumul_latencies2, latency2)
		cumul_drops2 = append(cumul_drops2, float64(drops2))

		// Record the time vs throughput for t=10 specifically
		if math.Abs(cbr_start-10.0) < 0.001 {
			// Write time_ticks1 (column 1) and thoughtput_ticks1 (column 2) to a csv file
			file, err := os.Create(basedir + "/results/exp03/exp03_" + suffix + "_" + queue + "_TCP.csv")
			if err != nil {
				panic(err)
			}
			defer file.Close()
			file.WriteString("time_ticks,throughput_ticks\n")
			for i, tick := range time_ticks1 {
				file.WriteString(fmt.Sprintf("%f,%f\n", tick, throughput_ticks1[i]))
			}
			file2, err2 := os.Create(basedir + "/results/exp03/exp03_" + suffix + "_" + queue + "_CBR.csv")
			if err2 != nil {
				panic(err2)
			}
			defer file2.Close()
			file2.WriteString("time_ticks,throughput_ticks\n")
			for i, tick := range time_ticks2 {
				file2.WriteString(fmt.Sprintf("%f,%f\n", tick, throughput_ticks2[i]))
			}
		}
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
		[]float64{avg_throughput1, std_throughput1, avg_latency1, std_latency1, avg_drops1, std_drops1,
			avg_throughput2, std_throughput2, avg_latency2, std_latency2, avg_drops2, std_drops2})

	end := time.Since(start).Round(time.Second)
	fmt.Printf("Finished %s with queue %s in %s\n", suffix, queue, end)

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
		for i := 0; i < len(result); i++ {
			line[i] = strconv.FormatFloat(result[i], 'f', 10, 64) // everything else is a float
		}
		w.Write(line)
	}
}

// Run Simulation 2 using ns2 and return a slice of traces. CBR always starts at t=0 here.
func Simulation03(agent string, queue string, fid int, from_node int, to_node int, cbr_start float64) []*pkg.Trace {
	split := strings.Split(agent, "/")
	suffix := split[len(split)-1]
	filename := "outfile_" + suffix + "_" + queue + ".tr"

	// Run a command from the shell
	cmd := exec.Command("ns", "../ns2/simulation03.tcl", agent, queue, strconv.FormatFloat(cbr_start, 'f', -1, 64),
		filename, "False")
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
