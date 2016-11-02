#!/bin/bash

for i in $(seq 1 $1)
do 
    echo $i
    echo -e "hello $i \n" | nc localhost 1337 &
    
done
