package main
import "./trie"
import "fmt"
import "time"
import "os"
import "bufio"
import "strings"


func main()  {
	

	err,FT := trie.Build()
	if(err == nil){
		res := []string{"DAH","ADH","BXW", "BQJ"}
		usr_flag := FT.CreateUrlEncodedflag(res)
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
		}
	}else{
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

