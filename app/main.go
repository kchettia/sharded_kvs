package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/kchettia/sharded_kvs/routes"
)

func main() {
	fmt.Println("Server Running...")
	router := mux.NewRouter()
	addr := os.Getenv("ADDRESS")
	view := os.Getenv("VIEW")
	routes.Request_handler(router, addr, view)
	fmt.Println(addr)
	http.ListenAndServe("0.0.0.0:13800", router)

}
