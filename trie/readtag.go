package trie

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CheckDN(FT *FrozenTrie) {
	var a [1]string
	a[0] = "D:\\Celzero\\blocklist script\\blocklistfiles\\privacy\\MNH.txt"
	//a[0] = "D:\\Celzero\\blocklist script\\blocklistfiles\\privacy\\NLH.txt"
	/*a[0] = "D:\\Celzero\\blocklist script\\blocklistfiles\\privacy\\ABY.txt"
	  a[1] = "D:\\Celzero\\blocklist script\\blocklistfiles\\privacy\\ADH.txt"
	  a[2] = "D:\\Celzero\\blocklist script\\blocklistfiles\\privacy\\ARL.txt"
	  a[3] = "D:\\Celzero\\blocklist script\\blocklistfiles\\privacy\\AUT.txt"
	  a[4] = "D:\\Celzero\\blocklist script\\blocklistfiles\\privacy\\BQJ.txt"
	  a[5] = "D:\\Celzero\\blocklist script\\blocklistfiles\\privacy\\BTY.txt"
	  a[6] = "D:\\Celzero\\blocklist script\\blocklistfiles\\privacy\\BXW.txt"
	  a[7] = "D:\\Celzero\\blocklist script\\blocklistfiles\\privacy\\DAH.txt"
	  a[8] = "D:\\Celzero\\blocklist script\\blocklistfiles\\privacy\\DDB.txt"
	  a[9] = "D:\\Celzero\\blocklist script\\blocklistfiles\\privacy\\DGJ.txt"*/

	var failcount = 0
	var totalcount = 0
	//var filecount =0
	for _, str := range a {
		file, err := os.Open(str)

		if err != nil {
			fmt.Printf("failed opening file: %s\n", err)
		}

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		fmt.Println("Reading File : ", str)
		//filecount = 0
		for scanner.Scan() {
			totalcount += 1
			var dn = scanner.Text()
			//fmt.Println(dn)
			//fmt.Println(str_uint8)
			status, _ := FT.DNlookup(strings.TrimSpace(dn), "")
			if !status {
				//fmt.Println("Fail not found in Trie")
				failcount += 1
			} else {
				//fmt.Println(arr)
			}
		}

		file.Close()
	}
	fmt.Println("failcount :", failcount)
	fmt.Println("Total Count :", totalcount)
}

func CheckDN1() {
	//var a [1]string
	//a[0] = "D:\\Celzero\\blocklist script\\blocklistfiles\\privacy\\ARL.txt"
	//a[0] = "D:\\Celzero\\blocklist script\\blocklistfiles\\privacy\\NLH.txt"
	/*a[0] = "D:\\Celzero\\blocklist script\\blocklistfiles\\privacy\\ABY.txt"
	  a[1] = "D:\\Celzero\\blocklist script\\blocklistfiles\\privacy\\ADH.txt"
	  a[2] = "D:\\Celzero\\blocklist script\\blocklistfiles\\privacy\\ARL.txt"
	  a[3] = "D:\\Celzero\\blocklist script\\blocklistfiles\\privacy\\AUT.txt"
	  a[4] = "D:\\Celzero\\blocklist script\\blocklistfiles\\privacy\\BQJ.txt"
	  a[5] = "D:\\Celzero\\blocklist script\\blocklistfiles\\privacy\\BTY.txt"
	  a[6] = "D:\\Celzero\\blocklist script\\blocklistfiles\\privacy\\BXW.txt"
	  a[7] = "D:\\Celzero\\blocklist script\\blocklistfiles\\privacy\\DAH.txt"
	  a[8] = "D:\\Celzero\\blocklist script\\blocklistfiles\\privacy\\DDB.txt"
	  a[9] = "D:\\Celzero\\blocklist script\\blocklistfiles\\privacy\\DGJ.txt"*/

	var failcount = 0
	var totalcount = 0
	var filecount = 0
	var str_uint8 = []uint8{}
	var status bool
	//res := []string{"DAH","ADH","BXW", "BQJ"}
	//usr_flag := FT.CreateUrlEncodedflag(res)

	var files []string
	FT, err := Build("./td.txt", "./rd.txt", "./basicconfig.json", "./filetag.json")
	if err != nil {
		panic(err)
	}
	root := "D:\\Celzero\\blocklist script\\blocklistfiles\\parentalcontrol\\"
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	for _, str := range files {
		content, err := os.ReadFile(str)
		if err != nil {
			fmt.Printf("failed opening file: %s\n", err)
			continue
		}

		filecount += 1
		fmt.Println("Reading File : ", str)
		for _, dn := range strings.Split(string(content), "\n") {
			totalcount += 1
			//fmt.Println(dn)
			//fmt.Printf("%d : %s",totalcount,dn)
			str_uint8, _ = TxtEncode(strings.TrimSpace(dn))
			status, _ = FT.lookup(str_uint8)
			if !status {
				//fmt.Println("Fail not found in Trie")
				failcount += 1
			}
		}

		if filecount == 5 {
			break
		}
	}
	fmt.Println("filecount :", filecount)
	fmt.Println("failcount :", failcount)
	fmt.Println("Total Count :", totalcount)
}
