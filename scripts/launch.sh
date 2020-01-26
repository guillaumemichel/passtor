#!/bin/bash

rm -r logs >/dev/null 2>&1
mkdir logs

cd ../server
go build

counter=0
n=${1-7}
debug=${2-1}
while [ $counter -lt $n ]
do
	let port=6000+$counter
	name='Passtor'$counter
	./server -name $name -addr 127.0.0.1:$port -peers 127.0.0.1:6000 -v $debug > '../scripts/logs/'$name'.txt' 2>&1 &
	((counter++))
done
