package bridge

import "os"

type Bridge interface {
	WaitForGettingOnline() error
	IsConnected() (bool, error)
	Deploy(file *os.File) error
	Run() error
}
