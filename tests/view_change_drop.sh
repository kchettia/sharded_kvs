#!/bin/bash

curl --request PUT --header "Content-Type: application/json" --data '{"view":"10.10.0.4:13800,10.10.0.5:13800"}' http://127.0.0.1:13801/kvs/view-change