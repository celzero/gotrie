package trie

import ( "fmt"
"io/ioutil"
)
import "bytes"
import "encoding/binary"
import "strconv"
import "strings"

var Debug = false
var RD = RankDirectory{}
var FT = FrozenTrie{}
var RD_buf = []uint16{}
var TD_buf = []uint16{}
var Blacklistconfigjson = "./filetag.json"
var W = 16
var L1 = 32*32;
var L2 = 32;
var NodeCount *int;
func Build()(error,FrozenTrie){
	var err error;
	TD_buf,err = read_file_u16("./td.txt")
	if(err != nil) {
		fmt.Println(err)
		return err,FT
	}
	RD_buf,err = read_file_u16("./rd.txt")
	if(err != nil) {
		fmt.Println(err)
		return err,FT
	}

	NodeCount,err = Read_nodecount("./node.txt")
	if(err != nil) {
		fmt.Println(err)
		return err,FT
	}
	fmt.Printf("TD_buf Length : %d\n",len(TD_buf))
	fmt.Printf("RD_buf Length : %d\n",len(RD_buf))

	RD.Init(RD_buf,TD_buf,*NodeCount * 2 + 1,L1,L2,nil)
	//RD.display()
	FT.Init(TD_buf,RD,*NodeCount)
	FT.LoadTag()


	
	/*
	var str_uint8 = []uint8{}
	var arr = []uint64{}
	var tag = []string{}
	var status bool
	str_uint8,err = TxtEncode("101.ru")
	fmt.Println(str_uint8)
	status,arr = FT.lookup(str_uint8)
	fmt.Println(status)
	fmt.Println(arr)
	if(status){
		tag = FT.FlagstoTag(arr)
		fmt.Println(tag)
	}*/
	return nil,FT
}



func read_file_u16(path string)([]uint16,error){
	content, err := ioutil.ReadFile(path)
	if (err != nil) {
		fmt.Println("Error At read file : build.go -> read_file_u16()")
		return nil,err
	}

	fmt.Println("file read successful : "+ path)
	fmt.Println("file byte length : ",len(content))
	r := bytes.NewReader(content)
	tmp16 := make([]uint16,len(content)/2)
	err = binary.Read(r, binary.LittleEndian, &tmp16)
	if(err != nil){
		fmt.Println("Error At byte to uint16 conversion : build.go -> read_file_u16()")
		return nil,err
	}
	return tmp16,err
}

func Read_nodecount(path string)(*int,error){
	content, err := ioutil.ReadFile(path)
	if (err != nil) {
		fmt.Println("Error At read file : build.go -> read_nodecount()")
		return nil,err
	}
	nodecount, _ := strconv.Atoi(strings.TrimSpace(string(content)))
	fmt.Println("file read successful : "+ path)
	fmt.Println("node count : ",nodecount)

	return &nodecount,err
}