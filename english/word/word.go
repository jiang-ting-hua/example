package word

import (
	"bufio"
	"english/conf"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	Words  = []*wordInfo{}
	Voices = map[string]*voice{}
)
//单词信息
type wordInfo struct {
	Word    string
	Chinese string
	VoiceFilePath string
}
//语音信息
type voice struct{
	FileName string
	Suffix string
	FilePath string
}
//获取当前目录下的文件
func getDirFile(dir string) (files []string,err error){
	f, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil,err
	}
	for _, file := range f {
		if !file.IsDir() {
			fileName := file.Name()
			files = append(files,fileName)
		}
	}
	return files,nil
}
//初始化语音文件
func InitVoice()(err error) {
	// 读取当前目录中的所有文件和子目录
	voicePath := conf.VoiceConf.GetFilePath()
	f, err := ioutil.ReadDir(voicePath)
	if err != nil {
		return err
	}
	// 获取文件，并输出它们的名字
	for _, file := range f {
		//当前目录下的文件
		if !file.IsDir() {
			fileName := file.Name()
			i := strings.LastIndex(fileName,".")
			if i == -1 {
				continue
			}
			name := fileName[:i]
			name = strings.ToLower(name)
			suffix := fileName[i:]
			Voices[fileName]=&voice{
				FileName: name,
				Suffix:   suffix,
				FilePath: voicePath+"/"+fileName,
			}
		}
		//当前目录下的目录
		if file.IsDir() {
			dirName := file.Name()
			dirPath := voicePath + "/" + dirName
			files,err := getDirFile(dirPath)
			if err != nil {
				return err
			}
			for _,v := range files{
				i := strings.LastIndex(v,".")
				if i == -1 {
					continue
				}
				name := v[:i]
				name = strings.ToLower(name)
				suffix := v[i:]
				Voices[name]=&voice{
					FileName: name,
					Suffix:   suffix,
					FilePath: dirPath+"/"+v,
				}
			}
		}
	}
	return nil
}
//初始化单词文本文件。
func InitWord() (err error) {
	f, err := os.OpenFile(conf.WordConf.GetFullPath(), os.O_RDONLY, 0)
	if err != nil {
		err = fmt.Errorf("open file(%s)failed:%w\n", conf.WordConf.GetFullPath(), err)
		return err
	}
	defer f.Close()
	fileScanner := bufio.NewScanner(f)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		line = strings.TrimSpace(line)
		i := strings.Index(line, ".")
		if i == -1 {
			continue
		}
		//n := line[:i]
		temp := line[i+1:]
		temp = strings.TrimSpace(temp)
		j := strings.Index(temp, " ")
		if j == -1 {
			continue
		}
		word := temp[:j]
		word = strings.TrimSpace(word)
		word = strings.ToLower(word)
		chinese := temp[j+1:]
		chinese = strings.TrimSpace(chinese)
		var voiceFilePath string
		if v, ok := Voices[word]; ok {
			voiceFilePath = v.FilePath
		}
		w := wordInfo{
			Word:    word,
			Chinese: chinese,
			VoiceFilePath:voiceFilePath,
		}
		Words = append(Words,&w)
	}
return nil
}
