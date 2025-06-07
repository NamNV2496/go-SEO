package utils

import (
	"bytes"
	"context"
	"html/template"
)

func BuildByTemplate(ctx context.Context, name, tpl string, request map[string]string) (string, error) {
	t, err := template.New(name).Parse(tpl)
	if err != nil {
		return "", err
	}
	result := new(bytes.Buffer)
	err = t.Execute(result, request)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}
