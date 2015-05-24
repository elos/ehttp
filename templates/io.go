package templates

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

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
