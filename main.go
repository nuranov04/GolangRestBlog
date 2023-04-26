package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("name")
	w.Write([]byte(fmt.Sprintf("Hello %s", name)))
}

func main() {
}
