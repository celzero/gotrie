package trie

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
)

var Debug = false
var W = 16
var L1 = 32 * 32
var L2 = 32

func Build(tdpath, rdpath, bcpath, ftpath string) (FT *FrozenTrie, err error) {
	// FIXME: add an integrity check for all four files which are
	// dependant on each other and need to be from the same "generation"
	FT = new(FrozenTrie)
	var RD = new(RankDirectory)
	TD_buf, err := read_file_u16(tdpath)
	if err != nil {
		fmt.Println(err)
		return
	}

	RD_buf, err := read_file_u16(rdpath)
	if err != nil {
		fmt.Println(err)
		return
	}

	bconfig, err := LoadBasicConfig(bcpath)
	if err != nil {
		return
	}

	if Debug {
		TD_buf_md5 := MD5Hex(castToBytes(TD_buf))
		RD_buf_md5 := MD5Hex(castToBytes(RD_buf))
		tdmd5, _ := bconfig["tdmd5"].(string)
		rdmd5, _ := bconfig["rdmd5"].(string)
		fmt.Printf("md5(TD): %s <-> %s | md5(RD): %s <-> %s | nc: %d\n", TD_buf_md5, tdmd5, RD_buf_md5, rdmd5, nodecount)
	}

	nodecount := int(bconfig["nodecount"].(float64))

	rdb := NewBStr(RD_buf)
	tdb := NewBStr(TD_buf)
	RD.Init(rdb, tdb, nodecount*2+1, L1, L2)
	FT.Init(tdb, RD, nodecount, ftpath)

	return
}

func read_file_u16(path string) (*[]uint16, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("err read file : build.go -> read_file_u16()")
		return nil, err
	}

	fmt.Printf("read from %s len %d\n", path, len(content))
	r := bytes.NewReader(content)
	tmp16 := make([]uint16, len(content)/2)
	err = binary.Read(r, binary.LittleEndian, &tmp16)
	if err != nil {
		fmt.Println("err byte2uint: build.go -> read_file_u16()")
		return nil, err
	}
	return &tmp16, err
}

func LoadBasicConfig(filepath string) (map[string]any, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	var jobj map[string]any
	err = json.Unmarshal(data, &jobj)
	if err != nil {
		fmt.Println("could not read basicconfig:", err)
		return nil, err
	}

	return jobj, err
}
