/***************************************************************************
 * DAG Scheduler simulator, Mike and Rick UofO CIS-410 Distributed Scheduling
 *
 ***************************************************************************/
package main

import fmt "fmt" // Package implementing formatted I/O.
import flag "flag" // Command line parsing
import parser "./parser" // structure and code for DAGs
//import sorter "./sorter" // topological sort
import vec "container/vector"
import gt "./getTimes" // get time info from dags
func main() { 
	
	// Command line flags	
	var numcores *int = flag.Int("n", 2, 
		"number of cores to use in the simulation [-n Int value]")
	var infname *string = flag.String("f", "infile.dag", 
		"filename to load the .dag from [-f filename.dag]")
	var algtype *string = flag.String("a", "t-level", 
		"algorithm type to use t-level, b-level, or ??")
	

	flag.Parse()

	fmt.Printf("simulating using %d cores\n", *numcores)
	fmt.Printf("loading DAG from %s\n", *infname)
	fmt.Printf("using %s scheduling algorithm\n", *algtype)
	
	
	// Read in the dag we want to schedule
	var dag = parser.ParseFile(*infname)
	
	// Schedule using t-level
	ScheduleTlevel(dag, *numcores)
	
	
	fmt.Printf("SeqTime %v\n", gt.SeqTime(dag))
	parser.PrintDAG(dag)
	fmt.Printf("CPTime %v\n", gt.CPTime(dag))
	parser.PrintDAG(dag)
} 

type Event struct {
	id int
	start int64
	end int64
}


// schedule using t-level Earliest Start time 1st 
func ScheduleTlevel(dag vec.Vector, ncpus int) (){
	var el int64
	var nt int
	var iccost bool
	cpu:=make([]int, len(dag)) // cpu is a slice as long as the dag
	for i:=0; i < len(dag); i++ {
		cpu[i]= -1
	}
	// Create a schedule as a vector of Event Vectors
	Schedule:= make([]vec.Vector, ncpus)
	// produce a topographically sorted DAG to work with
  //var dag =	&sorter.TopSort(*indag, 't')
  // we don't actually have to do a topo sort because all we use that for is
  // to get the first schedulable node, but we have that because we always have 
  // a root node of time 0, so we just start our t-levels from there.
	
	//update the t-level of the root node
 	dag.At(0).(*parser.Node).Lev=0
 	//recursively update the t-levels of it's children
 	tlUpdateChildren(dag,0)
 	
 	
	startTime := make([]int64, ncpus) // holds current start times	
	

	for i:=0; i <ncpus; i++ { // initialize all start times to 0
		startTime[i]=0
	}
	
	//cEvent.id=0
	//cEvent.start=0
	//cEvent.end=0
	
	// always schedule the start task onto cpu 0
	//Schedule[0].Push(cEvent)
	//cpu[0]=0


	for i:=0; i < len(dag); i++ {
		var coreChosen = -1
		var earliestAvail int64
		var earliestFeasable int64
		var efc int

		el=9223372036854775807 // largest signed integer
		earliestAvail=9223372036854775807
		earliestFeasable=9223372036854775807

		nt=0
		
		// Select the next task to Schedule
		// look through the dag for the task with the lowest t-level
		for j:=0; j < len(dag); j++ {
			if (cpu[j] < 0) && (dag.At(j).(*parser.Node).Lev < el) {
				nt=j
				el=dag.At(j).(*parser.Node).Lev
			}
		}
		
		// select a cpu to schedule it on
	  for c:=ncpus-1; c >= 0; c-- {
	  	if startTime[c] <= el {
	  		coreChosen = c
	  	}	
	  }
	  if coreChosen == -1 {
	  	for c:=ncpus-1; c >= 0; c-- {
	  		if startTime[c] <= earliestAvail {
	  			earliestAvail = startTime[c]
	  			
	  			coreChosen = c
	  		}	
	  	}
	  }

	  

	  esT := make([]int64, ncpus)
	  pet :=new(vec.Vector)
	  pCpu :=new(vec.Vector)
	  iccost=false
	  
	  //see if we have to account for communications time
    //check each parent to see what cpu it was scheduled on
	  //eA=0
	  for j:=0; j < len(dag.At(nt).(*parser.Node).Pl); j++ {
	  	cParentId:=dag.At(nt).(*parser.Node).Pl.At(j).(*parser.Rel).Id
	  	comC := dag.At(nt).(*parser.Node).Pl.At(j).(*parser.Rel).Cc
	  	cParent:=parser.GetIndexById(dag, cParentId)
	  	cpCore :=  cpu[cParent]
	  	cpIndex := findInSchedule(Schedule[cpCore], cParentId)
	  	
	  	pet.Push(Schedule[cpCore].At(cpIndex).(*Event).end + comC)
	  	pCpu.Push(cpCore)
	  	
	  	if cpu[cParent]!= coreChosen {
	  			iccost=true
	  	}

	  }

	  //find the core where the parentend+comm cost is least
	  efc=0
	  for c:=0; c < ncpus; c++ {
	  	esT[c]=startTime[c]
	  	 
	  	for p:=0; p < pCpu.Len(); p++ {
	  		if pCpu.At(p) != c {
	  			if pet.At(p).(int64) > esT[c] {
	  				esT[c] = pet.At(p).(int64)
	  			}
	  		}
	  	}
	  	if esT[c] < earliestFeasable {
	  		earliestFeasable = esT[c]
	  		efc=c
	  	}
	  }
	   
		el=startTime[coreChosen]
		// if we have to account for comm overhead chose this core
		if iccost {
	    el = earliestFeasable
			coreChosen = efc
	  }
	  
		// prepare the event
		ccEvent:=new(Event)
		ccEvent.id=dag.At(nt).(*parser.Node).Id
		ccEvent.start=el
		ccEvent.end=el+dag.At(nt).(*parser.Node).Ex

	  // update the available startTime for the core we chose
	  startTime[coreChosen] = ccEvent.end 

	  // schedule the task on a core
	  Schedule[coreChosen].Push(ccEvent)
	  cpu[nt]=coreChosen
	  
	  
	  if dag.At(nt).(*parser.Node).Lev!=el {
	  	//update the t-level of the current node
	  	dag.At(nt).(*parser.Node).Lev=el
	  	//update the t-levels of any children
	  	tlUpdateChildren(dag,nt)
	  }
	  
	  
	}

	PrintSchedule(Schedule)		
}

// Given the id of a task that has been scheduled on a processor
// return the index of that Task's Event entry
func findInSchedule(v vec.Vector, id int) (int){
	for i:=0; i<len(v); i++ {
		if v.At(i).(*Event).id == id {
				return i
		}
		
	}
	return -1
}

// Print out a schedule
// needed to verify that we're getting real schedules out of this
func PrintSchedule(Schedule []vec.Vector) () {
	for i:=0; i< len(Schedule); i++ {
		fmt.Printf("Schedule for core%d: %d tasks\n", i,len(Schedule[i]) )
		for j:=0; j < len(Schedule[i]); j++ {
			fmt.Printf("task %d, from %d to %d\n", Schedule[i].At(j).(*Event).id, 
				Schedule[i].At(j).(*Event).start, Schedule[i].At(j).(*Event).end)
		}
	}
	
}

// given a dag and a node within that dag update the t-levels of all children 
// of that node
func tlUpdateChildren(dag vec.Vector, nt int) {
	visited:=make([]bool, len(dag))
	for ii:=0; ii < len(dag); ii++ {
		visited[ii]=false
	}
		//recurse
	for jj:=0; jj<len((dag.At(nt).(*parser.Node)).Cl); jj++ {
		if !(visited[parser.GetIndexById(dag, 
			(dag.At(nt).(*parser.Node)).Cl.At(jj).(*parser.Rel).Id)]) {
				tlUpdate(dag, visited, parser.GetIndexById(dag, 
					(dag.At(nt).(*parser.Node)).Cl.At(jj).(*parser.Rel).Id))
		}
	}
}

// The recursive helper function for tlUpdateChildren
func tlUpdate(dag vec.Vector, visited []bool, nt int) {
	var max int64
	var pCost int64
	max=0
	if visited[nt] {
		return
	}
	for j:=0; j < len((dag.At(nt).(*parser.Node)).Pl); j++ {
		pIndex:=parser.GetIndexById(dag, 
			(dag.At(nt).(*parser.Node)).Pl.At(j).(*parser.Rel).Id)
		pLevel:=(dag.At(pIndex).(*parser.Node)).Lev
		linkW:= (dag.At(nt).(*parser.Node)).Pl.At(j).(*parser.Rel).Cc
		pCost=(dag.At(pIndex).(*parser.Node)).Ex
		
		if  ( pLevel + linkW +  pCost) > max {
			max = pLevel + linkW +  pCost 
		}
	}
		//actually set the level
		dag.At(nt).(*parser.Node).Lev = max	
		
	visited[nt]=true
	//recurse
	for jj:=0; jj<len((dag.At(nt).(*parser.Node)).Cl); jj++ {
		tlUpdate(dag, visited, parser.GetIndexById(dag, 
			(dag.At(nt).(*parser.Node)).Cl.At(jj).(*parser.Rel).Id))
	}
}
