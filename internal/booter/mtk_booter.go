package booter

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

const mtkClientLocation = "mtkclient"

type MtkBooter struct {
	preloaderPath string
}

func NewMtkBooter(preloaderPath string) *MtkBooter {
	return &MtkBooter{
		preloaderPath: preloaderPath,
	}
}

func (m *MtkBooter) Boot() error {
	cmd := exec.Command(path.Join(mtkClientLocation, ".venv", "bin", "python3"), path.Join(mtkClientLocation, "mtk.py"), "plstage", fmt.Sprintf("--preloader=%s", m.preloaderPath))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
