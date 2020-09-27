package main

import "github.com/celzero/gotrie/trie"
import "fmt"

//import "time"
//import "os"
//import "bufio"
//import "strings"

func main() {

	err, FT := trie.Build("./td.txt", "./rd.txt", "./basicconfig.json", "./filetag.json")
	if err == nil {
		//[33216 32768 8192 256 4]
		//6IeA6ICA4oCAxIAE
		//res := []string{"AMI","CQT","EOK","MTF"}
		//usr_flag := FT.CreateUrlEncodedflag(res)
		//fmt.Println(usr_flag)
		fmt.Println(FT.Urlenc_to_flag("w4DEgAQ="))
		fmt.Println(FT.DNlookup("google.com", "6IeA6ICA4oCAxIAE"))
		/*
			for{
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Enter text: ")
				text, _ := reader.ReadString('\n')
				text = strings.TrimSpace(text)
				if(text == "exit"){break;}
				start := time.Now()
				fmt.Println(FT.DNlookup(text,usr_flag))
				elapsed := time.Since(start)
				fmt.Printf("Time Diff %s\n",elapsed)
			}*/
	} else {
		fmt.Println("Error at trie Build")
	}

	/*err,FT := trie.Build()
	if(err == nil){
		res := []string{"DAH","ADH","BXW", "BQJ"}
		fmt.Println("Base64 to flag : ",FT.Urlenc_to_flag(FT.CreateUrlEncodedflag(res)))
	}*/

	/*
		start := time.Now()
		trie.CheckDN1()
		elapsed := time.Since(start)
		fmt.Printf("Time Diff %s\n",elapsed)*/

}
