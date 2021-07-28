#!/bin/bash

echo "Node 1"
curl --request GET --header "Content-Type: application/json" http://0.0.0.0:13801/kvs/key-count

echo "Node 2"
curl --request GET --header "Content-Type: application/json" http://0.0.0.0:13802/kvs/key-count

echo "Node 3"
curl --request GET --header "Content-Type: application/json" http://0.0.0.0:13803/kvs/key-count