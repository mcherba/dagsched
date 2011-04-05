/***************************************************************************
 * code for handling Directed Acyclic Graphs for Dist Scheduling hw1
 * Mike Cherba mcherba@gmail.com 04-04-11
 ***************************************************************************/
package dag
import fmt "fmt" // Package implementing formatted I/O.

// define a struct to store a child entry in our DAG
type ChildEntry struct {
	id int 		// the numerical id of the child node
	ccost int // the cost of communicating with this child 
	
}

// define a struct to store one entry in our DAG
type DagNode struct {
	id int 			// numerical id for the node
	weight int 	// the weight of the node
	
	/* the level of the node, from top or bottom dep on whether we are running 
	 * t-level or b-level */
	level int
 	
	/* an array of child pointers big enough to hold the most elements we use in 
	 * our test program
	 */
	children [50]int	  	
}

type Dag struct {
	g [50]DagNode;
}

var thisDag Dag

// DAG functions
//--------------------------------------------------------------------
func Itest () () {
	
	thisDag.g[0].id=1
	thisDag.g[0].weight=0
	thisDag.g[0].level=0
	//Dag.g.children[0
	return 
}

func PrintDag () () {
	fmt.Printf("%v\n", thisDag.g[0])
	return
}
