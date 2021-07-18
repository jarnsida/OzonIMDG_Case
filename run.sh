#!/bin/bash

echo -n "Введите порт подключения к серверу: "
read port

echo -n "Введите максимальное количество записей в базе: "
read maxEl

echo -n "Введите таймаут отклчения пользователей при отключении сервера: "
read timeOut

source ./cmd/env.sh $port $maxEl $timeOut
go run -race ./cmd/*.go