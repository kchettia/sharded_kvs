package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	kd "github.com/kchettia/sharded_kvs/key_distributor"
	"net/http"
	//"sort"
	"strings"
)

type KeyVal struct {
	Value string `json:"value"`
}
type view struct {
	Value string
}

var server_addr string

//var KVS = make(map[string]string)
var shards = make(map[string]map[string]string)

func shouldForward(key string) (bool, string) {
	closest_serv := kd.FindClosestNode(key)
	if closest_serv != server_addr {
		return true, closest_serv
	} else {
		return false, closest_serv
	}
}

func forwardRequest(addr, reqType string, reqBody []byte) (map[string]interface{}, int) {
	//fmt.Println("Forward", addr, reqType, reqBody)
	req, _ := http.NewRequest(reqType, addr, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		fmt.Println("Error with Forwarded Request")
		panic(err)
	}
	defer res.Body.Close()
	data := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		fmt.Println("Failed to Decode")
		panic(err)
	}
	return data, res.StatusCode

}
func responseFormatter(w http.ResponseWriter, response interface{}, httpStatus int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(response)
}

func put_key(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Received")
	vars := mux.Vars(r)
	key := vars["key"]
	//var data KeyVal
	if len(key) > 50 {
		fmt.Println("err1")
		response := PutErrorResp{
			Error:   "Key is too long",
			Message: "Error in PUT"}
		responseFormatter(w, response, 400)
		return
	}
	//err := json.NewDecoder(r.Body).Decode(&data)
	data := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&data)
	_, ok := data["value"]
	if err != nil || !ok {
		fmt.Println("err2", err)
		response := PutErrorResp{
			Error:   "Value is missing",
			Message: "Error in PUT"}
		responseFormatter(w, response, 400)
		return
	}
	value := fmt.Sprintf("%v", data["value"])
	if forwardFlag, closest_server := shouldForward(key); forwardFlag {
		//fmt.Println("forward")
		body := fmt.Sprintf("{\"value\":\"%s\"}", value)
		var jsonData = []byte(body)
		closest_server = fmt.Sprintf("http://%s/kvs/keys/%s", closest_server, key)
		response, status_code := forwardRequest(closest_server, "PUT", jsonData)
		responseFormatter(w, response, status_code)

	} else {
		//_, ok := KVS[key]
		//KVS[key] = data.Value

		ok := kd.AddKey(key, value)
		if ok {
			status_code := 200
			response := PutResp{
				Message:  "Updated successfully",
				Replaced: true,
				Address:  server_addr}
			responseFormatter(w, response, status_code)
		} else {
			status_code := 201
			response := PutResp{
				Message:  "Added successfully",
				Replaced: false,
				Address:  server_addr}
			responseFormatter(w, response, status_code)

		}

	}
}

func get_key(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	if forwardFlag, closest_server := shouldForward(key); forwardFlag {
		closest_server = "http://" + closest_server + "/kvs/keys/" + key
		response, status_code := forwardRequest(closest_server, "GET", nil)
		responseFormatter(w, response, status_code)
	} else {
		value, ok := kd.GetKey(key)
		if !ok {
			status_code := 404
			response := GetErrorResp{
				Message:   "Error in GET",
				DoesExist: false,
				Error:     "Key does not exist"}
			responseFormatter(w, response, status_code)
		} else {
			status_code := 200
			response := GetResp{
				Message:   "Retrieved successfully",
				DoesExist: true,
				Value:     value}
			responseFormatter(w, response, status_code)
		}
	}

}

func del_key(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	if forwardFlag, closest_server := shouldForward(key); forwardFlag {
		fmt.Println("forward")
		//closest_server = "http://" + closest_server + "/kvs/keys/" + key
		closest_server = fmt.Sprintf("http://%v/kvs/keys/%v", closest_server, key)
		response, status_code := forwardRequest(closest_server, "DELETE", nil)
		responseFormatter(w, response, status_code)
	} else {
		ok := kd.DelKey(key)
		if !ok {
			status_code := 404
			response := DelErrorResp{
				Message:   "Error in DELETE",
				DoesExist: false,
				Error:     "Key does not exist"}
			responseFormatter(w, response, status_code)
		} else {
			status_code := 200
			response := DelResp{
				Message:   "Deleted successfully",
				DoesExist: true,
				Address:   server_addr}
			responseFormatter(w, response, status_code)
		}
	}

}

func leader_view_change(w http.ResponseWriter, r *http.Request) {
	body := make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&body)
	view, ok := body["view"]
	if err != nil || !ok {
		response := ViewChangeErrorResp{
			Error:   "Value is missing",
			Message: "Error in View-Change"}
		responseFormatter(w, response, 400)
		return
	}
	newView := strings.Split(view, ",")
	mergedView := mergeViews(newView, kd.GetCurrentView())
	fmt.Println(mergedView)
	data := fmt.Sprintf("{\"view\":\"%v\"}", view)
	fmt.Println("DATA", data)
	for key, _ := range mergedView {
		if key == server_addr {
			continue
		}
		url := fmt.Sprintf("http://%v/kvs/follower-view-change", key)
		var jsonData = []byte(data)
		response, status_code := forwardRequest(url, "PUT", jsonData)
		fmt.Println(response, status_code)

	}
	shards, count := kd.DistributeKeys()
	for node_ip, shard := range shards {
		fmt.Fprintf(w, "%s : %d\n", node_ip, len(shard))
	}
	fmt.Fprintf(w, "count : %d\n", count)
	//kd.TempView(body["view"])

}
func follower_view_change(w http.ResponseWriter, r *http.Request) {
	fmt.Println("follower")
	body := make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&body)
	view, ok := body["view"]
	if err != nil || !ok {
		response := ViewChangeErrorResp{
			Error:   "Value is missing",
			Message: "Error in View-Change"}
		responseFormatter(w, response, 400)
		return
	}
	kd.Setview(view)
	shards, _ = kd.DistributeKeys()
	/* for node_ip, shard := range shards {
		fmt.Sprintf(w, "%s : %d\n", node_ip, len(shard))
	} */
	response := ViewChangeResp{
		Message: "View-Change Initiated"}
	responseFormatter(w, response, 200)
}

func mergeViews(a []string, b []string) map[string]bool {
	//	mergedView := make([]string)
	/*	for i := range a {
		mergedView = append(mergedView, b[i])
	}*/
	mergedView := make(map[string]bool)
	view := a
	for i := range b {
		view = append(view, b[i])
	}
	for j := range view {
		_, ok := mergedView[view[j]]
		if !ok {
			mergedView[view[j]] = true
		}
	}

	return mergedView
}

func key_count(w http.ResponseWriter, r *http.Request) {
	//key_count := kd.GetKeyCount()
	response := KeyCountResp{
		Message:  "Key count retrieved successfully",
		KeyCount: kd.GetKeyCount()}
	responseFormatter(w, response, 201)

}

func Request_handler(router *mux.Router, addr string, view string) {
	kd.Setview(view)
	server_addr = addr
	router.HandleFunc("/kvs/keys/{key}", get_key).Methods("GET")
	router.HandleFunc("/kvs/keys/{key}", put_key).Methods("PUT")
	router.HandleFunc("/kvs/keys/{key}", del_key).Methods("DELETE")
	router.HandleFunc("/kvs/key-count", key_count).Methods("GET")
	router.HandleFunc("/kvs/view-change", leader_view_change).Methods("PUT")
	router.HandleFunc("/kvs/follower-view-change", follower_view_change).Methods("PUT")
}
