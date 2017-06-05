package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
)

type Config struct {
	templateFilePath string
	defFilePath      string
	rpcFilePath      string
}

type Test struct {
	num int
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
	f := 3121.3
	s := "Hello!你好！World!"
	data := make([]byte, 64)

	binary.LittleEndian.PutUint64(data, math.Float64bits(f))
	copy(data[8:], []byte(s))

	fmt.Println(len(s), []byte(s))
	fmt.Println(data)

	f1 := math.Float64frombits(binary.LittleEndian.Uint64(data))
	s1 := string(data[8 : len(s)+8])

	fmt.Println(f1, s1, []byte(s1))
}
