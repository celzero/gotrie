package trie

import (
	"fmt"
	"io/ioutil"
)
import "bytes"
import "encoding/binary"
import "encoding/json"

var Debug = false
var W = 16
var L1 = 32 * 32
var L2 = 32

func Build(tdpath, rdpath, bcpath, ftpath string) (error, FrozenTrie) {
	// FIXME: add an integrity check for all four files which are
	// dependant on each other and need to be from the same "generation"
	var RD = RankDirectory{}
	var FT = FrozenTrie{}
	var RD_buf = []uint16{}
	var TD_buf = []uint16{}
	var err error
	TD_buf, err = read_file_u16(tdpath)
	if err != nil {
		fmt.Println(err)
		return err, FT
	}

	RD_buf, err = read_file_u16(rdpath)
	if err != nil {
		fmt.Println(err)
		return err, FT
	}

	bconfig, err := LoadBasicConfig(bcpath)
	if err != nil {
		return err, FT
	}

	TD_buf_md5 := MD5Hex(castToBytes(TD_buf))
	RD_buf_md5 := MD5Hex(castToBytes(RD_buf))

	nodecount := int(bconfig["nodecount"].(float64))
	tdmd5, _ := bconfig["tdmd5"].(string)
	rdmd5, _ := bconfig["rdmd5"].(string)
	fmt.Printf("md5(TD): %s <-> %s | md5(RD): %s <-> %s | nc: %d\n", TD_buf_md5, tdmd5, RD_buf_md5, rdmd5, nodecount)

	RD.Init(RD_buf, TD_buf, nodecount*2+1, L1, L2, nil)
	FT.Init(TD_buf, RD, nodecount)
	FT.LoadTag(ftpath)

	return nil, FT
}

func read_file_u16(path string) ([]uint16, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error At read file : build.go -> read_file_u16()")
		return nil, err
	}

	fmt.Println("read from %s len %d", path, len(content))
	r := bytes.NewReader(content)
	tmp16 := make([]uint16, len(content)/2)
	err = binary.Read(r, binary.LittleEndian, &tmp16)
	if err != nil {
		fmt.Println("Error At byte to uint16 conversion : build.go -> read_file_u16()")
		return nil, err
	}
	return tmp16, err
}

func LoadBasicConfig(filepath string) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	var jobj map[string]interface{}
	err = json.Unmarshal(data, &jobj)
	if err != nil {
		fmt.Println("could not read basicconfig:", err)
		return nil, err
	}

	return jobj, err
}
