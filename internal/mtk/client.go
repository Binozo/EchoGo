package mtk

type Client interface {
	BootDevice(preloaderPath string) error
}

func NewDefaultClient() Client {
	return &PythonClient{}
}
