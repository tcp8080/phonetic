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

func BasicScan(filename string) []string{
    file, err := os.Open(filename)
    if err != nil {
      fmt.Println(err)
      return nil
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanLines)

    // This is our buffer now
    var lines []string

    for scanner.Scan() {
      lines = append(lines, scanner.Text())
    }

    // for _, line := range lines {
    //   fmt.Println(line)
    //   fmt.Println("-----------------------------------------\n")
    // }
    return lines
}

func WordScan(line string) []string{

    scanner := bufio.NewScanner(strings.NewReader(line))
    scanner.Split(bufio.ScanWords)

    var words []string
    for scanner.Scan() {
        word := scanner.Text()
        if (strings.HasSuffix(word, ".") || strings.HasSuffix(word, ",") ||strings.HasSuffix(word, "'")||strings.HasSuffix(word, ";")){
            word = word[0 : len(word)-1]
        }
        // if strings.HasSuffix(word, "’s"){
        //     word = word[0 : len(word)-2]
        // }
        // if strings.HasPrefix(word, "“"){
        //     word = word[3 : len(word)-3]
        // }
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
                            "TV", "broadcast", "adding", "were", "being", "we", "are", "that", "would", "the", "from",
                            "is", "recovery", "economic", "is", "Mr", "come", "go", "back", "many", "as", "clear", "big", "small",
                            "what", "where", "when", "how", "any", "ways", "have", "some", "to", "this", "that", "whether", "which",
                        "way", "take", "parts", "but", "still", "weak", "good", "their", "at", "more", "use", "less", "with", "everyone",
                    "has", "only", "level", "data", "it", "was", "three"}

    for _, name := range wordlist{
    	//if word==name {
        if strings.EqualFold(word, name) {
    		return true
    	}
    }

    if (strings.Contains(word, "é") || strings.Contains(word, "-") || strings.Contains(word, "%") || strings.Contains(word, "(") ||strings.Contains(word, "'") ||
        strings.Contains(word, "0") || strings.Contains(word, "1") || strings.Contains(word, "2") || strings.Contains(word, ")") ||strings.Contains(word, "”") ||
        strings.Contains(word, "3") || strings.Contains(word, "4") || strings.Contains(word, "5") || strings.Contains(word, "“") ||
        strings.Contains(word, "6") || strings.Contains(word, "7") || strings.Contains(word, "8") || strings.Contains(word, "9") ){
        return true
    }
    return false
}

func getphonetic(db *sql.DB, word string)string {
    rows, _ := db.Query("SELECT id, phonetic FROM stardict where word='" + word + "'")
    //fmt.Println("SELECT id, phonetic FROM stardict where word='" + word + "'")
    if rows == nil {
        return "##NOTHING FOUND With Word: " + word +"##"
    }
    var id int
    var phonetic string
    i:=0
    for rows.Next() {
        rows.Scan(&id, &phonetic)
        //fmt.Println(strconv.Itoa(id) + ": " + word + " ["+ phonetic + "] ")
        if phonetic == "" {
            return "##NOT Filled IN DB##"
        }
        return "["+ phonetic + "]"
        i = i+1
        if i>10 {break}
    }
    rows.Close()

    return "##NOT FOUND##"
}

func getHtmlWord(word string, phonetc string) string{
    content := "<ruby>"
    if strings.HasPrefix(phonetc, "[") { // phonetic need to print
        content = content + word + "<rt>" + phonetc + "</rt>"
    } else {
        bnumber := len(phonetc)
        content = content + word + "<rt>" + strings.Repeat(" ", bnumber) + "</rt>"
    }
	content = content + "</ruby>&nbsp;"
    //fmt.Println("------------------" + content+"---------------------")
	return content
}
func getHtmlSentence(db *sql.DB, sentence string) string{
    s := "<p><span style=\"display:inline-block;width:60%;word-wrap:break-word;white-space:normal;\">"
    wordlist := strings.Split(sentence, " ")
    for _, c := range wordlist {
        word := strings.Trim(c, " ")
        p := word
        if NoNeedPhoneitcWords(word){ 
            p = word
        } else {
            p = getphonetic(db, word)
        }
       
        if strings.Contains(p, "##") { // no phonetic found
            p = word
        }
        
        s = s + getHtmlWord(word, p)
    } 
    
    return s + "</span></p>"
}

func main(){
    db, _ := sql.Open("sqlite3", "stardict.db")
	defer db.Close()

    // lines := BasicScan("test.txt")
    // for _, line := range lines{
    //     wordlist := WordScan(line)
    //     for _, word := range wordlist {
    //         if NoNeedPhoneitcWords(word){
    //             fmt.Print(word + " ")
    //         } else {
    //             fmt.Print(word + getphonetic(db, word) + " ")
    //         }
    //     }
    //    fmt.Println()
    // }

    lines := BasicScan("test.txt")
    s := "　　<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" /><body style=\"background-color:powderblue;\">"
    for _, line := range lines{
        s = s + getHtmlSentence(db, line)
    }
    s = s+"</body>"
    // fmt.Println(s)
    f, e := os.Create("test111.html")
    if e != nil {
        panic(e)
    }
    defer f.Close()
    _, err := f.WriteString(s)
    if err != nil {
        panic(err)
    }
    f.Sync()
}