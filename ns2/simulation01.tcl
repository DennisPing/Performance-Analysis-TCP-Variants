# Create a simulator object
set ns [new Simulator]

# **************************************************************************
# Parse CLI arguments

if { $argc != 6 } {
    puts "You only provided $argc arguments. Simulation requires 6 arguments."
    puts "ns simulation01.tcl <tcpAgent> <tcpStart> <cbrStart> <cbrRate> <outputFile> <verbose>"
    exit 1
}

# If argument parsing is broken, use this for loop for debugging.
# set k 1
# foreach  i  $argv  {
#     puts "arg $k is $i"
#     incr k
# }

set tcpAgent [lindex $argv 0]
set tcpStart [lindex $argv 1]
set cbrStart [lindex $argv 2]
set cbrRate [lindex $argv 3]
set outFile [lindex $argv 4]
set verbose [lindex $argv 5]

if { $verbose == True } {
    puts "Simulation started with parameters:"
    puts "tcpAgent = $tcpAgent"
    puts "tcpStart = $tcpStart sec"
    puts "cbrStart = $cbrStart sec"
    puts "cbrRate = $cbrRate mbps"
    puts "outFile = $outFile"
}

set validAgents [list "Agent/TCP" "Agent/TCP/Reno" "Agent/TCP/Newreno" "Agent/TCP/Vegas"]
set tcpAgentIndex [lsearch $validAgents $tcpAgent]
if { $tcpAgentIndex < 0 } {
    puts "Error: Invalid TCP agent '$tcpAgent'. Valid agents are: "
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

# Setup a TCP Tahoe connection between n1 and n4
set tcp [new $tcpAgent]
$tcp set window_ 100
set sink [new Agent/TCPSink]
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
$ns attach-agent $n2 $udp
$ns attach-agent $n3 $null
$ns connect $udp $null
$udp set fid_ 2

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
$ns at $cbrStart "$cbr start"
$ns at $tcpStart "$ftp start"
# stop time is at $tcpStart + 30 sec
$ns at [expr {$tcpStart+30.0}] "$cbr stop"
$ns at [expr {$tcpStart+30.0}] "$ftp stop"

# Detach tcp and sink agents
$ns at [expr {$tcpStart+29.5}] "$ns detach-agent $n1 $tcp ; $ns detach-agent $n4 $sink"

# Call the finish procedure when the simulation is done
$ns at [expr {$tcpStart+30.0}] "finish"

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
