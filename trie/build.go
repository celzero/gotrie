package trie

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var Debug = false
var W = 16
var L1 = 32 * 32
var L2 = 32

func Build(tdpath, rdpath, bcpath, ftpath string, usemmap bool) (ftrie *FrozenTrie, err error) {
	td16, td8, err := readBinary(tdpath, usemmap)
	if err != nil {
		return
	}

	rd16, rd8, err := readBinary(rdpath, usemmap)
	if err != nil {
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
		err = errors.New(hstatus)
		return
	}

	rdb := NewBStr(rd16)
	tdb := NewBStr(td16)
	rdir := NewRankDir(rdb, tdb, nodecount*2+1, L1, L2)
	ftrie = NewFrozenTrie(tdb, rdir, nodecount, ftpath)

	return
}

// mmapBinary mmaps the binary file and returns the uint16 and byte slices
func mmapBinary(path string) (*[]uint16, *[]byte, error) {
	if path == "" {
		return nil, nil, errors.New("trie: empty path")
	}

	if !lilbo(false) {
		return nil, nil, errors.New("trie: mmap on little endian only")
	}

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	// Get the file size
	fi, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}
	size := fi.Size()

	data, err := syscall.Mmap(int(file.Fd()), 0, int(size), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("trie: mmaped", path, "size", len(data))

	return bytesToUint16(&data), &data, nil
}

func readBinary(path string, usemmap bool) (*[]uint16, *[]byte, error) {
	isLittleEndian := lilbo(false)
	if usemmap && isLittleEndian {
		return mmapBinary(path)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("trie: err read file : build.go -> read_file_u16()")
		return nil, nil, err
	}

	fmt.Printf("trie: read from %s len %d\n", path, len(content))

	// works only on little endian machines: go.dev/play/p/50t1HxCr9DV
	if isLittleEndian {
		return bytesToUint16(&content), &content, nil
	}

	r := bytes.NewReader(content)
	tmp16 := make([]uint16, len(content)/2)
	err = binary.Read(r, binary.LittleEndian, &tmp16)
	if err != nil {
		return nil, nil, fmt.Errorf("trie: err: %v", err)
	}
	return &tmp16, &content, nil
}

func readBasicConfig(filepath string) (map[string]any, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var jobj map[string]any
	err = json.Unmarshal(data, &jobj)
	if err != nil {
		return nil, fmt.Errorf("trie: err reading basicconfig: %v", err)
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
