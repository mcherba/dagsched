package parser

import fmt "fmt"
import "io/ioutil"
import "os"
import "strconv"
import vector "container/vector"

// method to be called by main
func ParseFile(fName string){
	
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
	
	// get node count	
	for i:=bIdx; i<len(buf); i++ {
		
		if buf[i] != space && buf[i] != nL {
			s += string(buf[i])
		} else {
			if s != "" {
				nCount, _ = strconv.Atoi(s)
				fmt.Printf("node count: %d\n", nCount)
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
							case 4: nT.cc, _ = strconv.Atoi(s)
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
	for i:=0; i<nCount; i++ {
		(nTArray.At(i).(*nodeTemp)).printNT()
	}
	
	
}

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
	cc int
}

// another practice func, but checks correctness
func (n *nodeTemp) printNT() {
	fmt.Printf("ID: %d, TY: %s, CC: %d, CL: ", n.id, n.ty, n.cc)
	for i:=0; i<(n.cl).Len(); i++ {
		fmt.Printf("%d ", (n.cl.At(i)).(int))
	}
	fmt.Printf("\n")
}

	
	
	
	
	
	
	
	
	
	
	
