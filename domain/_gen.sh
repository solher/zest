#!/bin/bash

for i in `ls ./*.go`; do sed 's/package domain/package ressources/g' $i > `echo $i | sed "s/.go/_temp.go/g"`; done
mv *_temp.go ../ressources
cd ../ressources
gen -f
rm ./*_temp.go
