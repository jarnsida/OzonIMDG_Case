#!/bin/bash
if [ -n "$1" ]
then
export PORT=$1
else
export PORT=8080
fi

# Max elements limit
if [ -n "$2" ]
then
export MAX_MEMORY=$2
else
export MAX_MEMORY=100000000
fi


# Connection close time out
if [ -n "$3" ]
then
export CON_CLOSE_TO=$3
else
export CON_CLOSE_TO=5
fi



# Logger settings
export LOG_LEVEL=debug

