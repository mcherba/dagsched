package sorter

import fmt "fmt"
import vec "container/vector"
import par "./parser"

// returns a topologically sorted vector of Node's
func TopSort(dag vec.Vector, s byte) (vec.Vector) {

	sortDag := new(vec.Vector)
	tempDag := dag.Copy()
	destDag := tempDag.Copy()
	setVec  := new(vec.Vector)
	
	// t-level gets regular top sort
	if s == 't' {
		setVec.Push(tempDag.At(0))
		destDag.Delete(0)
		for i:=setVec.Len(); i>0; i=setVec.Len() {
			n := (setVec.Pop().(*par.Node)).Copy()
			sortDag.Push(n)
			for j:=0; j<(n.Cl).Len(); j++ {
				c := ((n.Cl).At(j).(*par.Rel)).Id
				for k:=0; k<destDag.Len(); k++ {
					if (destDag.At(k).(*par.Node)).Id == c {
						for l:=0; l<(destDag.At(k).(*par.Node)).Pl.Len(); l++ {
							if (destDag.At(k).(*par.Node)).Pl.At(l).(*par.Rel).Id == n.Id {
								(destDag.At(k).(*par.Node)).Pl.Delete(l)
								break
							}
						}
					}
				}
			}
			for j:=0; j<destDag.Len(); j++ {
				if (destDag.At(j).(*par.Node)).Pl.Len() == 0 {
					c := (destDag.At(j).(*par.Node)).Id
					for k:=0; k<tempDag.Len(); k++ {
						if (tempDag.At(k).(*par.Node)).Id == c {
							setVec.Push(tempDag.At(k))
							break
						}
					}
					destDag.Delete(j)
					j--
				}
			}	
		}
	// b-level gets reverse top sort
	} else if s == 'b' {
		setVec.Push(tempDag.At(tempDag.Len()-1))
		destDag.Delete(destDag.Len()-1)
		for i:=setVec.Len(); i>0; i=setVec.Len() {
			n := (setVec.Pop().(*par.Node)).Copy()
			sortDag.Push(n)
			for j:=0; j<(n.Pl).Len(); j++ {
				c := ((n.Pl).At(j).(*par.Rel)).Id
				for k:=0; k<destDag.Len(); k++ {
					if (destDag.At(k).(*par.Node)).Id == c {
						for l:=0; l<(destDag.At(k).(*par.Node)).Cl.Len(); l++ {
							if (destDag.At(k).(*par.Node)).Cl.At(l).(*par.Rel).Id == n.Id {
								(destDag.At(k).(*par.Node)).Cl.Delete(l)
								break
							}
						}
					}
				}
			}
			for j:=0; j<destDag.Len(); j++ {
				if (destDag.At(j).(*par.Node)).Cl.Len() == 0 {
					c := (destDag.At(j).(*par.Node)).Id
					for k:=0; k<tempDag.Len(); k++ {
						if (tempDag.At(k).(*par.Node)).Id == c {
							setVec.Push(tempDag.At(k))
							break
						}
					}
					destDag.Delete(j)
					j--
				}
			}	
		}
	} else {
		fmt.Printf("Error")
	}

	return sortDag.Copy()

}
	
