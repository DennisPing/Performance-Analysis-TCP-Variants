# Create a simulator object
set ns [new Simulator]

# **************************************************************************
# Parse CLI arguments

if { $argc != 5} {
    puts "You only provided $argc arguments. Simulation requires 5 arguments."
    puts "ns simulation01.tcl <tcpAgent> <queueType> <cbrStart> <outputFile> <verbose>"
    exit 1
}

# If argument parsing is broken, use this for loop for debugging.
# set k 1
# foreach  i  $argv  {
#     puts "arg $k is $i"
#     incr k
# }

set tcpAgent [lindex $argv 0]
set queueType [lindex $argv 1]
set cbrStart [lindex $argv 2]
set outFile [lindex $argv 3]
set verbose [lindex $argv 4]

if { $verbose == True } {
    puts "Simulation started with parameters:"
    puts "tcpAgent = $tcpAgent"
    puts "queueType = $queueType"
    puts "tcpStart = 0 sec"
    puts "cbrStart = $cbrStart sec"
    puts "cbrRate = 8 mbps"
    puts "outFile = $outFile"
}

set validAgents [list "Agent/TCP/Reno" "Agent/TCP/Sack1"]
set tcpAgentIndex [lsearch $validAgents $tcpAgent]
if { $tcpAgentIndex < 0 } {
    puts "Error: Invalid TCP agent '$tcpAgent'. Valid agents are: "
    foreach agent $validAgents {
        puts "\t$agent"
    }
    exit 1
}

if { $tcpAgentIndex < 1 } {
    set sink [new Agent/TCPSink]
} else {
    set sink [new Agent/TCPSink/Sack1]
}

set validQueues [list "DropTail" "RED"]
set queueIndex [lsearch $validQueues $queueType]
if { $queueIndex < 0 } {
    puts "Error: Invalid queue type '$queueType'. Valid types are: "
    foreach agent $validQueues {
        puts "\t$agent"
    }
    exit 1
}

# **************************************************************************
# Run the simulation

# Open the trace file
set traceFile [open $outFile w]
$ns trace-all $traceFile

# Define a "finish" procedure
proc finish {} {
    global ns traceFile
    $ns flush-trace
    close $traceFile
    exit 0
}

# Create six nodes
set n1 [$ns node]
set n2 [$ns node]
set n3 [$ns node]
set n4 [$ns node]
set n5 [$ns node]
set n6 [$ns node]

# Create links between the nodes
# Bandwidth for all links is 10 Mbps, round trip time is 10 ms
$ns duplex-link $n1 $n2 10Mb 10ms $queueType
$ns duplex-link $n2 $n3 10Mb 10ms $queueType
$ns duplex-link $n3 $n4 10Mb 10ms $queueType
$ns duplex-link $n3 $n6 10Mb 10ms $queueType
$ns duplex-link $n2 $n5 10Mb 10ms $queueType

# Set the queue limit
$ns queue-limit $n2 $n3 50

# Setup a TCP connection between n1 and n4
set tcp [new $tcpAgent]
$tcp set window_ 100
# We already set up sink
$ns attach-agent $n1 $tcp
$ns attach-agent $n4 $sink
$ns connect $tcp $sink
$tcp set fid_ 1

# Setup an FTP over TCP connection
set ftp [new Application/FTP]
$ftp attach-agent $tcp
$ftp set type_ FTP

# Setup a UDP connection at n2 to n3
set udp [new Agent/UDP]
set null [new Agent/Null]
$ns attach-agent $n1 $udp
$ns attach-agent $n4 $null
$ns connect $udp $null
$udp set fid_ 2

# Setup a CBR over the UDP connection
# Documentation: https://www.isi.edu/nsnam/ns/doc/node510.html 
set cbr [new Application/Traffic/CBR]
$cbr attach-agent $udp
$cbr set type_ CBR
$cbr set rate_ 8Mb

# Schedule events for the CBR agent
$ns at $cbrStart "$cbr start"
$ns at 0 "$ftp start"
$ns at 30.0 "$cbr stop"
$ns at 30.0 "$ftp stop"

# Detach tcp and sink agents
$ns at 29.5 "$ns detach-agent $n1 $tcp ; $ns detach-agent $n4 $sink"

# Call the finish procedure when the simulation is done
$ns at 30.0 "finish"

if { $verbose == True } {
    puts " "
    # Print TCP packet size
    puts "TCP packet size = [$tcp set packetSize_] bytes"

    # Print CBR packet size
    puts "CBR packet size = [$cbr set packetSize_] bytes"
    # Print CBR interval
    puts "CBR interval = [$cbr set interval_] sec"
    # Print CBR rate in megabits/sec
    set mbps [expr [$cbr set rate_] / 1000000]
    puts "CBR rate = $mbps mbps"
    puts " "
}

# Run the simulation
$ns run
