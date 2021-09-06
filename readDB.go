package main

import (
    "os"
    "bufio"
    //"bytes"
    "database/sql"
    "fmt"
    "strings"
    //"strconv"
    //"time"
    //"math/rand"
    _ "github.com/mattn/go-sqlite3"
)


func WordScan(filename string) []string{
    file, err := os.Open(filename)
    if err != nil {
      fmt.Println(err)
      return nil
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanWords)

    var words []string
    //words := []string{}
    for scanner.Scan() {
        word := scanner.Text()
        if (strings.HasSuffix(word, ".") || strings.HasSuffix(word, ",") ||strings.HasSuffix(word, "'")){
            word = word[0 : len(word)-1]
        }
        if strings.HasSuffix(word, "’s"){
            word = word[0 : len(word)-2]
        }
        if strings.HasPrefix(word, "\""){
            word = word[1 : len(word)-1]
        }
        words = append(words, word)
    }

    // fmt.Println("word list:")
    // for k,v := range words {
    //     fmt.Println(strconv.Itoa(k) + ":" + v)
    // }
    return words

}

func NoNeedPhoneitcWords(word string) bool {

    wordlist := [...]string{"one", "a", "the", "and", "said", "had", "in", "by", "state", "they", "than", "after", "year",
                            "an", "soldiers", "less", "won", "on", "television", "president", "of", "west", "be", "week",
                            "borders", "closed", "country", "group", "taking", "our", "for", "third", "man", "he", "own",
                            "TV", "broadcast", "adding", "were", "being", "we", "are", "that", "would"}

    for _, name := range wordlist{
    	if word==name {
    		return true
    	}
    }

    if strings.Contains(word, "é"){
        return true
    }

    return false
}

func getphonetic(db *sql.DB, word string)string {
    rows, _ := db.Query("SELECT id, phonetic FROM stardict where word='" + word + "'")
    //fmt.Println("SELECT id, phonetic FROM stardict where word='" + word + "'")
    var id int
    var phonetic string
    i:=0
    for rows.Next() {
        rows.Scan(&id, &phonetic)
        //fmt.Println(strconv.Itoa(id) + ": " + word + " ["+ phonetic + "] ")
        return "["+ phonetic + "]"
        i = i+1
        if i>10 {break}
    }
    rows.Close()

    return "##NOT FOUND##"
}

func main(){
    db, _ := sql.Open("sqlite3", "stardict.db")
	defer db.Close()

    wordlist := WordScan("test.txt")
    for _, word := range wordlist{
        if NoNeedPhoneitcWords(word){
            fmt.Print(word + " ")
        } else {
            fmt.Print(word + getphonetic(db, word) + " ")
        }
        
    }
    
}