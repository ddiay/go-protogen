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

type MethodDef struct {
	methodName string
	params     []FieldDef
}

type MsgDef struct {
	from    string
	to      string
	methods []MethodDef
}

type Def struct {
	structDefs []StructDef
	msgDefs    []MsgDef
}

func LoadStructs(text string) []StructDef {
	var structDefs []StructDef
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
	return structDefs
}

func LoadFuncs(text string) []FuncDef {
	var funcDefs []FuncDef

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
	return funcDefs
}

func loadMethod(text string) []MethodDef {
	var methodDefs []MethodDef

	return methodDefs
}

func loadMsgs(text string) []MsgDef {
	var msgDefs []MsgDef

	reg := regexp.MustCompile(`msg[\s]+([\w]+)[\s]+->[\s]+([\w]+)\{((?s:.*?))\}`)
	// regField := regexp.MustCompile(`[\s]+([\w]+)[\s]+([\w]+)`)

	matches := reg.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		methodDefs := loadMethod(match[2])
		msgDef := MsgDef{
			from:    match[1],
			to:      match[2],
			methods: methodDefs,
		}

		msgDefs = append(msgDefs, msgDef)
	}

	fmt.Printf("msg:\n%q\n", matches)
	return msgDefs
}

func LoadDef(path string) (*Def, error) {
	def := &Def{}
	defData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(defData)
	def.structDefs = LoadStructs(content)
	def.msgDefs = loadMsgs(content)
	return def, nil
}
