#!/bin/bash

rm -r logs >/dev/null 2>&1
mkdir logs

pkill -f server 2> /dev/null

cd ../server
go build

pkill -f server 2> /dev/null

counter=0
n=${1-7}
debug=${2-1}
while [ $counter -lt $n ]
do
	let port=6000+$counter
	name='Passtor'$counter
	nohup tilix -t $name -s $name --window-style=disable-csd-hide-toolbar -e bash -c './server -name '$name' -addr 127.0.0.1:'$port' -peers 127.0.0.1:6000 -v '$debug'; exec bash' > /dev/null 2>&1 &
	((counter++))
done
