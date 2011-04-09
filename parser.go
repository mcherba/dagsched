package parser

import fmt "fmt"
import "io/ioutil"
import "os"
import "strconv"
import vector "container/vector"

// returns a vector representing a dag, each element is of type Node (see "type Node struct..." below
func ParseFile(fName string) (vector.Vector){
	
	// get file contents
	buf, err := ioutil.ReadFile(fName)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	
	// some things
	var   s      string = ""
	var   bIdx   int    = 0
	var   nCount int    = 0
	const space  byte   = ' '
	const null   byte   = '-'
	const nL     byte   = '\n'
	nTArray := new(vector.Vector)
	
	// get to node count
	for i:=0; i<len(buf); i++ {
	
		if buf[i] != space && buf[i] != nL {
			s += string(buf[i])
		} else {
			if checkNC(s) {
				s = ""
				bIdx = i
				break
			} else {
				s = ""
			}
		}
	}
	
	// get the node count	
	for i:=bIdx; i<len(buf); i++ {
		
		if buf[i] != space && buf[i] != nL {
			s += string(buf[i])
		} else {
			if s != "" {
				nCount, _ = strconv.Atoi(s)
				//fmt.Printf("node count: %d\n", nCount)
				s = ""
				bIdx = i
				break
			}
		}
	}	
	
	// for each node, skip, get id int, childList vector, compCost int, skip
	for i:=0; i<nCount; i++ {
		nT := new(nodeTemp)
		for j:=0; j<6; j++ {
			for k:=bIdx; k<len(buf); k++ {
				if buf[k] != space && buf[k] != nL {
					s += string(buf[k])
				} else {
					if s != "" {
						switch j {
							case 0: // NODE - unused
							case 1: nT.id, _ = strconv.Atoi(s)
							case 2: nT.cl    = getCVec(s)
							case 3: nT.ty    = s
							case 4: nT.cc, _ = strconv.Atoi64(s)
							case 5: // PAR COST - unused
						}
						s = ""
						bIdx = k
						break
					}
				}
			}
		}
		nTArray.Push(nT)
	}
	
	// check for correctness
	//for i:=0; i<nCount; i++ {
	//	(nTArray.At(i).(*nodeTemp)).printNT()
	//}
	
	// vector to be returned, note that child lists are set up
	// in the following for loops (indexed by 'i')
	dagVec := new(vector.Vector)
	
	// push root
	for i:=0; i<nCount; i++ {
		temp := (nTArray.At(i).(*nodeTemp))
		if temp.ty == "ROOT" {
			n    := new(Node)
			n.Id  = 0
			n.Ty  = "ROOT"
			n.Ex  = 0
			n.Lev = -1
			for j:=0; j<(temp.cl).Len(); j++ {
				r   := new(Rel)
				r.Id = ((temp.cl).At(j)).(int)
				r.Cc = 0
				(n.Cl).Push(r)
			}
			dagVec.Push(n)
			break
		}
	}
	
	// push computation nodes
	for i:=0; i<nTArray.Len(); i++ {
		temp := nTArray.At(i).(*nodeTemp)
		if temp.ty == "COMPUTATION" {
			n    := new(Node)
			n.Id  = temp.id
			n.Ty  = temp.ty
			n.Ex  = temp.cc
			n.Lev = -1
			for j:=0; j<(temp.cl).Len(); j++ {
				r   := new(Rel)
				tId := ((temp.cl).At(j)).(int)
				for k:=0; k<nTArray.Len(); k++ {
					tNo := nTArray.At(k).(*nodeTemp)
					if tNo.id == tId {
						if tNo.ty == "TRANSFER" {
							r.Id = (tNo.cl).At(0).(int)
							r.Cc = tNo.cc
						} else if tNo.ty == "END" {
							r.Id = tNo.id
							r.Cc = 0
						}
					break
					}
				}
				(n.Cl).Push(r)
			}
			dagVec.Push(n)
		}
	}
	
	// push end node (no child list)
	for i:=0; i<nTArray.Len(); i++ {
		temp := nTArray.At(i).(*nodeTemp)
		if temp.ty == "END" {
			n    := new(Node)
			n.Id  = temp.id
			n.Ty  = temp.ty
			n.Ex  = 0
			n.Lev = -1
			dagVec.Push(n)
			break
		}
	}
		
	// set up parent lists
	for i:=0; i<dagVec.Len(); i++ {
		tId  := (dagVec.At(i).(*Node)).Id
		for j:=0; j<dagVec.Len(); j++ {
			temp := dagVec.At(j).(*Node)
			for k:=0; k<(temp.Cl).Len(); k++ {
				relT := (temp.Cl).At(k).(*Rel)
				if relT.Id == tId {
					r := new(Rel)
					r.Id = temp.Id
					r.Cc = relT.Cc
					(dagVec.At(i).(*Node)).Pl.Push(r)
				}
			}
		}
	}		
	
	// return dag vector
	return dagVec.Copy()
}

// Node structure, the elements of the vectors cl and pl are of type Rel, see "type Rel struct..." below
type Node struct {
	Id int
	Ty string
	Ex int64
	Cl vector.Vector
	Pl vector.Vector
	Lev int64
}

// Rel(ative) structure
type Rel struct {
	Id int
	Cc int64
}

// Copy returns a copy of a node
func (n *Node) Copy() (*Node) {
	t    := new(Node)
	t.Id  = n.Id
	t.Ty  = n.Ty
	t.Ex  = n.Ex
	t.Cl  = (n.Cl).Copy()
	t.Pl  = (n.Pl).Copy()
	t.Lev = n.Lev
	return t
}

// PrintNode prints a Node
func (n *Node) PrintNode() {
	fmt.Printf("NODE:       %d\n type:      %s\n exTime:    %d\n level:     %d\n", n.Id, n.Ty, n.Ex, n.Lev)
	p := n.Pl
	c := n.Cl
	fmt.Printf(" parList:   ")
	for j:=0; j<p.Len(); j++ {
		r := p.At(j).(*Rel)
		fmt.Printf("(%d, %d) ", r.Id, r.Cc)
	}
	fmt.Printf("\n childList: ")
	for j:=0; j<c.Len(); j++ {
		r := c.At(j).(*Rel)
		fmt.Printf("(%d, %d) ", r.Id, r.Cc)
	}
	fmt.Printf("\n")
}

// PrintDAG prints a dag represented by v
func PrintDAG(v vector.Vector) {
	for i:=0; i<v.Len(); i++ {
		(v.At(i).(*Node)).PrintNode()
	}
}

// *****
// Helper functions
// *****

// made to get some practice writing func's
func checkNC(s string) bool {
	if s == "NODE_COUNT" { return true }
	return false
}

// turn string of children id's into a vector
func getCVec(s string) vector.Vector {
	var t string = ""
	v := new(vector.Vector)
	for i:=0; i<len(s); i++ {
		if string(s[i]) != "," {
			t += string(s[i])
		} else {
			if t != "" {
				r, _ := strconv.Atoi(t)
				v.Push(r)
				t = ""
			}
		}
	}
	if t != "" {
		r, _ := strconv.Atoi(t)
		v.Push(r)
	}
	return v.Copy()
}

// intermediate node structure
type nodeTemp struct {
	id int
	cl vector.Vector
	ty string
	cc int64
}

// another practice func, but checks correctness
func (n *nodeTemp) printNT() {
	fmt.Printf("ID: %d, TY: %s, CC: %d, CL: ", n.id, n.ty, n.cc)
	for i:=0; i<(n.cl).Len(); i++ {
		fmt.Printf("%d ", (n.cl.At(i)).(int))
	}
	fmt.Printf("\n")
}


	
// get the index of a node in a vector given nodes ID
func GetIndexById (v vector.Vector,  id int) (int) {
	for i:=0; i<v.Len(); i++ {
		if v.At(i).(*Node).Id == id {
			return i
		}
	}
	return -1
}

// reset levels of a node to -1
func (n *Node) ResetLevel () {
	n.Lev = -1
}
	
	
	
	
	
	
	
	
	
