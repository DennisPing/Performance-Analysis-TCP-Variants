# Performance-Analysis-TCP-Variants

Matthew Jones  
Dennis Ping

## Background

This project was initially done in Python and was ported to Go for self-learning purposes.

All 3 experiments run with multithreading for faster computation. All 3 experiments run for 100 trails per test scenario in order to obtain statistical reliability. All graphs are generated with Python.

## Requirements

* Go 1.15+
* Python 3.6+

## How to Build

* Build all 3 experiments
    ```txt
    make
    ```

* Build individual experiments
    ```txt
    make exp01
    make exp02
    make exp03
    ```

## How to Run

* cd into the bin directory
    ```txt
    ./exp01
    ./exp02
    ./exp03
    ```

## How to Generate Graphs

* Install Python dependencies
    ```txt
    python3 -m pip install -r requirements.txt
    ```

* cd into the graph directory
    ```txt
    python3 graph_exp01.py
    python3 graph_exp02.py
    python3 graph_exp03.py
    ```

## Project Setup

```txt
.
├── bin                 <-- Go binaries
│   ├── exp01
│   ├── exp02
│   └── exp03
├── cmd                 <-- Experiment 1, 2, 3 Go code
│   ├── exp01
│   │   └── main.go
│   ├── exp02
│   │   └── main.go
│   └── exp03
│       └── main.go
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
├── pkg                 <-- Shared Go code
│   ├── recorder.go
│   ├── stats.go
│   └── trace.go
├── README.md
├── res                 <-- Other resources
└── results             <-- Experiment 1, 3, 3 results
    ├── exp01
    ├── exp02
    └── exp03
```