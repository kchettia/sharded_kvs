#!/bin/bash

docker stop node3 && docker rm node3

docker run -d -p 13803:13800 --net=kv_subnet --ip=10.10.0.6 --name="node3" \
      -e ADDRESS="10.10.0.6:13800" -e VIEW="10.10.0.4:13800,10.10.0.5:13800,10.10.0.6:13800" \
      kvs:3.0

sleep 2
curl --request PUT --header "Content-Type: application/json" --data '{"view":"10.10.0.4:13800,10.10.0.5:13800,10.10.0.6:13800"}' http://127.0.0.1:13801/kvs/view-change
