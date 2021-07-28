module github.com/kchettia/sharded_kvs/routes

go 1.16

replace github.com/kchettia/sharded_kvs/key_distributor => ./key_distributor

require (
	github.com/gorilla/mux v1.8.0
	github.com/kchettia/sharded_kvs/key_distributor v0.0.0
)
