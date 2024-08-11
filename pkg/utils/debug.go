package utils

import (
	"encoding/json"
	"fmt"
)

func Debug(data any) {
	//print json struct
	byte, _ := json.MarshalIndent(data, "", "\t")
	fmt.Println(string(byte))
}

func Output(data any) []byte {
	byte, _ := json.Marshal(data)
	return byte
}
