package main

import (
	"errors"
	"io/ioutil"
	"strings"
)

func genFields(templateStruct *TemplateType, template *Template, structDef *StructDef) string {
	var strFields []string
	for _, fieldDef := range structDef.fields {
		fieldT, ok := template.Typemap[fieldDef.fieldType]
		var strField string
		if !ok {
			strField = "\t" + strings.Join(fieldT.Field, "\n\t")
		} else {
			strField = "\t" + strings.Join(templateStruct.Field, "\n\t")
		}
		strField = strings.Replace(strField, "<FIELD>", fieldDef.fieldName, -1)
		strField = strings.Replace(strField, "<TYPE>", fieldDef.fieldType, -1)
		strFields = append(strFields, strField)
	}
	return strings.Join(strFields, "\n")
}

func genStructs(templateStruct *TemplateType, template *Template, def *Def) string {
	var strStructs []string
	var strStruct string
	for _, structDef := range def.structDefs {
		strStruct = strings.Join(templateStruct.Declaration, "\n")

		strFields := genFields(templateStruct, template, &structDef)

		strStruct = strings.Replace(strStruct, "<TYPE>", structDef.structName, -1)
		strStruct = strings.Replace(strStruct, "<FIELDS>", strFields, -1)
		strStructs = append(strStructs, strStruct)
	}

	return strings.Join(strStructs, "\n")
}

func genProxy() {

}

func genRpc(template *Template, def *Def) (string, error) {
	var content string

	templateStruct, ok := template.Typemap["struct"]
	if !ok {
		return content, errors.New("struct not found in template")
	}

	content += strings.Join(template.Head, "\n") + "\n"

	content += genStructs(&templateStruct, template, def)

	content += strings.Join(template.Tail, "\n")

	return content, nil
}

func SaveRpc(template *Template, def *Def, path string) error {
	content, err := genRpc(template, def)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}
