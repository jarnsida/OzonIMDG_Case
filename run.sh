#!/bin/bash
if [ -n "$1" ]
then
export port=$1
else
echo -n "Введите порт подключения к серверу: "
read port
fi

if [ -n "$2" ]
then
export maxEl=$2
else
echo -n "Введите максимальное количество записей в базе: "
read maxEl
fi

if [ -n "$3" ]
then
export timeOut=$3
else
echo -n "Введите таймаут отклчения пользователей при отключении сервера: "
read timeOut
fi



source ./cmd/env.sh $port $maxEl $timeOut
go run -race ./cmd/*.go

