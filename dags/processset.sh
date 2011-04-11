#!/bin/bash
pushd $1
mv -f $1_$3_$2.csv $1_$3_$2.csv.old
for n in {2..10}
do
	nn=$(($n * 5))
	dname=$(printf "%.2dtasks" "$nn")
	pushd $dname
	for nt in {1..25}
	do
		fname=$(printf "%.3d.dag" "$nt")
		
	  ../../../dagsched -f $fname -a $2 -n $3 >> ../../$1_$3_$2.csv
	done
	popd
done
popd
