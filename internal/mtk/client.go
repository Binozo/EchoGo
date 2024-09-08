package mtk

import (
	"os"
	"os/exec"
	"path"
)

const mtkClientLocation = "mtkclient"

func Boot() error {
	cmd := exec.Command(path.Join(mtkClientLocation, ".venv", "bin", "python3"), path.Join(mtkClientLocation, "mtk.py"), "plstage", "--preloader=preloader_no_hdr.bin")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	//output, err := cmd.CombinedOutput()
	//fmt.Println(string(output))
	return cmd.Run()
}
