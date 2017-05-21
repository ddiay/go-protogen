package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

type FieldDef struct {
	fieldType string
	fieldName string
}

type StructDef struct {
	structName string
	fields     []FieldDef
}

type FuncDef struct {
	funcName string
	params   []FieldDef
}

var structDefs []StructDef
var funcDefs []FuncDef

func ParseDef(text string) {
	regStruct := regexp.MustCompile(`struct[\s]+([\w]+)[\s]+\{((?s:.*?))\}`)
	regField := regexp.MustCompile(`[\s]+([\w]+)[\s]+([\w]+)`)

	matchStructs := regStruct.FindAllStringSubmatch(text, -1)

	for _, matchStruct := range matchStructs {
		matchFields := regField.FindAllStringSubmatch(matchStruct[2], -1)
		structDef := StructDef{
			structName: matchStruct[1],
		}
		for _, matchField := range matchFields {
			fieldDef := FieldDef{
				fieldType: matchField[1],
				fieldName: matchField[2],
			}
			structDef.fields = append(structDef.fields, fieldDef)
		}
		structDefs = append(structDefs, structDef)
	}
	fmt.Printf("%v\n", structDefs)

	regFunc := regexp.MustCompile(`func[\s]+([\w]+)[\s]*\((.*?)\)`)
	regParam := regexp.MustCompile(`[\s]*([\w]+)[\s]+([\w]+)[,]*`)

	matchFuncs := regFunc.FindAllStringSubmatch(text, -1)
	for _, matchFunc := range matchFuncs {
		matchParams := regParam.FindAllStringSubmatch(matchFunc[2], -1)
		funcDef := FuncDef{
			funcName: matchFunc[1],
		}
		for _, matchParam := range matchParams {
			fieldDef := FieldDef{
				fieldType: matchParam[1],
				fieldName: matchParam[2],
			}
			funcDef.params = append(funcDef.params, fieldDef)
		}
		funcDefs = append(funcDefs, funcDef)
	}
	fmt.Printf("%v\n", funcDefs)
}

func 

func main() {
	content, err := ioutil.ReadFile("protocol/test.def")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ParseDef(string(content))
}
