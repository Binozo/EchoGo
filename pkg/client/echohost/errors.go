package echohost

import "errors"

var ErrAlexaNotConnected = errors.New("not connected to Alexa")
var ErrAlexaBootTimeout = errors.New("timeout waiting for Alexa boot")
