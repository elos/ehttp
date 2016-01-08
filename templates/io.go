package templates

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

// EnsureDirectoryExists creates a directory with permission
// bits 0755 if that directory doesn't already exist.
// That is: drwxr-xr-x permissions
func EnsureDirectoryExists(path string) (err error) {
	err = os.Mkdir(path, 0755)
	if err != nil {
		// If it already exists, we are fine.
		if os.IsExist(err) {
			err = nil
		}
	}
	return
}

// ExecuteAndWriteGoFile generates a *.go file based on the template 'n', using
// the engine 'e' and writes it to the file 'file' using the context data 'd'
func ExecuteAndWriteGoFile(e *Engine, n Name, file string, d interface{}) error {
	if err := ExecuteAndWrite(e, n, file, d); err != nil {
		return err
	}

	if err := exec.Command("goimports", "-w=true", file).Run(); err != nil {
		return fmt.Errorf("error running goimports on file: %s", err)
	}

	return nil
}

func ExecuteAndWrite(e *Engine, n Name, file string, d interface{}) error {
	var buf bytes.Buffer

	if err := e.Execute(&buf, n, d); err != nil {
		return fmt.Errorf("error executing template: %s", err)
	}

	if err := ioutil.WriteFile(file, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("error writing file: %s", err)
	}

	log.Printf("wrote %s", file)

	return nil
}
