package test

import "github.com/ddiay/go-gameserver/acctserver"

type foo struct {
	num int
	
}

func Marshalfoo(self foo, data []byte) {
}

func Unmarshalfoo(self foo, data []byte) {
}

type foo1 struct {
	
	name string
}

func Marshalfoo1(self foo1, data []byte) {
}

func Unmarshalfoo1(self foo1, data []byte) {
}

type ServerInfo struct {
	id int
	name string
	status int
}

func MarshalServerInfo(self ServerInfo, data []byte) {
}

func UnmarshalServerInfo(self ServerInfo, data []byte) {
}
