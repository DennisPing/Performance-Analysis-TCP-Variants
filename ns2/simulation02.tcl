# Create a simulator object
set ns [new Simulator]

# **************************************************************************
# Parse CLI arguments

if { $argc != 6 } {
    puts "You only provided $argc arguments. Simulation requires 6 arguments."
    puts "ns simulation01.tcl <tcpAgent1> <tcpAgent2> <tcpStart2> <cbrRate> <outputFile> <verbose>"
    exit 1
}

# If argument parsing is broken, use this for loop for debugging.
# set k 1
# foreach  i  $argv  {
#     puts "arg $k is $i"
#     incr k
# }

set tcpAgent1 [lindex $argv 0]
set tcpAgent2 [lindex $argv 1]
set tcpStart2 [lindex $argv 2]
set cbrRate [lindex $argv 3]
set outFile [lindex $argv 4]
set verbose [lindex $argv 5]

if { $verbose == True } {
    puts "Simulation started with parameters:"
    puts "tcpAgent1 = $tcpAgent1"
    puts "tcpAgent2 = $tcpAgent2"
    puts "tcpStart1 = 4 sec"
    puts "tcpStart2 = $tcpStart2 sec"
    puts "cbrStart = 0 sec"
    puts "cbrRate = $cbrRate mbps"
    puts "outFile = $outFile"
}

set validAgents [list "Agent/TCP" "Agent/TCP/Reno" "Agent/TCP/Newreno" "Agent/TCP/Vegas"]
set tcpAgentIndex1 [lsearch $validAgents $tcpAgent1]
if { $tcpAgentIndex1 < 0 } {
    puts "Error: Invalid TCP agent '$tcpAgent1'. Valid agents are: "
    foreach agent $validAgents {
        puts "\t$agent"
    }
    exit 1
}
set tcpAgentIndex2 [lsearch $validAgents $tcpAgent2]
if { $tcpAgentIndex2 < 0 } {
    puts "Error: Invalid TCP agent '$tcpAgent2'. Valid agents are: "
    foreach agent $validAgents {
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
$ns duplex-link $n1 $n2 10Mb 10ms DropTail
$ns duplex-link $n2 $n3 10Mb 10ms DropTail
$ns duplex-link $n3 $n4 10Mb 10ms DropTail
$ns duplex-link $n3 $n6 10Mb 10ms DropTail
$ns duplex-link $n2 $n5 10Mb 10ms DropTail

# Set the queue limit
$ns queue-limit $n2 $n3 50

# Setup a TCP connection between n1 and n4
set tcp1 [new $tcpAgent1]
$tcp1 set window_ 100
set sink1 [new Agent/TCPSink]
$ns attach-agent $n1 $tcp1
$ns attach-agent $n4 $sink1
$ns connect $tcp1 $sink1
$tcp1 set fid_ 1

# Setup an FTP over TCP connection
set ftp1 [new Application/FTP]
$ftp1 attach-agent $tcp1
$ftp1 set type_ FTP

# Setup a TCPconnection between n5 and n6
set tcp2 [new $tcpAgent2]
$tcp2 set window_ 100
set sink2 [new Agent/TCPSink]
$ns attach-agent $n5 $tcp2
$ns attach-agent $n6 $sink2
$ns connect $tcp2 $sink2
$tcp2 set fid_ 2

# Setup an FTP over TCP connection
set ftp2 [new Application/FTP]
$ftp2 attach-agent $tcp2
$ftp2 set type_ FTP

# Setup a UDP connection at n2 to n3
set udp [new Agent/UDP]
set null [new Agent/Null]
$ns attach-agent $n2 $udp
$ns attach-agent $n3 $null
$ns connect $udp $null
$udp set fid_ 3

# Setup a CBR over the UDP connection
# Documentation: https://www.isi.edu/nsnam/ns/doc/node510.html 
set cbr [new Application/Traffic/CBR]
$cbr attach-agent $udp
$cbr set type_ CBR
# $cbr set interval_ 0.01
# $cbr set packetSize_ 1000
$cbr set rate_ ${cbrRate}Mb
# Note: rate_ and interval_ are mutually exclusive. Just choose one.
# Note: Either use just (rate_) or (interval_ and packetSize_)
# Note: TCL does not appear to be case sensitive. 10Mb and 10mb are the same.
# Warning: TCL is space sensitive. 10Mb is correct while 10 Mb is not.

# Schedule events for the CBR agent
$ns at 0 "$cbr start"
$ns at 4 "$ftp1 start"
$ns at $tcpStart2 "$ftp2 start"
$ns at 30.0 "$cbr stop"
$ns at 30.0 "$ftp1 stop"
$ns at 30.0 "$ftp2 stop"

# Detach tcp and sink agents
$ns at 29.5 "$ns detach-agent $n1 $tcp1 ; $ns detach-agent $n4 $sink1; $ns detach-agent $n5 $tcp2 ; $ns detach-agent $n6 $sink2"

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
