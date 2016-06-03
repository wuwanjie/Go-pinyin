package pinyingo

import (
	"bufio"
	"encoding/json"
	//"fmt"
	"github.com/yanyiwu/gojieba"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	STYLE_NORMAL       = 1
	STYLE_TONE         = 2
	STYLE_INITIALS     = 3
	STYLE_FIRST_LETTER = 4
	USE_SEGMENT        = true
	NO_SEGMENT         = false
	use_hmm            = true
)

var phrasesDict map[string]string
var commonWord map[string]bool
var reg *regexp.Regexp
var INITIALS []string = strings.Split("b,p,m,f,d,t,n,l,g,k,h,j,q,x,r,zh,ch,sh,z,c,s", ",")
var keyString string
var jieba *gojieba.Jieba
var sympolMap = map[string]string{
	"ā": "a1",
	"á": "a2",
	"ǎ": "a3",
	"à": "a4",
	"ē": "e1",
	"é": "e2",
	"ě": "e3",
	"è": "e4",
	"ō": "o1",
	"ó": "o2",
	"ǒ": "o3",
	"ò": "o4",
	"ī": "i1",
	"í": "i2",
	"ǐ": "i3",
	"ì": "i4",
	"ū": "u1",
	"ú": "u2",
	"ǔ": "u3",
	"ù": "u4",
	"ü": "v0",
	"ǘ": "v2",
	"ǚ": "v3",
	"ǜ": "v4",
	"ń": "n2",
	"ň": "n3",
	"": "m2",
}

func init() {
	keyString = getMapKeys()
	reg = regexp.MustCompile("([" + keyString + "])")
	dictPath := getDictPath()

	//初始化时将gojieba实例化到内存
	jieba = gojieba.NewJieba(dictPath+"jieba.dict.utf8", dictPath+"hmm_model.utf8", dictPath+"user.dict.utf8")

	// 初始化常用词
	initPhrases()
	loadPhrases()

	// 加载常用字
	commonWord = make(map[string]bool, 10000)
	loadCommonWord()
}

//func GeneDict() {
//	fHandle, err := os.OpenFile("/home/wuwanjie/heici.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm|os.ModeTemporary)
//	defer fHandle.Close()
//	if err != nil {
//		util.LogInfo(fmt.Sprintf("%v", err))
//		return
//	}
//	bufW := bufio.NewWriter(fHandle)
//
//	for k, v := range dict {
//		if len(v) == 0 {
//			continue
//		}
//		str := normalStr(v)
//		fmt.Printf("%s %s\n", v, str)
//		bufW.WriteString(fmt.Sprintf("dict[0x%x]  = \"%s\"\n", k, str))
//		bufW.Flush()
//	}
//}

func getMapKeys() string {
	keyString := ""
	for key, _ := range sympolMap {
		keyString += key
	}
	return keyString
}

func normalStr(str string) string {
	tmp := reg.ReplaceAllStringFunc(str, replaceFunc)
	return tmp
}

func replaceFunc(str string) string {
	_, ok := sympolMap[str]
	if !ok {
		return str
	}
	return string([]byte(sympolMap[str])[0])
}

//获取文件所在的根目录
func getDictPath() string {
	//currentPath, _ := os.Getwd()
	return "/home/q/data/itachi/common/"
}

func initPhrases() {
	//currentPath, _ := os.Getwd()
	f, err := os.Open(getDictPath() . "phrases-dict")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&phrasesDict); err != nil {
		log.Fatal(err)
	}
}

func loadPhrases() {
	f, err := os.Open(getDictPath() . "idf.utf8")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	bufRead := bufio.NewReader(f)
	var shuziReg *regexp.Regexp
	shuziReg = regexp.MustCompile("^[一|二|三|四|五|六|七|八|九]*$")
	for {
		line, err := bufRead.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		line = strings.Trim(line, "\r\n")
		fields := strings.Split(strings.Trim(line, " "), " ")
		if len(fields[0]) <= 0 {
			continue
		}
		if shuziReg.MatchString(fields[0]) {
			continue
		}
		phrasesDict[fields[0]] = "1"
	}
}

func loadCommonWord() {
	f, err := os.Open(getDictPath() . "common.utf8")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	bufRead := bufio.NewReader(f)
	for {
		line, err := bufRead.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		line = strings.Trim(line, "\r\n")
		for len(line) > 0 {
			r, size := utf8.DecodeRuneInString(line)
			commonWord[string(r)] = true
			line = line[size:]
		}
	}
}
