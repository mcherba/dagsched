/***************************************************************************
 * DAG Scheduler simulator, Mike and Rick UofO CIS-410 Distributed Scheduling
 *
 ***************************************************************************/
package main

import fmt "fmt" // Package implementing formatted I/O.
import flag "flag" // Command line parsing
import parser "./parser" // structure and code for DAGs
import sorter "./sorter" // topological sort
import vec "container/vector"
func main() { 
	
	// Command line flags	
	var numcores *int = flag.Int("n", 2, "number of cores to use in the simulation [-n Int value]")
	var infname *string = flag.String("f", "infile.dag", "filename to load the .dag from [-f filename.dag]")
	var algtype *string = flag.String("a", "t-level", "algorithm type to use t-level, b-level, or ??")
	

	flag.Parse()

	fmt.Printf("simulating using %d cores\n", *numcores)
	fmt.Printf("loading DAG from %s\n", *infname)
	fmt.Printf("using %s scheduling algorithm\n", *algtype)
	
	
	// Read in the dag we want to schedule
	var dag = parser.ParseFile(*infname)
	//parser.PrintDAG(dag)
	
	// topo sort it
	//var tsdag = sorter.TopSort(dag, 't')
	//parser.PrintDAG(tsdag)
	//fmt.Printf("%v", tsdag.At(0).Id)
	tlevel(dag)
	fmt.Printf("-------------------\n\n")
	tlevelsched(dag, 2)
	
	slist:=sorter.TSort(dag)
	fmt.Printf("\n\n%v\n\n, last=%d", slist, slist.Last())
	
} 

type Event struct {
	id int
	start int64
	end int64
}

func tlevel (indag vec.Vector) () {
	var TopList = sorter.TopSort(indag, 't')
	var max int64
	parser.PrintDAG(TopList)
	//initialize the level of the root node to 0
	(TopList.At(0).(*parser.Node)).Lev = 0
	
	// for each node in the sorted list
	for i:=1; i<len(TopList); i++ {
		max=0
		// for each parent node of the present node
		for j:=0; j<len((TopList.At(i).(*parser.Node)).Pl); j++ {
			nodeID:=(TopList.At(i).(*parser.Node)).Id
			pId:= (TopList.At(i).(*parser.Node)).Pl.At(j).(*parser.Rel).Id
			linkW:= (TopList.At(i).(*parser.Node)).Pl.At(j).(*parser.Rel).Cc
			pIndex:=parser.GetIndexById(TopList, pId)
			pLevel:=(TopList.At(pIndex).(*parser.Node)).Lev
			pCost:=(TopList.At(pIndex).(*parser.Node)).Ex
			fmt.Printf("(nodeId %d: i %d, j %d, pId %d, pLevel %d, parent index:%d)\t", nodeID, i, j,pId, pLevel, pIndex)
			fmt.Printf(" Link to parent cost: %d ", pCost)
			fmt.Printf(" linkw %d, pCost %d,  cp %d \n", linkW, pCost, pLevel + linkW +  pCost)
			if  ( pLevel + linkW +  pCost) > max {
				max = pLevel + linkW +  pCost
			}
		}
		(TopList.At(i).(*parser.Node)).Lev = max
		fmt.Printf("(%d)\n", max) 
	}
}

// schedule using t-level Earliest Start time 1st 
func tlevelsched (indag vec.Vector, ncpus int) () {
	//cpu:= new(vec.Vector)  // holds the cpu assigned to a task
	cpu:=make([]int, len(indag))
	
	var TopList = 	sorter.TopSort(indag, 't')
	var max int64
	var ncpu int
	var pCost int64

	
	//initialize the level of the root node to 0
	(TopList.At(0).(*parser.Node)).Lev = 0
	//always schedule the root node on cpu 0
	cpu[0]=0
	// for each node in the sorted list
	for i:=1; i<len(TopList); i++ {
		max=0
		ncpu=0
		// for each parent node of the present node
		for j:=0; j<len((TopList.At(i).(*parser.Node)).Pl); j++ {
			

			linkW:= (TopList.At(i).(*parser.Node)).Pl.At(j).(*parser.Rel).Cc
			
			pId:= (TopList.At(i).(*parser.Node)).Pl.At(j).(*parser.Rel).Id
			pIndex:=parser.GetIndexById(TopList, pId)
			pLevel:=(TopList.At(pIndex).(*parser.Node)).Lev
			
			for k:=0; k < ncpus; k++ {
				if k != cpu[pIndex] {	
					pCost=(TopList.At(pIndex).(*parser.Node)).Ex
				} else {
					pCost=0
				}
				if  ( pLevel + linkW +  pCost) > max {
					max = pLevel + linkW +  pCost
					ncpu = k
				}
			}
			
		}
		(TopList.At(i).(*parser.Node)).Lev = max
		cpu[i]=(ncpu)
		
		nodeID:=(TopList.At(i).(*parser.Node)).Id
		fmt.Printf("Scheduling %d on cpu%d with t-level %d\n", nodeID, ncpu, max) 
	}
	
}
