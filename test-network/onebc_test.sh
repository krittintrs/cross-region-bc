#!/bin/bash

data=$1

echo -e "index\tisFound\ttime"

./onebc_execute.sh 1
./onebc_execute.sh $(($data/10*2))
./onebc_execute.sh $(($data/10*5))
./onebc_execute.sh $(($data/10*7))
./onebc_execute.sh $data
./onebc_execute.sh 1
./onebc_execute.sh $(($data/10*2))
./onebc_execute.sh $(($data/10*5))
./onebc_execute.sh $(($data/10*7))
./onebc_execute.sh $data
./onebc_execute.sh $(($data/10*12))
./onebc_execute.sh $(($data/10*15))
./onebc_execute.sh $(($data/10*18))
./onebc_execute.sh $(($data/10*20))