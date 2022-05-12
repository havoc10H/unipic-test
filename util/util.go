package util

import (
	"bytes"
	"html/template"
	"os"

	"github.com/sysdevguru/unipic/config"
)

// ParseTemplate generates html file
func ParseTemplate(filePath string, data interface{}) error {
	t, err := template.ParseFiles(filePath)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}

	// generate html
	f, err := os.Create(config.Global.Config.IndexPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.Write([]byte(buf.String())); err != nil {
		return err
	}

	return nil
}