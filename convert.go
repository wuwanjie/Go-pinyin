package pinyingo

import (
	//"fmt"
	//"golang/util"
	"strings"
	"unicode/utf8"
)

type Options struct {
	style     int
	segment   bool
	heteronym bool
}

func (this *Options) perStr(pinyinStrs string) string {
	//util.LogInfo(fmt.Sprintf("[pinyinStrs] %#v", pinyinStrs))
	//util.LogInfo(fmt.Sprintf("[INITIALS] %#v", INITIALS))
	switch this.style {
	case STYLE_INITIALS:
		for i := 0; i < len(INITIALS); i++ {
			if strings.Index(pinyinStrs, INITIALS[i]) == 0 {
				return INITIALS[i]
			}
		}
		return ""
	case STYLE_TONE:
		ret := strings.Split(pinyinStrs, ",")
		//util.LogInfo(fmt.Sprintf("[STYLE_TONE] %#v", ret))
		return ret[0]
	case STYLE_NORMAL:
		ret := strings.Split(pinyinStrs, ",")
		//util.LogInfo(fmt.Sprintf("[STYLE_NORMAL] %#v", ret))
		return normalStr(ret[0])
	}
	return ""
}

func (this *Options) doConvert(strs string) []string {
	//获取字符串的长度
	bytes := []byte(strs)
	pinyinArr := make([]string, 0)
	nohans := ""
	var tempStr string
	var single string
	for len(bytes) > 0 {
		r, w := utf8.DecodeRune(bytes)
		bytes = bytes[w:]
		single = get(int(r))
		// 中文字符判断
		tempStr = string(r)
		if len(single) == 0 {
			nohans += tempStr
		} else {
			if len(nohans) > 0 {
				pinyinArr = append(pinyinArr, nohans)
				nohans = ""
			}
			pinyinArr = append(pinyinArr, this.perStr(single))
		}
	}
	//处理末尾非中文的字符串
	if len(nohans) > 0 {
		pinyinArr = append(pinyinArr, nohans)
	}
	return pinyinArr
}

func (this *Options) Convert(strs string) []string {
	retArr := make([]string, 0)
	if this.segment {
		jiebaed := jieba.Cut(strs, use_hmm)
		for _, item := range jiebaed {
			mapValuesStr, exist := phrasesDict[item]
			mapValuesArr := strings.Split(mapValuesStr, ",")
			if exist {
				for _, v := range mapValuesArr {
					retArr = append(retArr, this.perStr(v))
				}
			} else {
				converted := this.doConvert(item)
				for _, v := range converted {
					retArr = append(retArr, v)
				}
			}
		}
	} else {
		retArr = this.doConvert(strs)
	}

	return retArr
}

func (this *Options) GetCovertDict(strs string) ([]string, map[string]string) {
	length := len(strs)
	dictHan2Pin := make(map[string]string, length)

	//获取字符串的长度
	bytes := []byte(strs)
	bolcks := make([]string, 0)
	var single string
	var tmp string
	for len(bytes) > 0 {
		r, w := utf8.DecodeRune(bytes)
		bytes = bytes[w:]
		//util.LogInfo(fmt.Sprintf("[pinyin] %v", int(r)))
		single = get(int(r))
		tmp = string(r)
		// 汉字
		if len(single) != 0 {
			tmp = this.perStr(single)
			dictHan2Pin[string(r)] = tmp
		}
		blocks = append(blocks, string(r))
	}
	return blocks, dictHan2Pin

}

func (this *Options) doConvertAndGeneDict(strs string) ([]string, map[string]string, map[string]string) {
	length := len(strs)
	dictHan2Pin := make(map[string]string, length)
	dictPin2Han := make(map[string]string, length)

	//获取字符串的长度
	bytes := []byte(strs)
	pinyinArr := make([]string, 0)
	var single string
	var tmp string
	for len(bytes) > 0 {
		r, w := utf8.DecodeRune(bytes)
		bytes = bytes[w:]
		//util.LogInfo(fmt.Sprintf("[pinyin] %v", int(r)))
		single = get(int(r))
		tmp = string(r)
		// 汉字
		if len(single) != 0 {
			//util.LogInfo(fmt.Sprintf("[pinyin] %v", single))
			tmp = this.perStr(single)
			//util.LogInfo(fmt.Sprintf("[pinyin] %v", tmp))
			dictPin2Han[tmp] = string(r)
			dictHan2Pin[string(r)] = tmp
		}
		pinyinArr = append(pinyinArr, tmp)
	}
	return pinyinArr, dictHan2Pin, dictPin2Han
}

func (this *Options) ConvertAndGeneDict(strs string) ([]string, map[string]string, map[string]string) {
	return this.doConvertAndGeneDict(strs)
}

func NewPy(style int, segment bool) *Options {
	return &Options{style, segment, false}
}
