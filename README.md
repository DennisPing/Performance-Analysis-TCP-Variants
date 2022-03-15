# Performance-Analysis-TCP-Variants

Matthew Jones  
Dennis Ping

## Background

This project was initially done in Python and was ported to Go for self-learning purposes.

All 3 experiments run with multithreading for faster computation. All 3 experiments run for 100 trails per test scenario in order to obtain statistical reliability. The graphs are generated with Python.

## Requirements

* Go 1.15+
* Python 3.7+
    ```
    matplotlib
    numpy
    pandas
    ```

## How to Build

* Build all 3 experiments
    ```
    make
    ```

* Build individual experiments
    ```
    make exp01
    make exp02
    make exp03
    ```

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
├── pkg                 <-- shared packages
│   ├── stats.go
│   └── trace.go
├── README.md
└── results             <-- Experiment results
    ├── exp01
    ├── exp02
    └── exp03
```