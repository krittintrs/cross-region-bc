#!/bin/bash

data=$1

echo "[ INDEXER - $data data ]"

# regionalCC1
./indexer_execute.sh HP1 1
./indexer_execute.sh HP3 $(($data/10*2))
./indexer_execute.sh HP4 $(($data/10*5))
./indexer_execute.sh HP1 $(($data/10*10))
./indexer_execute.sh HP3 $(($data/10*15))

# regionalCC2
./indexer_execute.sh HP2 1
./indexer_execute.sh HP5 $(($data/10*2))
./indexer_execute.sh HP8 $(($data/10*3))
./indexer_execute.sh HP2 $(($data/10*6))
./indexer_execute.sh HP5 $(($data/10*12))

# regionalCC3
./indexer_execute.sh HP6 1
./indexer_execute.sh HP7 $(($data/10*2))
./indexer_execute.sh HP6 $(($data/10*4))
./indexer_execute.sh HP7 $(($data/10*8))

# invalid HP case
./indexer_execute.sh HP9 1
./indexer_execute.sh HP10 $data