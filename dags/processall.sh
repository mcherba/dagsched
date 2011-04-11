#!/bin/bash
for n in {2..5}
do
	./processset.sh $1 tl $n
	./processset.sh $1 bl $n
	./processset.sh $1 c $n
done	
