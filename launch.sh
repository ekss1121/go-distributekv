#!/bin/bash
set -e
# kill all lingering process
trap 'killall distributeKV' SIGINT

cd $(dirname $0)
killall distributeKV || true
sleep 0.1

go install -v
../../bin/distributeKV -launch-mode="DB" -db-location=$PWD/atlanta.db -partition-config=$PWD/static_partition.toml -http-addr=127.0.0.1:8080 -partition=Atlanta &
../../bin/distributeKV -launch-mode="DB" -db-location=$PWD/texas.db -partition-config=$PWD/static_partition.toml -http-addr=127.0.0.1:8081 -partition=Texas &
../../bin/distributeKV -launch-mode="DB" -db-location=$PWD/sanfrancisco.db -partition-config=$PWD/static_partition.toml -http-addr=127.0.0.1:8082 -partition=SanFrancisco &
../../bin/distributeKV -launch-mode="Proxy" -partition-config=$PWD/static_partition.toml &

wait