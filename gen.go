package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
)

func genFieldsNew(replaceFieldName string, tpl *Template, fields []FieldDef) string {
	var strReplaced []string
	for _, f := range fields {
		tplType, ok := tpl.Typemap[f.fieldType]
		if !ok {
			tplType, ok = tpl.Typemap["struct"]
			if !ok {
				return ""
			}
		}

		tplReflectValue := reflect.ValueOf(tplType)
		replaceField := tplReflectValue.FieldByName(replaceFieldName)
		strBlocks := replaceField.Interface().([]string)

		strTemp := "\t" + strings.Join(strBlocks, "\n\t")

		strTemp = strings.Replace(strTemp, "<TYPE>", f.fieldType, -1)
		strTemp = strings.Replace(strTemp, "<FIELD>", f.fieldName, -1)
		strTemp = strings.Replace(strTemp, "<KEYTYPE>", f.keyType, -1)
		strTemp = strings.Replace(strTemp, "<VALUETYPE>", f.valueType, -1)
		strReplaced = append(strReplaced, strTemp)
	}

	return strings.Join(strReplaced, "\n")

}
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

func genStructs(templateStruct *TemplateType, tpl *Template, def *Def) string {
	var strStructs []string
	var strStruct string
	for _, structDef := range def.structDefs {
		strStruct = strings.Join(templateStruct.Declaration, "\n")

		strStruct = strings.Replace(strStruct, "<TYPE>", structDef.structName, -1)
		strStruct = strings.Replace(strStruct, "<FIELDS>", genFieldsNew("Field", tpl, structDef.fields), -1)
		strStruct = strings.Replace(strStruct, "<MEMBERSERIALIZE>", genFieldsNew("MemberSerialize", tpl, structDef.fields), -1)
		strStruct = strings.Replace(strStruct, "<MEMBERDESERIALIZE>", genFieldsNew("MemberDeserialize", tpl, structDef.fields), -1)
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
		strParam = strings.Replace(strBody, "<FIELD>", field.fieldName, -1)
		strParam = strings.Replace(strParam, "<TYPE>", field.fieldType, -1)
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
		strParam = strings.Replace(strBody, "<FIELD>", field.fieldName, -1)
		strParam = strings.Replace(strParam, "<TYPE>", field.fieldType, -1)
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
