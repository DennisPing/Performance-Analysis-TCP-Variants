package pkg

import (
	"bufio"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Trace struct {
	event       string
	time        float64
	from        int
	to          int
	packet_type string
	packet_size int
	fid         int
	seq         int
	pid         int
}

// String function for Trace struct
func (t *Trace) String() string {
	str := t.event + " " +
		strconv.FormatFloat(t.time, 'f', -1, 64) + " " +
		strconv.Itoa(t.from) + " " +
		strconv.Itoa(t.to) + " " +
		t.packet_type + " " +
		strconv.Itoa(t.packet_size) + " " +
		"-------" + " " +
		strconv.Itoa(t.fid) + " " +
		strconv.Itoa(t.seq) + " " +
		strconv.Itoa(t.pid)
	return str
}

// Parse the trace file and return a slice of Trace structs
func ParseTraceFile(file string) ([]*Trace, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var traces []*Trace
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, " ")
		time, _ := strconv.ParseFloat(fields[1], 64)
		from, _ := strconv.Atoi(fields[2])
		to, _ := strconv.Atoi(fields[3])
		packet_size, _ := strconv.Atoi(fields[5])
		fid, _ := strconv.Atoi(fields[7])
		seq, _ := strconv.Atoi(fields[10])
		pid, _ := strconv.Atoi(fields[11])

		trace := &Trace{
			event:       fields[0],
			time:        time,
			from:        from,
			to:          to,
			packet_type: fields[4],
			packet_size: packet_size,
			fid:         fid,
			seq:         seq,
			pid:         pid,
		}
		traces = append(traces, trace)
	}
	return traces, err
}

// Get a slice of traces of flow id 'fid'
func FilterByFid(traces []*Trace, fid int) []*Trace {
	var filtered []*Trace
	for _, trace := range traces {
		if trace.fid == fid {
			filtered = append(filtered, trace)
		}
	}
	return filtered
}

// Get a slice of traces of type 'packet_type' (tcp, cbr, ack)
func FilterByType(traces []*Trace, packet_type string) []*Trace {
	var filtered []*Trace
	for _, trace := range traces {
		if trace.packet_type == packet_type {
			filtered = append(filtered, trace)
		}
	}
	return filtered
}

// Calculate throughput vs time given a TCP flow start time
// Return slice times, slice throughputs, and average throughput
func CalculateThroughput(traces []*Trace, from_node int, to_node int, flow_start float64, window_size float64) ([]float64, []float64, float64) {
	var time_ticks []float64
	var throughput_ticks []float64

	var recv_times []float64
	for _, trace := range traces {
		if trace.event == "r" && trace.from == from_node && trace.to == to_node {
			recv_times = append(recv_times, trace.time)
		}
	}

	sort.Float64s(recv_times)

	var head int           // The index of the window head
	var tail int           // The index of the window tail
	var win_throughput int // The number of bytes in the current window
	var tot_throughput int // The total running throughput

	for head < len(recv_times) {
		// If a new trace enters the window
		if recv_times[head] < recv_times[tail]+window_size {
			win_throughput += traces[head].packet_size
			tot_throughput += traces[head].packet_size
			time_ticks = append(time_ticks, recv_times[head])
			head++
			// If a trace leaves the window
		} else if recv_times[head] > recv_times[tail]+window_size {
			win_throughput -= traces[tail].packet_size
			time_ticks = append(time_ticks, recv_times[tail]+window_size)
			tail++
			// If a trace enters the window and leaves the window at the same time
		} else {
			win_throughput += traces[head].packet_size - traces[tail].packet_size
			tot_throughput += traces[head].packet_size
			time_ticks = append(time_ticks, recv_times[head])
			head++
			tail++
		}
		throughput_ticks = append(throughput_ticks, (float64(win_throughput))/window_size/125000) // In Mbps
	}
	avg_throughput := (float64(tot_throughput) / (Max(time_ticks) - Min(time_ticks) - window_size)) / 125000 // In Mbps
	return time_ticks, throughput_ticks, avg_throughput
}

// Calculate latency vs time given a TCP flow start time
// Return slice times, slice throughputs, and average latency
func CalculateLatency(traces []*Trace, from_node int, to_node int, flow_start float64) ([]float64, []float64, float64) {
	var time_ticks []float64
	var latency_ticks []float64

	// A hashmap with {key, value} of {pid, Trace}
	start_traces := make(map[int]*Trace) // Traces of event '+'
	end_traces := make(map[int]*Trace)   // Traces of event 'r'

	// Iterate backwards through the trace list
	// If you encounter a 'd', skip the next 2 traces
	// The next 2 traces will always be '+' and 'r' at the exact same time
	for i := len(traces) - 1; i >= 0; i-- {
		trace := traces[i]
		if trace.event == "d" {
			i -= 3
			continue
		} else if trace.event == "+" && trace.from == from_node && trace.to == to_node {
			start_traces[trace.pid] = trace
		} else if trace.event == "r" && trace.from == from_node && trace.to == to_node {
			end_traces[trace.pid] = trace
		}
	}

	// In some cases, the number of start_traces > end_traces. That's ok.
	// A failed hashmap lookup is just a no-op
	for _, end_trace := range end_traces {
		start_trace, ok := start_traces[end_trace.pid] // Find the matches
		if ok {
			latency := end_trace.time - start_trace.time
			time_ticks = append(time_ticks, end_trace.time)
			latency_ticks = append(latency_ticks, latency)
		}
	}
	return time_ticks, latency_ticks, Mean(latency_ticks)
}

// Count the number of dropped packets. The trace should already be filtered by fid
func CountDrops(traces []*Trace) int {
	var drops int
	for _, trace := range traces {
		if trace.event == "d" {
			drops++
		}
	}
	return drops
}
