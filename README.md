# log-k-decomp
[![](https://img.shields.io/github/v/release/cem-okulmus/log-k-decomp)](https://github.com/cem-okulmus/log-k-decomp/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/cem-okulmus/log-k-decomp.svg)](https://pkg.go.dev/github.com/cem-okulmus/log-k-decomp)
[![Go Report Card](https://goreportcard.com/badge/github.com/cem-okulmus/log-k-decomp?)](https://goreportcard.com/report/github.com/cem-okulmus/log-k-decomp)

log-k-decomp implements a novel parallel algorithm to compute Hypertree Decompositions based on the structural information of CQs or CSPs. This can then be used to evaluate them in provably polynomial time.


## How to build 
Needs Go 1.14 to be installed first. Files to install it for Linux, macOS or Windows can be found here: <https://go.dev/dl/>. 

Command to produce exectuable: `go build` 

## Using the command line tool
Run `./log-k-decomp -h` to see currently supported command and options. Hypergraphs need to be encoded in HyperBench format, more info here: <http://hyperbench.dbai.tuwien.ac.at/downloads/manual.pdf>.

Only the '-graph' and '-width' flags need to be specified for a run, though the tool provides plenty of customisation options, ranging from providing additional logs to subtle modifications to the underlying algorithm. For detailed information on the log-k-decomp algorith, we refer to the paper. 


## Publication

[[1]](https://dl.acm.org/doi/abs/10.1145/3517804.3524153) G. Gottlob, M. Lanzinger, C. Okulmus, R. Pichler: Fast Parallel Hypertree Decompositions in Logarithmic Recursion Depth. Proceedings of the 41st ACM SIGMOD-SIGACT-SIGAI Symposium on Principles of Database Systems, (PODS), June 2022 
