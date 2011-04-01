package main

import fmt "fmt" // Package implementing formatted I/O.
import flag "flag" // Command line parsing

func main() { 
	
	// Command line flags	
	var numcores *int = flag.Int("n", 2, "number of cores to use in the simulation [-n Int value]")
	var infname *string = flag.String("f", "infile.dag", "filename to load the .dag from [-f filename.dag]")
	var algtype *string = flag.String("a", "t-level", "algorithm type to use t-level, b-level, or ??")

	flag.Parse()

	fmt.Printf("simulating using %d cores\n", *numcores)
	fmt.Printf("loading DAG from %s\n", *infname)
	fmt.Printf("using %s scheduling algorithm\n", *algtype)	
} 
