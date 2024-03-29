package trie

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"unsafe"
)

var Debug = false
var W = 16
var L1 = 32 * 32
var L2 = 32

func Build(tdpath, rdpath, bcpath, ftpath string) (ftrie *FrozenTrie, err error) {
	td16, td8, err := readBinary(tdpath)
	if err != nil {
		fmt.Println(err)
		return
	}

	rd16, rd8, err := readBinary(rdpath)
	if err != nil {
		fmt.Println(err)
		return
	}

	bconfig, err := readBasicConfig(bcpath)
	if err != nil {
		return
	}

	nodecount := int(bconfig["nodecount"].(float64))

	tdmd5hex := MD5Hex(td8)
	rdmd5hex := MD5Hex(rd8)
	tdmd5, _ := bconfig["tdmd5"].(string)
	rdmd5, _ := bconfig["rdmd5"].(string)
	hstatus := fmt.Sprintf("md5 mismatch: %s <-> %s | %s <-> %s", tdmd5hex, tdmd5, rdmd5hex, rdmd5)

	if tdmd5hex != tdmd5 || rdmd5hex != rdmd5 {
		fmt.Println(hstatus)
		err = errors.New(hstatus)
		return
	}

	rdb := NewBStr(rd16)
	tdb := NewBStr(td16)
	rdir := NewRankDir(rdb, tdb, nodecount*2+1, L1, L2)
	ftrie = NewFrozenTrie(tdb, rdir, nodecount, ftpath)

	return
}

func readBinary(path string) (*[]uint16, *[]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("err read file : build.go -> read_file_u16()")
		return nil, nil, err
	}

	fmt.Printf("read from %s len %d\n", path, len(content))

	// works only on little endian machines: go.dev/play/p/50t1HxCr9DV
	if lilbo(false) {
		return bytesToUint16(&content), &content, nil
	}

	r := bytes.NewReader(content)
	tmp16 := make([]uint16, len(content)/2)
	err = binary.Read(r, binary.LittleEndian, &tmp16)
	if err != nil {
		fmt.Println("err byte2uint: build.go -> read_file_u16()")
		return nil, nil, err
	}
	return &tmp16, &content, nil
}

func readBasicConfig(filepath string) (map[string]any, error) {
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

// little endian byte order?
// stackoverflow.com/a/53286786
func lilbo(unknown bool) bool {
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		return true
	case [2]byte{0xAB, 0xCD}:
		return false
	default:
		return unknown
	}
}
