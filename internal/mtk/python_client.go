package mtk

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

type PythonClient struct {
}

const mtkClientLocation = "mtkclient"

func (p *PythonClient) BootDevice(preloaderPath string) error {
	cmd := exec.Command(path.Join(mtkClientLocation, ".venv", "bin", "python3"), path.Join(mtkClientLocation, "mtk.py"), "plstage", fmt.Sprintf("--preloader=%s", preloaderPath))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
