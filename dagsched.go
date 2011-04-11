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
import gt "./getTimes" // get time info from dags
func main() { 
	
	// Command line flags	
	var numcores *int = flag.Int("n", 2, 
		"number of cores to use in the simulation [-n Int value]")
	var infname *string = flag.String("f", "infile.dag", 
		"filename to load the .dag from [-f filename.dag]")
	var algtype *string = flag.String("a", "tl", 
		"algorithm type to use tl=t-level, bl=b-level, or ??")
	

	flag.Parse()

	//fmt.Printf("simulating using %d cores\n", *numcores)
	//fmt.Printf("loading DAG from %s\n", *infname)
	//fmt.Printf("using %s scheduling algorithm\n", *algtype)
	
	
	// Read in the dag we want to schedule
	var dag = parser.ParseFile(*infname)
	switch *algtype {
	case "tl": 
		// Schedule using t-level
		ScheduleTlevel(dag, *numcores)

	case "bl":
		// schedule using b-level
		ScheduleBlevel(dag, *numcores)

	}
	
	//fmt.Printf("fwd: %v\n", sorter.TSort (dag, 'f'))
	//fmt.Printf("rev: %v\n", sorter.TSort (dag, 'r'))	
	
	//fmt.Printf("SeqTime %v\n", gt.SeqTime(dag))
	//parser.PrintDAG(dag)
	//fmt.Printf("CPTime %v\n", gt.CPTime(dag))
	//parser.PrintDAG(dag)
} 

type Event struct {
	id int
	start int64
	end int64
}


// schedule using t-level Earliest Start time 1st 
func ScheduleTlevel(indag vec.Vector, ncpus int) (){
	var el int64
	var nt int
	var iccost bool
	// produce a topographically sorted DAG to work with
	var dag =	sorter.TopSort(indag, 't')
	cpu:=make([]int, len(dag)) // cpu is a slice as long as the dag
	for i:=0; i < len(dag); i++ {
		cpu[i]= -1
	}
	// Create a schedule as a vector of Event Vectors
	Schedule:= make([]vec.Vector, ncpus)

	//update the t-level of the root node
 	dag.At(0).(*parser.Node).Lev=0
 	//recursively update the t-levels of it's children
 	tlUpdateChildren(dag,0)
 	
 	
	startTime := make([]int64, ncpus) // holds current start times	
  // initialize all start times to 0
	for i:=0; i <ncpus; i++ { 
		startTime[i]=0
	}
	
	//iterate over the dag, till we've done every node.
	for i:=0; i < len(dag); i++ {
		var coreChosen = -1
		var earliestAvail int64
		var earliestFeasable int64
		var efc int
	  esT := make([]int64, ncpus)
	  pet :=new(vec.Vector)
	  pCpu :=new(vec.Vector)
	  iccost=false

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
	  // if none is obvious then we'll record some info and default to 0
	  if coreChosen == -1 {
	  	for c:=ncpus-1; c >= 0; c-- {
	  		if startTime[c] <= earliestAvail {
	  			earliestAvail = startTime[c]
	  			coreChosen = c
	  		}	
	  	}
	  }
	  
	  //see if we have to account for communications time
    //check each parent to see what cpu it was scheduled on
	  //eA=0
	  for j:=0; j < len(dag.At(nt).(*parser.Node).Pl); j++ {
	  	// set up scratch variables, otherwise the statement to capture 
	  	// what we want is too long and complex.
	  	cParentId:=dag.At(nt).(*parser.Node).Pl.At(j).(*parser.Rel).Id
	  	comC := dag.At(nt).(*parser.Node).Pl.At(j).(*parser.Rel).Cc
	  	cParent:=parser.GetIndexById(dag, cParentId)
	  	cpCore :=  cpu[cParent]
	  	cpIndex := findInSchedule(Schedule[cpCore], cParentId)
	  	
	  	// Here's what we want
	  	pet.Push(Schedule[cpCore].At(cpIndex).(*Event).end + comC)
	  	pCpu.Push(cpCore)
	  	
	  	// check and see if we need to account for comm costs.
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

	//	PrintSchedule(Schedule)			
	makespan:=findScheduleEnd(Schedule)
	fmt.Printf("%d,", (len(dag)-2))  // don't count the 0 time start and end tasks
	fmt.Printf("%d,", gt.SeqTime(dag))
	fmt.Printf("%d,", makespan)
	fmt.Printf("%d,", gt.CPTime(dag))
	fmt.Printf("%v\n", float64(SumExTime(Schedule)) / float64(makespan*int64(ncpus)))
}

func ScheduleBlevel(indag vec.Vector, ncpus int) (){
	var el int64
	var nt int
	var iccost bool
	// produce a topographically sorted DAG to work with
	var dag =	sorter.TopSort(indag, 'b')
	//update the b-level of the end node
 	dag.At(0).(*parser.Node).Lev=0
 	//recursively update the b-levels of it's parents
 	blUpdateChildren(dag,0)
 	
 	//parser.PrintDAG(dag)
	
 	cpu:=make([]int, len(dag)) // cpu is a slice as long as the dag
	for i:=0; i < len(dag); i++ {
		cpu[i]= -1
	}
 	// Create a schedule as a vector of Event Vectors
	Schedule:= make([]vec.Vector, ncpus)
 	startTime := make([]int64, ncpus) // holds current start times	
  // initialize all start times to 0
	for i:=0; i <ncpus; i++ { 
		startTime[i]=0
	}
 	
 	
 	// Iterate over the dag, each time selecting the node with the highest b-level
 	for i:=0; i < len(dag); i++ {
 		
 		var coreChosen = -1
		var earliestAvail int64
		var earliestFeasable int64
		var efc int
	  esT := make([]int64, ncpus)
	  pet :=new(vec.Vector)
	  pCpu :=new(vec.Vector)
	  iccost=false

		el=0 // largest signed integer
		earliestAvail=9223372036854775807
		earliestFeasable=9223372036854775807
		nt=0
				
		// Select the next task to Schedule
		// look through the dag for the task with the highest b-level
		for j:=0; j < len(dag); j++ {
			if (cpu[j] < 0) && (dag.At(j).(*parser.Node).Lev >= el) {
				nt=j
				el=dag.At(j).(*parser.Node).Lev
			}
		}
		

		// select a cpu to schedule it on
		/*
	  for c:=ncpus-1; c >= 0; c-- {
	  	if startTime[c] <= el {
	  		coreChosen = c
	  	}	
	  }
	  */
	  // if none is obvious then we'll record some info and default to 0
	  //if coreChosen == -1 {
	  	for c:=ncpus-1; c >= 0; c-- {
	  		if startTime[c] <= earliestAvail {
	  			earliestAvail = startTime[c]
	  			coreChosen = c
	  		}	
	  	}
	  //}
	  
	  //fmt.Printf("Task %d, Core chosen=%d\n",dag.At(nt).(*parser.Node).Id, coreChosen)
	  
	  //fmt.Printf("A\n")
	  //see if we have to account for communications time
    //check each parent to see what cpu it was scheduled on
	  //eA=0
	  if coreChosen == -1 {
	  	coreChosen=0
	  }
	  for j:=0; j < len(dag.At(nt).(*parser.Node).Pl); j++ {
	  	
	  	// set up scratch variables, otherwise the statement to capture 
	  	// what we want is too long and complex.
	  	cParentId:=dag.At(nt).(*parser.Node).Pl.At(j).(*parser.Rel).Id
	  	comC := dag.At(nt).(*parser.Node).Pl.At(j).(*parser.Rel).Cc
	  	cParent:=parser.GetIndexById(dag, cParentId)
	  	cpCore :=  cpu[cParent]
	  	/*
	  	fmt.Printf("Iteration %d ncpus=%v  cpCore=%v\n", j, ncpus, cpCore)
	  	if cpCore== -1 {
	  			PrintSchedule(Schedule)
	  	}
	  	fmt.Printf("cParent=%d, cParentId=%d\n", cParent, cParentId)
	  	fmt.Printf("Schedule[cpCore]=%v  cParentd=%v\n", Schedule[cpCore], cParentId)
	  	*/
	  	cpIndex := findInSchedule(Schedule[cpCore], cParentId)
	  	
	  	// Here's what we want
	  	pet.Push(Schedule[cpCore].At(cpIndex).(*Event).end + comC)
	  	pCpu.Push(cpCore)
	  	
	  	// check and see if we need to account for comm costs.
	  	if cpu[cParent]!= coreChosen {
	  			iccost=true
	  	}

	  }
	  
 	  //see if we have to account for communications time
    //check each parent to see what cpu it was scheduled on
	  //eA=0
	  /*
	  for j:=0; j < len(dag.At(nt).(*parser.Node).Cl); j++ {
	  	fmt.Printf("Task %d has child nodes=%d\n", dag.At(nt).(*parser.Node).Id, len(dag.At(nt).(*parser.Node).Cl))
	  	// set up scratch variables, otherwise the statement to capture 
	  	// what we want is too long and complex.
	  	cParentId:=dag.At(nt).(*parser.Node).Cl.At(j).(*parser.Rel).Id
	  	comC := dag.At(nt).(*parser.Node).Cl.At(j).(*parser.Rel).Cc
	  	cParent:=parser.GetIndexById(dag, cParentId)
	  	cpCore :=  cpu[cParent]
	  	fmt.Printf("ncpus=%v  cpCore=%v\n", ncpus, cpCore)
	  	fmt.Printf("Schedule[cpCore]=%v  cChildId=%v\n", Schedule[cpCore], cParentId)
	  	cpIndex := findInSchedule(Schedule[cpCore], cParentId)
	  	
	  	// Here's what we want
	  	pet.Push(Schedule[cpCore].At(cpIndex).(*Event).end + comC)
	  	pCpu.Push(cpCore)
	  	
	  	// check and see if we need to account for comm costs.
	  	if cpu[cParent]!= coreChosen {
	  			iccost=true
	  	}

	  }
*/
	  //fmt.Printf("B\n")
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
		//fmt.Printf("Core %d got Event: %v\n", coreChosen, ccEvent)

	  // update the available startTime for the core we chose
	  startTime[coreChosen] = ccEvent.end 

	  // schedule the task on a core
	  Schedule[coreChosen].Push(ccEvent)
	  cpu[nt]=coreChosen
	  
	  
	}	
	//PrintSchedule(Schedule)
	// Print the output
	makespan:=findScheduleEnd(Schedule)
	fmt.Printf("%d,", (len(dag)-2))  // don't count the 0 time start and end tasks
	fmt.Printf("%d,", gt.SeqTime(dag))
	fmt.Printf("%d,", makespan)
	fmt.Printf("%d,", gt.CPTime(dag))
	fmt.Printf("%v\n", float64(SumExTime(Schedule)) / float64(makespan*int64(ncpus)))
	//fmt.Printf("%v\n",  SumExTime(Schedule))
	//fmt.Printf("\n")
	
	//parser.PrintDAG(dag)
	
	//parser.PrintDAG(dag)
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

// returns a sum of the time 
func SumExTime (S []vec.Vector) (int64){

		var tsum int64
		for i:=0; i < len(S); i++ {
			if S[i].Len() > 0 { 
				tsum += S[i].Last().(*Event).end - S[i].At(0).(*Event).start
			}
		}
	  
		return tsum
}

// Returns the end time of the last task scheduled
func findScheduleEnd(S []vec.Vector) (int64){
	var en int64
	en=0
	for i:=0; i < len(S); i++ {
		if en < S[i].Last().(*Event).end {
			en = S[i].Last().(*Event).end
		}
	}
	return en
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

// given a dag and a node within that dag update the b-levels of all children 
// of that node
func blUpdateChildren(dag vec.Vector, nt int) {
	var max int64
	var cCost int64
	var linkW int64

	// for each node
	for jj:=0; jj<len(dag); jj++ {
		max=0
		// for each child node
		for kk:=0; kk<len((dag.At(jj).(*parser.Node)).Cl); kk++ {
			cIndex:=parser.GetIndexById(dag, 
				(dag.At(jj).(*parser.Node)).Cl.At(kk).(*parser.Rel).Id)
			cLevel:=(dag.At(cIndex).(*parser.Node)).Lev
			linkW= (dag.At(jj).(*parser.Node)).Cl.At(kk).(*parser.Rel).Cc
			cCost=(dag.At(cIndex).(*parser.Node)).Ex	
			if cCost + cLevel > max {
				max= cCost +cLevel
			}
		}
		//actually set the level
		dag.At(jj).(*parser.Node).Lev = max	+linkW
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

// The recursive helper function for blUpdateChildren
func blUpdate(dag vec.Vector, visited []bool, nt int) {
	var max int64
	var cCost int64
	var linkW int64
	max=0
	if visited[nt] {
		return
	}
	for j:=0; j < len((dag.At(nt).(*parser.Node)).Cl); j++ {
	  cIndex:=parser.GetIndexById(dag, 
			(dag.At(nt).(*parser.Node)).Cl.At(j).(*parser.Rel).Id)
		cLevel:=(dag.At(cIndex).(*parser.Node)).Lev
		linkW= (dag.At(nt).(*parser.Node)).Cl.At(j).(*parser.Rel).Cc
		cCost=(dag.At(cIndex).(*parser.Node)).Ex
		
		if  ( cLevel +  cCost) > max {
			max = cLevel + cCost 
		}
		
	}
		//actually set the level
		dag.At(nt).(*parser.Node).Lev = max + linkW	
		
	visited[nt]=true
	//recurse
	for jj:=0; jj<len((dag.At(nt).(*parser.Node)).Pl); jj++ {
		blUpdate(dag, visited, parser.GetIndexById(dag, 
			(dag.At(nt).(*parser.Node)).Pl.At(jj).(*parser.Rel).Id))
	}
}


