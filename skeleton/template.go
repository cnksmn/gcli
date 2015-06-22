package skeleton

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/tcnksm/gcli/helper"
)

type Template struct {
	// Path is the path to this template.
	Path string

	// OutputPathTmpl is the template for outputPath.
	OutputPathTmpl string
}

// Exec evaluate this template and write it to provided file.
// At First, it reads template content.
// Then, it generates output file path from output path template and its data.
// Then, it creates directory if not exist from output path.
// Then, it opens output file.
// Finally, it evaluates template contents and generate it to output file.
// If output file is gocode, run go fmt.
//
// It returns an error if any.
func (t *Template) Exec(data interface{}) error {
	// Asset function is generated by go-bindata
	contents, err := Asset(t.Path)
	if err != nil {
		return err
	}

	outputPath, err := processPathTmpl(t.OutputPathTmpl, data)
	if err != nil {
		return err
	}

	// Create directory if necessary
	dir, _ := filepath.Split(outputPath)
	if dir != "" {
		if err := mkdir(dir); err != nil {
			return err
		}
	}

	wr, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	if err := execute(string(contents), wr, data); err != nil {
		return err
	}

	if strings.HasSuffix(outputPath, ".go") {
		helper.GoFmt(outputPath, nil)
	}

	return nil
}

// processPathTmpl evaluates output path template string
// and generate real absolute path. Any errors that occur are returned.
func processPathTmpl(pathTmpl string, data interface{}) (string, error) {
	var outputPathBuf bytes.Buffer
	if err := execute(pathTmpl, &outputPathBuf, data); err != nil {
		return "", err
	}

	return outputPathBuf.String(), nil
}

// mkdir makes the named directory.
func mkdir(dir string) error {
	if _, err := os.Stat(dir); err == nil {
		return nil
	}

	return os.MkdirAll(dir, 0777)
}

// execute evaluates template content with data
// and write them to writer.
func execute(content string, wr io.Writer, data interface{}) error {
	name := ""
	funcs := funcMap()

	tmpl, err := template.New(name).Funcs(funcs).Parse(content)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(wr, data); err != nil {
		return err
	}

	return nil
}

func funcMap() template.FuncMap {
	return template.FuncMap{
		"date":  dateFunc(),
		"title": strings.Title,
	}
}
