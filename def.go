package main

import (
	"io/ioutil"
	"regexp"
)

const (
	BASETYPE = 1
	LIST     = 2
	MAP      = 3
	STRUCT   = 4
)

type FieldDef struct {
	fieldType string
	fieldName string
	keyType   string
	valueType string
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

func loadFields(text string) []FieldDef {
	var fieldDefs []FieldDef

	// regBaseType := regexp.MustCompile(`[\s]+([\w]+)[\s]+([\w]+)`)
	// regList := regexp.MustCompile(`[\s]+(list\[([\w]+)\])[\s]+([\w]+)`)
	// regMap := regexp.MustCompile(`[\s]+(map\[([\w]+)\]([\w]+))[\s]+([\w]+)`)
	reg := regexp.MustCompile(`[\s]*(((map)\[([\w]+)\]([\w]+))|((list)\[([\w]+)\])|([\w]+))[\s]+([\w]+)[,]*`)

	matches := reg.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		fieldDef := FieldDef{
			fieldType: match[1],
			fieldName: match[10],
		}

		if match[3] == "map" {
			fieldDef.fieldType = match[3]
			fieldDef.keyType = match[4]
			fieldDef.valueType = match[5]
		} else if match[7] == "list" {
			fieldDef.fieldType = match[7]
			fieldDef.valueType = match[8]
		} else {
			fieldDef.fieldType = match[9]
		}
		fieldDefs = append(fieldDefs, fieldDef)
	}

	return fieldDefs
}

func loadStructs(text string) []StructDef {
	var structDefs []StructDef
	regStruct := regexp.MustCompile(`struct[\s]+([\w]+)[\s]+\{((?s:.*?))\}`)

	matchStructs := regStruct.FindAllStringSubmatch(text, -1)

	for _, matchStruct := range matchStructs {
		structDef := StructDef{
			structName: matchStruct[1],
			fields:     loadFields(matchStruct[2]),
		}
		structDefs = append(structDefs, structDef)
	}
	return structDefs
}

func loadFuncs(text string) []FuncDef {
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

func loadParams(text string) []FieldDef {
	var fieldDefs []FieldDef

	regBaseType := regexp.MustCompile(`[\s]*([\w]+)[\s]+([\w]+)[,]*`)
	regList := regexp.MustCompile(`[\s]*(list\[([\w]+)\])[\s]+([\w]+)[,]*`)
	regMap := regexp.MustCompile(`[\s]*(map\[([\w]+)\]([\w]+))[\s]+([\w]+)[,]*`)

	matches := regBaseType.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		fieldDef := FieldDef{
			fieldType: match[1],
			fieldName: match[2],
		}

		fieldDefs = append(fieldDefs, fieldDef)
	}

	matches = regList.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		fieldDef := FieldDef{
			fieldType: match[1],
			fieldName: match[3],
			valueType: match[2],
		}
		fieldDefs = append(fieldDefs, fieldDef)
	}

	matches = regMap.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		fieldDef := FieldDef{
			fieldType: match[1],
			fieldName: match[4],
			keyType:   match[2],
			valueType: match[3],
		}
		fieldDefs = append(fieldDefs, fieldDef)
	}

	return fieldDefs

}

func loadMethod(text string) []MethodDef {
	var methodDefs []MethodDef

	reg := regexp.MustCompile(`[\s]*([\w]+)[\s]*\((.*?)\)`)

	matches := reg.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		methodDef := MethodDef{
			methodName: match[1],
			params:     loadFields(match[2]),
		}

		methodDefs = append(methodDefs, methodDef)
	}

	return methodDefs
}

func loadMsgs(text string) []MsgDef {
	var msgDefs []MsgDef

	reg := regexp.MustCompile(`msg[\s]+([\w]*)[\s]*(->)[\s]*([\w]*)[\s]*\{((?s:.*?))\}`)
	// regField := regexp.MustCompile(`[\s]+([\w]+)[\s]+([\w]+)`)

	matches := reg.FindAllStringSubmatch(text, -1)
	// fmt.Printf("msg:\n%q\n", matches)
	for _, match := range matches {
		methodDefs := loadMethod(match[4])
		msgDef := MsgDef{
			from:    match[1],
			to:      match[3],
			methods: methodDefs,
		}

		msgDefs = append(msgDefs, msgDef)
	}

	return msgDefs
}

func LoadDef(path string) (*Def, error) {
	def := &Def{}
	defData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(defData)
	def.structDefs = loadStructs(content)
	def.msgDefs = loadMsgs(content)
	return def, nil
}
