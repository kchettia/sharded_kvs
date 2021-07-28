package key_distributor

import (
	"fmt"
	"github.com/spaolacci/murmur3"
	"hash/fnv"
	"sort"
	"strings"
)

type Shard map[string]string

var count int
var shard = make(map[string]string)
var numVirtualNodes = 100
var view []string

//var temp_view []string
var nodeMapkeys []uint64
var nodeMap map[uint64]string

func AddKey(key string, value string) bool {
	_, ok := shard[key]
	shard[key] = value
	if ok {
		return true
	} else {
		return false
	}
}
func DelKey(key string) bool {
	_, ok := shard[key]
	if ok {
		delete(shard, key)
		return true
	} else {
		return false
	}
}
func GetKey(key string) (string, bool) {
	value, ok := shard[key]
	if ok {
		return value, true
	} else {
		return "", false
	}
}
func GetKeyCount() int {
	return len(shard)
}
func Setview(current_view string) {
	fmt.Println("numVN:", numVirtualNodes)
	view = strings.Split(current_view, ",")
	sort.Strings(view)
	//fmt.Println(view, len(view))
	createNodeMap()
}
func GetCurrentView() []string {
	return view
}
func DistributeKeys() (map[string]map[string]string, int) {
	shards := make(map[string]map[string]string)
	for _, node_ip := range view {
		//fmt.Println(node_ip)
		shards[node_ip] = make(map[string]string)
	}
	for key, value := range shard {
		ip := FindClosestNode(key)
		//fmt.Println(ip, key)
		shards[ip][key] = value
	}

	return shards, count
	/*for _, node_ip := range view {
		fmt.Println(node_ip, " :", len(shards[node_ip]))
	}*/
}

func createNodeMap() {
	nodeMap = make(map[uint64]string)
	nodeMapkeys = make([]uint64, numVirtualNodes*len(view))
	for i, node_ip := range view {
		for j := 0; j < numVirtualNodes; j++ {
			vnid := fmt.Sprintf("%s%d%d", node_ip, i, j) //Virtual node id
			hashed_vnid := hash(vnid)
			//fmt.Println(hashed_vnid)
			_, ok := nodeMap[hashed_vnid]
			if ok {
				fmt.Println(vnid, " Collision\n")
			} else {
				nodeMap[hashed_vnid] = node_ip
				nodeMapkeys[i*numVirtualNodes+j] = hashed_vnid
			}
		}
	}
	sort.Slice(nodeMapkeys, func(i, j int) bool { return nodeMapkeys[i] < nodeMapkeys[j] })
	//fmt.Println("NodeMapKeys", nodeMapkeys)
}

/* func hash(s string) uint64 {
	h := murmur3.New64()
	h.Write([]byte(s))
	return h.Sum64()
} */

func hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	hash_1 := fmt.Sprintf("%v%s", h.Sum64(), s)
	h2 := murmur3.New64()
	h2.Write([]byte(hash_1))
	return h2.Sum64()
}

func FindClosestNode(key string) string {
	//hashed_key := hash(key) % nodeMapkeys[len(nodeMapkeys)-1]
	hashed_key := hash(key)
	index := search(nodeMapkeys, hashed_key, 1, len(nodeMapkeys)-1)
	//fmt.Println("Key: ", key, "hash: ", hashed_key, "Node: ", nodeMapkeys[index])
	return nodeMap[nodeMapkeys[index]]
}

func search(arr []uint64, target uint64, low int, high int) int {
	if target <= arr[0] || target > arr[len(arr)-1] {
		count = count + 1
		return 0
	} else {
		return binarySearch(arr, target, 1, len(arr)-1)
	}
}

func binarySearch(arr []uint64, target uint64, low int, high int) int {

	if high == low {
		return high
	}
	mid := (low + high) / 2
	if arr[mid] > target {
		return binarySearch(arr, target, low, mid)
	} else if arr[mid] < target {
		return binarySearch(arr, target, mid+1, high)
	} else {
		return mid
	}
}
