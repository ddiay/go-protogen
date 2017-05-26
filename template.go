package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type TemplateType struct {
	Variable    []string
	Field       []string
	Param       []string
	Serialize   []string
	Deserialize []string
}

type Template struct {
	Language string
	Head     []string
	Tail     []string
	Typemap  map[string]TemplateType
}

func LoadTemplate(path string) (*Template, error) {
	template := &Template{}

	jsonData, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	err = json.Unmarshal(jsonData, template)
	if err != nil {
		return nil, err
	}
	return template, nil
}
