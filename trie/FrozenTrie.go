package trie


import "math"
import "fmt"
import (
    "encoding/json"
	"io/ioutil"
	"strings"
	b64 "encoding/base64"
	"net/url"
)
type FrozenTrie struct{
	data BS
	directory RankDirectory
	extraBit int
	bitslen int
	letterStart int
	valuesStart int
	valuesIndexLength int
	valuesDirBitsLength int
	rflags []string
	fdata map[string]interface{}
	bcache *Cache
	blen *int
	fcache *Cache
	flen *int
	blimt int
	flimt int

	usr_flag string
	usr_bl []string
}

func (FT *FrozenTrie) Init(trieData []uint16, rdir RankDirectory, nodeCount int){
	FT.data = BS{}
	FT.data.Init(trieData)

	FT.directory = rdir
	FT.extraBit = 1 //(config.compress && !config.unroll) ? 1 : 0;
	FT.bitslen = 9 + FT.extraBit //((config.base32) ? 6 : 9) + this.extraBit;
	FT.letterStart = nodeCount * 2 + 1;

	FT.valuesStart = FT.letterStart + (nodeCount * FT.bitslen); // + 1;

	FT.valuesIndexLength = int(math.Ceil(math.Log2(float64(nodeCount))));

	FT.valuesDirBitsLength = int(math.Ceil(math.Log2(float64(FT.data.length - FT.valuesStart))));

	FT.bcache = New()
	FT.fcache = New()
	var a = 0
	FT.blen = &a
	var b = 0
	FT.flen = &b
	FT.blimt = 2500
	FT.flimt = 2500

	FT.fdata = make(map[string]interface{})

	FT.usr_flag = ""
	FT.usr_bl  = []string{}
}

func (FT FrozenTrie) getNodeByIndex(index int)FrozenTrieNode{
	FTN := FrozenTrieNode{}
	FTN.Init(FT,index)
	return FTN
}

func (FT FrozenTrie) getRoot()FrozenTrieNode {
	return FT.getNodeByIndex(0)
}

func (FT FrozenTrie) lookup(word []uint8)(bool,[]uint32){	
	var node = FT.getRoot()	
	var emptyreturn []uint32
	
	for i := 0; i < len(word); i++ {
		var isFlag = -1;
		var that FrozenTrieNode
		for {
			that = node.getChild(isFlag + 1);
			if (!that.flag()) {break;}
			isFlag += 1;
			if(!((isFlag + 1) < node.getChildCount())){
				break;
			}
		}
		var minChild = isFlag;
		if (Debug){
			fmt.Printf("            count: %d i: %d  w: %d  nl: %d  flag: %d\n",node.getChildCount() , i ,word[i] ,node.letter() ,isFlag)
		} 
		if ((node.getChildCount()-1) <= minChild) {			
			return false,emptyreturn;
		}
		//if(config.compress === true && !config.unroll)
		var high = node.getChildCount();
		var low = isFlag;
		var child FrozenTrieNode
		for (high-low) > 1 {
			var probe = (high + low) / 2 | 0
			child = node.getChild(probe)
			var prevchild *FrozenTrieNode

			
			if(probe > isFlag){
				var tmp =  node.getChild(probe - 1)
				prevchild = &tmp
			}else{
				prevchild = nil
			}

			if (Debug){
				fmt.Printf("            current: %d l: %d h: %d w: %d\n",child.letter() , low ,high ,word[i])
				//return false,emptyreturn
			}

			if(child.compressed() || (prevchild != nil && (prevchild.compressed() && !prevchild.flag()))){
				var startchild []FrozenTrieNode 
				var endchild []FrozenTrieNode
				var temp FrozenTrieNode
				var start = 0;
				var end = 0;
				startchild = append(startchild,child)
				start = start + 1

				for {
					temp = node.getChild(probe - start)
					if (!temp.compressed()){ break;}
					if (temp.flag()){ break; }
					startchild = append(startchild,temp)
					start = start + 1
				}
				if(Debug){
					fmt.Printf("  check: letter : %d  word : %d start: %d\n",startchild[start - 1].letter(),word[i],start)
				}
				
				if (uint8(startchild[start - 1].letter()) > word[i]) {
					if (Debug) {
						fmt.Printf("            shrinkh start: %d s: %d w: %d\n",startchild[start-1].letter() , start , word[i])
					}
										
					high = probe - start + 1;
					if (high - low <= 1) {	
						if(Debug){
							fmt.Printf("    (high - low ): %d c: %d h: %d l: %d cl: %d w: %d pr: %d\n",(high - low),node.getChildCount(), high, low,child.letter(), word[i], probe)
						}												
						return false,emptyreturn;
					}
					continue;
				}

				if (child.compressed()) {
					for {
						end = end + 1
						temp = node.getChild(probe + end);
						endchild = append(endchild,temp);
						if (!temp.compressed()){ break; }
					}
				}

				if(uint8(startchild[start - 1].letter()) < word[i]){
					if (Debug) {
						fmt.Printf("            shrinkh start: %d s: %d w: %d\n",startchild[start-1].letter() , start , word[i])
					}
					low = probe + end;
					if (high - low <= 1) {
						if(Debug){
							fmt.Printf("    (high - low ): %d c: %d h: %d l: %d cl: %d w: %d pr: %d\n",(high - low),node.getChildCount(), high, low,child.letter(), word[i], probe)
						}
						return false,emptyreturn;
					}
					continue;
				}


				for ii, jj := 0, len(startchild)-1; ii < jj; ii, jj = ii+1, jj-1 {
					startchild[ii], startchild[jj] = startchild[jj], startchild[ii]
				}
				var nodes = append(startchild,endchild...)
				var comp []uint8
				for inc:=0;inc<len(nodes);inc++ {
					comp = append(comp,uint8(nodes[inc].letter()))
				}

				
				var sliceend = i + len(comp)
				if(sliceend > len(word)){
					sliceend = len(word)
				}
				var w = word[i:sliceend]

				if(Debug){
					fmt.Printf("p: %d comp: %v w: %v\n",probe,comp,w)
				}
				if(len(w) < len(comp)){
					return false,emptyreturn;
				}
				for inc:=0;inc<len(comp);inc++ {
					if(w[inc] != comp[inc]){
						return false,emptyreturn;
					}
				}
				child = nodes[len(nodes) - 1];
				i += len(comp) - 1; 
				break;
			}else{
				if ( uint8(child.letter()) == word[i] ) {
						break;
				} else if ( word[i] > uint8(child.letter()) ) {
					low = probe;
				} else {
					high = probe;
				}
			}

			if (high - low <= 1) {
				return false,emptyreturn;
			}	
		}
		if(Debug){
			fmt.Printf("        next: %d\n" ,child.letter())
		}
		node = child;
	}
	if(node.final()){
		return node.final(),node.value()
	}
	return false,emptyreturn
}


func(FT *FrozenTrie) LoadTag()error{
	data, err := ioutil.ReadFile(Blacklistconfigjson)
    if err != nil {
	  fmt.Print(err)
	  return err
	}
	
	var obj map[string]interface{}
	err = json.Unmarshal(data, &obj)
    if err != nil {
		fmt.Println("error:", err)
		return err
	}
	FT.rflags = make([]string,len(obj))
	for key, _ := range obj {
		indata := obj[key].(map[string]interface{})
		var index = int(indata["value"].(float64))
		var data = indata["uname"].(string)
		FT.rflags[index] = data 
		
		FT.fdata[key] = indata
	}
	return nil
}

func(FT FrozenTrie) FlagstoTag(flags []uint32)[]string{
	// flags has to be an array of 16-bit integers.
	header := uint16(flags[0]);
	tagIndices := []int{};
	values := []string{};
	var mask uint16
	mask = 0x8000
	//fmt.Println(FT.rflags)
	for i := 0 ; i < 16; i++ {
		if ((header << i) == 0) {break;}
		if ((header & mask) == mask) {
			tagIndices = append(tagIndices,i);
		}
		mask = mask >> 1;
	}
	// flags.length must be equal to tagIndices.length
	if (len(tagIndices) != (len(flags) - 1)) {
		fmt.Printf("%v %v flags and header mismatch (bug in upsert?)",tagIndices, flags);
		return values;
	}
	for i := 1; i < len(flags); i++ {
		//fmt.Println("i : ",i)
		var flag = uint16(flags[i]);
		var index = tagIndices[i-1]
		mask = 0x8000
		for j := 0 ; j < 16; j++ {
			if ((flag << j) == 0){ break; }
			if ((flag & mask) == mask) {
				var pos = (index * 16) + j;
				if(Debug){
					fmt.Printf("pos %d  index/tagIndices %d %v j/i %d %d\n",pos,index,tagIndices,j,i)
				}
				//console.log("pos " , pos, "index/tagIndices", index, tagIndices, "j/i", j , i);
				values = append(values,FT.rflags[pos]);
			}
			mask = mask >> 1;
		}
	}
	return values;
}




func (FT *FrozenTrie) DNlookup(dn string,usr_flag string)(bool,[]string){

	if(FT.usr_flag == "" || FT.usr_flag != usr_flag){
		//fmt.Println("User config saved : ")
		FT.usr_flag = usr_flag
		FT.usr_bl = FT.Urlenc_to_flag(FT.usr_flag)
	}


	dn = strings.TrimSpace(dn)
	bvalue := FT.bcache.Get(dn)
	fvalue := FT.fcache.Get(dn)
	var blfname []string
	var retlist []string
	var found = false
	if(bvalue != nil){
		//fmt.Println("Return frm B-Cache : ")
		blfname = strings.Split(bvalue.(string), "-")
		found,retlist = Find_Lista_Listb(FT.usr_bl,blfname)	
		return found,retlist
	}
	if(fvalue != nil){
		//fmt.Println("Return frm F-Cache : ")
		return found,blfname
	}


	var arr = []uint32{}
	var tag = []string{}
	var status bool
	str_uint8,_ := TxtEncode(dn)
	status,arr = FT.lookup(str_uint8)
	if(status){
		tag = FT.FlagstoTag(arr)

		found,retlist = Find_Lista_Listb(FT.usr_bl,tag)

		if(*FT.blen >= FT.blimt){
			FT.bcache.Evict(1)
			*(FT.blen)--
			//fmt.Println("B-Cache Full Evicted : 1")
		}
		FT.bcache.Set(dn,strings.Join(tag,"-"))
		*(FT.blen)++
		//fmt.Println("Add to B-Cache lenght : ",*FT.blen)
		return found,retlist
	} else {
		if(*FT.flen >= FT.flimt){
			FT.fcache.Evict(1)
			*(FT.flen)--
			//fmt.Println("F-Cache Full Evicted : 1")
		}
		FT.fcache.Set(dn,"")
		*(FT.flen)++
		//fmt.Println("Add to F-Cache lenght : ",*FT.flen)
		return found,retlist
	}	
}


func (FT FrozenTrie) CreateUrlEncodedflag(fl []string)(string){
	var val = 0
	var res = ""
	for _,flag := range fl{
		indata := FT.fdata[flag].(map[string]interface{})
		//fmt.Println(indata)
		val = int(indata["value"].(float64))
		//header := 0
		index := ((val / 16) | 0)
		pos := val % 16
		dataIndex1 := 0
		h:=0
		if(len(res)>=1){
			h = DEC16(res,0)
		}

		//fmt.Println("Mask Bottom : ",uint(FT.data.MaskBottom[16][16 - index]))
		dataIndex := FT.data.countSetBits(h & int(FT.data.MaskBottom[16][16 - index])) + 1;


		n :=0
		if((((h >> (15 - index)) & 0x1) != 1)){
			n=0
		}else{
			n = DEC16(res,dataIndex)
		}

		upsertData := (n != 0)
		h |= 1 << (15 - index);
		n |= 1 << (15 - pos);

		//fmt.Println(Flag_to_uint(res))
		if(len(res)>=2){
			dataIndex1 = dataIndex
			if(upsertData){
				dataIndex1 = dataIndex1 + 1
			}	
			res = CHR16(h) + FlagSubstring(res,1,dataIndex) + CHR16(n) + FlagSubstring(res,dataIndex1,0)	
		}else{
			res = CHR16(h) + CHR16(n)
		}

		//fmt.Println("h : ",h)
		//fmt.Println("n : ",n)
		//fmt.Println("dataindex : ",dataIndex)
		//fmt.Println("dataindex1 : ",dataIndex1)
		//fmt.Println("Index : ",index)
		//fmt.Println("Pos : ",pos)
	}
	//temp := Flag_to_uint(res)
	//fmt.Println(temp)
	//fmt.Println(FT.FlagstoTag(temp))
	//fmt.Println("base 64 encode string :",b64.StdEncoding.EncodeToString([]byte(res)))
	//fmt.Println("url encode string :" ,url.QueryEscape(b64.StdEncoding.EncodeToString([]byte(res))))
	return url.QueryEscape(b64.StdEncoding.EncodeToString([]byte(res)))
}

func (FT FrozenTrie) Urlenc_to_flag(str string)([]string){
	str,_ = url.QueryUnescape(str)
	buf,_ := b64.StdEncoding.DecodeString(str)
	str = string(buf)
	return FT.FlagstoTag(Flag_to_uint(str))
}