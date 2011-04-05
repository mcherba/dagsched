#!/bin/bash
mkdir $1
pushd $1
for n in {2..10}
do
	nn=$(($n * 5))
	dname=$(printf "%.2dtasks" "$nn")
	mkdir $dname
	pushd $dname
	for nt in {1..25}
	do
		fname=$(printf "%.3d.dag" "$nt")
		daggen -n $nn --mindata 10 --maxdata 1500 --fat .73  --density .37 regular .2 --ccr 3 --jump 3 --maxalpha 0 > $fname
	done
	popd
done
popd
