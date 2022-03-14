# Performance-Analysis-TCP-Variants

## Project Setup

```txt
.
├── bin                 <-- Go binaries
│   ├── exp01
│   ├── exp02
│   └── exp03
├── cmd                 <-- Experiment 1, 2, 3 code
│   ├── exp01
│   │   └── main.go
│   ├── exp02
│   │   └── main.go
│   └── exp03
│       └── main.go
├── common              <-- Common packages
│   ├── stats.go
│   └── trace.go
├── go.mod
├── graph               <-- Graph results with Python
│   ├── graph_exp01.py
│   ├── graph_exp02.py
│   └── graph_exp03.py
├── LICENSE
├── Makefile
├── ns2                 <-- ns2 simulation scripts
│   ├── simulation01.tcl
│   ├── simulation02.tcl
│   └── simulation03.tcl
├── README.md
└── results             <-- Experiment results
    ├── exp01
    │   ├── exp01_Newreno.csv
    │   ├── exp01_Reno.csv
    │   ├── exp01_Tahoe.csv
    │   └── exp01_Vegas.csv
    ├── exp02
    └── exp03

```