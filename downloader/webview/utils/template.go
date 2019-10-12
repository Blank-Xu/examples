package utils

import (
	"html/template"
	"io/ioutil"
)

func OpenTpl(filename string) (*template.Template, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return template.New("tpl").Parse(string(data))
}
