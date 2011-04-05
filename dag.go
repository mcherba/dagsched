/***************************************************************************
 * code for handling Directed Acyclic Graphs for Dist Scheduling hw1
 * Mike Cherba mcherba@gmail.com 04-04-11
 ***************************************************************************/
package dag
import fmt "fmt" // Package implementing formatted I/O.

const maxnodes = 52 // this is set to 52 to allow for an additional entry and exit node for algorithmic simplicity 

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
	children [maxnodes]ChildEntry	  	
}

type Dag struct {
	g [maxnodes]DagNode;
}

var thisDag Dag


// DAG functions
//--------------------------------------------------------------------
/***************************************************************************
 * Initialize an the empty DAG.  Takes nothing, currently returns nothing, 
 * though it should return an error.
 ***************************************************************************/
func init () () {
	var i uint = 0
	var j uint = 0
	
	for i=0; i < maxnodes; i++ {
		thisDag.g[i].id=-1
		thisDag.g[i].weight=-1
		thisDag.g[i].level=-1
		for j=0; j < maxnodes; j++ {
			thisDag.g[i].children[j].id=-1
			thisDag.g[i].children[j].ccost=-1
		}
	}
	return
	
}

/***************************************************************************
 * Initialize aSample DAG for use testing algrothms 
 ***************************************************************************/
func Itest () () {
	
	thisDag.g[0].id=0
	thisDag.g[0].weight = 0
	thisDag.g[0].level = 0
	thisDag.g[0].children[0].id = 1
	thisDag.g[0].children[0].ccost = 0
	thisDag.g[0].children[1].id = 2
	thisDag.g[0].children[1].ccost = 0
	thisDag.g[1].id=1
	thisDag.g[1].weight = 2053741237
	thisDag.g[1].level = 1
	thisDag.g[1].children[0].id = 5
	thisDag.g[1].children[0].ccost = 2376224
	thisDag.g[2].id=2
	thisDag.g[2].weight = 1073741824
	thisDag.g[2].level = 1
	thisDag.g[2].children[0].id = 6
	thisDag.g[2].children[0].ccost = 8388608
	thisDag.g[3].id=5
	thisDag.g[3].weight = 2053741237
	thisDag.g[3].level = 2
	thisDag.g[3].children[0].id = 7
	thisDag.g[3].children[0].ccost = 0
	thisDag.g[4].id=6
	thisDag.g[4].weight = 2053741237
	thisDag.g[4].level = 2
	thisDag.g[4].children[0].id = 7
	thisDag.g[4].children[0].ccost = 0
	thisDag.g[5].id=7
	thisDag.g[5].weight = 0
	thisDag.g[5].level = 3
	
	//Dag.g.children[0
	return 
}

/***************************************************************************
 * Print the Active DAG.  But only print elements that we've put something in. 
 ***************************************************************************/
func PrintDag () () {
	var i uint = 0
	for i =0; i < maxnodes; i++ {
			if thisDag.g[i].id >= 0 {
				printNode(thisDag.g[i])
			}
		}
	//fmt.Printf("%v\n", thisDag.g[0])
	return
}

/***************************************************************************
 * Print a single node in the DAG.  Only print Child Nodes we've used, not the 
 * Whole array. 
 ***************************************************************************/
func printNode (ln DagNode) () {
	var i uint = 0
	
	if ln.id != -1 {
		fmt.Printf("%d, %d, %d [ ", ln.id, ln.weight, ln.level)
		for i =0; i < maxnodes; i++ {
			if ln.children[i].id >= 0 {
				fmt.Printf("%v ", ln.children[i])
			}
		}
		fmt.Printf("]\n")
	}
}
