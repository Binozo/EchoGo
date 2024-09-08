package wsbindings

import (
	"fmt"
	"github.com/Binozo/EchoGo/v2/pkg/constants"
	"github.com/gorilla/websocket"
)

type baseConnection struct {
	con *websocket.Conn
}

func (b *baseConnection) Close() {
	b.con.Close()
}

type Connector interface {
	Connect() error
	Close()
}

func connect(route string) (*websocket.Conn, error) {
	builtRoute := fmt.Sprintf("ws://localhost:%d%s", constants.Port, route)
	c, _, err := websocket.DefaultDialer.Dial(builtRoute, nil)
	return c, err
}
