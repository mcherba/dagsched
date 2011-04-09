package getTimes

import (
	vec "container/vector"
	p "./parser"
)

// takes vector of nodes, returns sequential schedule length
func SeqTime (v vec.Vector) (int64) {
	var s int64
	s = 0
	for i:=0; i<v.Len(); i++ {
		s += (v.At(i).(*p.Node)).Ex
	}
	return s
}

// takes vector of nodes, returns critical path length
func CPTime (v vec.Vector) (int64) {
	var sLoc int
	sLoc = -1
	for i:=0; i<v.Len(); i++ {
		v.At(i).(*p.Node).ResetLevel()
		if v.At(i).(*p.Node).Ty == "ROOT" {
			sLoc = i
		}
	}
	
	setCPTime(&v, sLoc)
	
	return v.At(sLoc).(*p.Node).Lev
}

// helper for func CPTime	
func setCPTime (v *vec.Vector, s int) {
	var t *p.Node
	t = v.At(s).(*p.Node)	
	if t.Cl.Len() == 0 {
		t.Lev = t.Ex
	} else {
		var m int64
		m = 0
		for i:=0; i<t.Cl.Len(); i++ {
			tId := t.Cl.At(i).(*p.Rel).Id
			idx := p.GetIndexById(*v, tId)
			if v.At(idx).(*p.Node).Lev < 0 {
				setCPTime(v, idx)
			}
			temp := t.Ex + v.At(idx).(*p.Node).Lev
			if temp > m {
				m = temp
			}
		}
		t.Lev = m
	}
}
					
		
		
		
		
