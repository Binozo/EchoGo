package server

import (
	"fmt"
	"github.com/Binozo/EchoGo/v2/pkg/constants"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{} // use default options

func Serve() error {
	r := mux.NewRouter()

	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("EchoGo"))
	})

	r.HandleFunc("/buttons", buttonHandler)
	r.HandleFunc("/microphone", micHandler)
	r.HandleFunc("/speaker", speakerHandler)
	r.HandleFunc("/led", ledHandler)

	http.Handle("/", r)
	return http.ListenAndServe(fmt.Sprintf(":%d", constants.Port), r)
}
