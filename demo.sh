#!/usr/bin/env bash

set -e
PORT=6379

echo
echo ">>>PING / GET / SET<<<"
redis-cli -p $PORT <<EOF
PING
SET user test-user
GET user
TYPE user
EOF

echo
echo ">>>LISTS<<<"
redis-cli -p $PORT <<EOF
LPUSH fruits apple mango banana
LRANGE fruits 0 -1
LPOP fruits
LRANGE fruits 0 -1
EOF

echo
echo ">>>TTL / EXPIRATION<<<"
redis-cli -p $PORT <<EOF
SET temp value EX 1
GET temp
EOF

sleep 2

redis-cli -p $PORT <<EOF
GET temp
TYPE temp
EOF

echo
echo ">>>STREAMS<<<"
redis-cli -p $PORT <<EOF
XADD posts * title godis-is-good
XADD posts * title not-joking
TYPE posts
XRANGE STREAMS posts - +
XREAD STREAMS posts 0-0
EOF

echo
echo ">>>CONFIG<<<"
redis-cli -p $PORT <<EOF
CONFIG GET dir
CONFIG SET dir /tmp
CONFIG GET dir
EOF
