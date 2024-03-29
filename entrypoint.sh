#!/bin/sh
#Entry point docker container

cd /app

if [ -z "${CONFIG}" ];then
    CONFIG="bot.toml"
fi

if [ ! -z "${NUMBER}" ];then
    ./registar -a $COUNT -c $CONFIG -n $NUMBER
    exit 0
fi

if [ ! -z "${SERVICE}" ];then
    ./registar -a $COUNT -c $CONFIG -s $SERVICE
    exit 0
fi

exit 1

