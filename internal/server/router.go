package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
)

const Port = 8092

var upgrader = websocket.Upgrader{} // use default options

func Serve() error {
	r := mux.NewRouter()

	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("EchoGo"))
	})

	r.HandleFunc("/buttons", buttonHandler)

	http.Handle("/", r)
	return http.ListenAndServe(fmt.Sprintf(":%d", Port), r)
}
