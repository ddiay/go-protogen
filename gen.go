package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

func genFields(templateStruct *TemplateType, template *Template, structDef *StructDef) string {
	var strFields []string
	for _, fieldDef := range structDef.fields {
		fieldT, ok := template.Typemap[fieldDef.fieldType]
		var strField string
		if !ok {
			strField = "\t" + strings.Join(templateStruct.Field, "\n\t")
		} else {
			strField = "\t" + strings.Join(fieldT.Field, "\n\t")
		}
		fmt.Println(fieldDef.fieldType, ok, strField)
		strField = strings.Replace(strField, "<TYPE>", fieldDef.fieldType, -1)
		strField = strings.Replace(strField, "<FIELD>", fieldDef.fieldName, -1)
		strField = strings.Replace(strField, "<KEYTYPE>", fieldDef.keyType, -1)
		strField = strings.Replace(strField, "<VALUETYPE>", fieldDef.valueType, -1)
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

func genParams(tpl *Template, fields []FieldDef) string {
	var strParams []string
	var strParam string

	for _, field := range fields {
		tplField, ok := tpl.Typemap[field.fieldType]
		if !ok {
			tplField, ok = tpl.Typemap["struct"]
			if !ok {
				break
			}
		}

		strBody := strings.Join(tplField.Field, "")
		strParam = strings.Replace(strBody, "<FIELDNAME>", field.fieldName, -1)
		strParam = strings.Replace(strParam, "<FIELDTYPE>", field.fieldType, -1)
		strParam = strings.Replace(strParam, "<KEYTYPE>", field.keyType, -1)
		strParam = strings.Replace(strParam, "<VALUETYPE>", field.valueType, -1)

		strParams = append(strParams, strParam)
	}

	return strings.Join(strParams, ", ")
}

func genVars(tpl *Template, fields []FieldDef) string {
	var strParams []string
	var strParam string

	for _, field := range fields {
		tplField, ok := tpl.Typemap[field.fieldType]
		if !ok {
			tplField, ok = tpl.Typemap["struct"]
			if !ok {
				break
			}
		}

		strBody := strings.Join(tplField.Param, "")
		strParam = strings.Replace(strBody, "<FIELDNAME>", field.fieldName, -1)
		strParam = strings.Replace(strParam, "<FIELDTYPE>", field.fieldType, -1)
		strParam = strings.Replace(strParam, "<KEYTYPE>", field.keyType, -1)
		strParam = strings.Replace(strParam, "<VALUETYPE>", field.valueType, -1)
		strParams = append(strParams, field.fieldName)
	}

	return strings.Join(strParams, ", ")
}

func genMethods(tpl *Template, tplMsg *TemplateMsg, msgDef *MsgDef) string {
	var strResult string
	strBody := strings.Join(tplMsg.Body, "\n")

	for _, method := range msgDef.methods {
		strMethod := strBody
		strMethod = strings.Replace(strMethod, "<METHOD>", method.methodName, -1)
		strMethod = strings.Replace(strMethod, "<FROMTYPE>", msgDef.from, -1)
		strMethod = strings.Replace(strMethod, "<TOTYPE>", msgDef.to, -1)
		strMethod = strings.Replace(strMethod, "<DESERIALIZE>", "", -1)
		strMethod = strings.Replace(strMethod, "<SERIALIZE>", "", -1)
		strMethod = strings.Replace(strMethod, "<PARAMS>", genParams(tpl, method.params), -1)
		strMethod = strings.Replace(strMethod, "<VARS>", genVars(tpl, method.params), -1)
		strResult += strMethod + "\n"
	}

	return strResult
}

func genMsg(tpl *Template, msgDef *MsgDef) string {
	var strMsg string

	tplFrom, ok := tpl.Msgmap["from"]
	if !ok {
		return ""
	}

	tplTo, ok := tpl.Msgmap["to"]
	if !ok {
		return ""
	}

	strMsg += genMethods(tpl, &tplFrom, msgDef)
	strMsg += genMethods(tpl, &tplTo, msgDef)

	return strMsg
}

func genProxies(template *Template, def *Def) string {
	var strMsgs []string
	var strMsg string

	for _, msgDef := range def.msgDefs {
		strMsg = genMsg(template, &msgDef)
		strMsgs = append(strMsgs, strMsg)
	}

	return strings.Join(strMsgs, "\n")
}

func genCallbacks(template *Template, def *Def) string {
	var strCallbacks []string

	return strings.Join(strCallbacks, "\n")
}

func genRpc(template *Template, def *Def) (string, error) {
	var content string

	templateStruct, ok := template.Typemap["struct"]
	if !ok {
		return content, errors.New("struct not found in template")
	}

	content += strings.Join(template.Head, "\n") + "\n"
	content += genStructs(&templateStruct, template, def) + "\n"
	content += genProxies(template, def) + "\n"
	content += genCallbacks(template, def) + "\n"
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
