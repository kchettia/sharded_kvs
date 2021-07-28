module github.com/kchettia/sharded_kvs

go 1.16

replace github.com/kchettia/sharded_kvs/routes => ./routes

replace github.com/kchettia/sharded_kvs/key_distributor => ./key_distributor

require (
	github.com/gorilla/mux v1.8.0
	github.com/kchettia/sharded_kvs/key_distributor v0.0.0
	github.com/kchettia/sharded_kvs/routes v0.0.0
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	lukechampine.com/uint128 v1.1.1 // indirect
)
