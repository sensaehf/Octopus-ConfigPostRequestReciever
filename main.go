package main

import (
	"net/http"
)

func recieveInfo(w http.ResponseWriter, req *http.Request) {

}

func main() {

	http.HandleFunc("/", recieveInfo)

	http.ListenAndServe(":8090", nil)
}
