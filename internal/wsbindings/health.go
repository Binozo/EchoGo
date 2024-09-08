package wsbindings

import (
	"fmt"
	"github.com/Binozo/EchoGo/v2/pkg/constants"
	"net/http"
)

func CheckHealth() (isOnline bool) {
	res, err := http.Get(fmt.Sprintf("http://localhost:%d/", constants.Port))
	if err != nil {
		return false
	}
	defer res.Body.Close()
	return true
}
