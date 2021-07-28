#! /usr/bin/bash
sudo docker stop $(sudo docker ps -a -q)
sudo docker rm $(sudo docker ps -a -q)
sudo docker network create --subnet=10.10.0.0/16 kv_subnet || true
sudo docker build -t kvs:3.0 . 
sudo docker run -d -p 13801:13800 --net=kv_subnet --ip=10.10.0.4 --name="node1" -e ADDRESS="10.10.0.4:13800" -e VIEW="10.10.0.4:13800,10.10.0.5:13800" kvs:3.0
sudo docker run -d -p 13803:13800 --net=kv_subnet --ip=10.10.0.6 --name="node3" -e ADDRESS="10.10.0.6:13800" -e VIEW="10.10.0.6:13800" kvs:3.0
sudo docker run -p 13802:13800 --net=kv_subnet --ip=10.10.0.5 --name="node2" -e ADDRESS="10.10.0.5:13800" -e VIEW="10.10.0.4:13800,10.10.0.5:13800" kvs:3.0

