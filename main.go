package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func GenGo(structDefs []StructDef) {
	var content string

	// var line string
	// for i, def := range structDefs {

	// }

	ioutil.WriteFile("rpc.go", []byte(content), 0644)
}

func GenRPCFiles(text string) {
	structDefs := LoadStructs(text)
	GenGo(structDefs)
}

type Config struct {
	templateFilePath string
	defFilePath      string
	rpcFilePath      string
}

func main() {
	if len(os.Args) < 3 {
		return
	}

	template, err := LoadTemplate(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	def, err := LoadDef(os.Args[2])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("template:\n%v\n", template)
	fmt.Printf("def:\n%v\n", def)

	err = SaveRpc(template, def, os.Args[3])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// GenRPCFiles(string(content))
}
