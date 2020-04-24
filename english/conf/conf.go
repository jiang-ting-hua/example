package conf

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	ConfigFile string = `./config.ini` //配置文件
)
//单词表文件配置
type wordConf struct {
	FilePath string `conf:"file_path"`
	FileName string `conf:"file_name"`
}
//语音文件配置
type voiceConf struct {
	FilePath string `conf:"file_path"`
	FileType string `conf:"file_type"`
}
//学习配置
type study struct {
	LastNum int `conf:"last_num"`
	FontColor string `conf:"font_color"`
	Repeat int `conf:"repeat"`
}

var (
	WordConf  = wordConf{}
	VoiceConf = voiceConf{}
	StudyConf = study{}
)

func (w *wordConf) GetFilePath() string {
	return w.FilePath
}
func (w *wordConf) GetFileName() string {
	return w.FileName
}
func (w *wordConf) GetFullPath() string {
	return w.FilePath + w.FileName
}
func (v *voiceConf) GetFilePath() string {
	return v.FilePath
}
func (v *voiceConf) GetFileType() string {
	return v.FileType
}
func (s *study) GetLastNum() int {
	return s.LastNum
}
func (s *study) GetfontColor() string {
	return s.FontColor
}
func (s *study) GetRepeatNum() int {
	return s.Repeat
}
func init() {
	//读取数据库登录配置文件,保存于结构conf.DbConfig{}中.
	err := ParseConf(ConfigFile, &WordConf)
	if err != nil {
		fmt.Printf("init(),读取配置文件失败: %s/n", err)
		return
	}
	err = ParseConf(ConfigFile, &VoiceConf)
	if err != nil {
		fmt.Printf("init(),读取配置文件失败: %s/n", err)
		return
	}
	err = ParseConf(ConfigFile, &StudyConf)
	if err != nil {
		fmt.Printf("init(),读取配置文件失败: %s/n", err)
		return
	}
}
//设置配置文件.
func SetConf(groupName string, k string, v string) (err error) {
	//1.打开文件
	var index int = 0
	f, err := os.OpenFile(ConfigFile, os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("打开配置文件%s失败:%w/n", ConfigFile, err)
		return err
	}
	defer f.Close()
	var newConf string
	group := ""
	fileScanner := bufio.NewScanner(f)
	for fileScanner.Scan() {
		index++
		line := fileScanner.Text()
		// 以#或;开头视为注释,空行和注释不读取
		if line == "" {
			newConf = newConf + line + "\r\n"
			continue
		}
		if strings.HasPrefix(line, "#") {
			newConf = newConf + line + "\r\n"
			continue
		}
		if strings.HasPrefix(line, ";") {
			newConf = newConf + line + "\r\n"
			continue
		}
		//line = strings.TrimSpace(line)
		//检查是否前缀是[,后缀是]的分组,并取出group组名称.
		if len(line) > 2 && line[0:1] == "[" && line[len(line)-1:] == "]" {
			group = line[1 : len(line)-1]
			group = strings.TrimSpace(group)
			newConf = newConf + line + "\r\n"
			continue
		}
		group = strings.ToUpper(group)
		groupName := strings.ToUpper(groupName)
		if group != groupName {
			newConf = newConf + line + "\r\n"
			continue
		}
		//判断是不是具体配置项,判断是不是有等号.
		index := strings.Index(line, "=")
		if index == -1 {
			newConf = newConf + line + "\r\n"
			continue
		}
		//按照等号=分割,左边是KEY,右边是VALUE
		key := line[:index]
		value := line[index+1:]
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if len(key) == 0 {
			newConf = newConf + line + "\r\n"
			continue
		}
		if key == k{
			newConf = newConf + k+" = "+v + "\r\n"
			continue
		}else{
			newConf = newConf + line + "\r\n"
		}
	}
	f.Seek(0,0)
	_, err = f.Write([]byte(newConf))
	if err!=nil {
		return err
	}
	return nil
}

//从配置文件中,取得数据,返回给一个结构体.
func ParseConf(fileName string, result interface{}) (err error) {
	t := reflect.TypeOf(result)
	v := reflect.ValueOf(result)
	//Elem()是获取引用类型指针的指向的对象,如果是值类型不需要。
	tElem := t.Elem()
	vElem := v.Elem()
	//result必须是一个指针
	if t.Kind() != reflect.Ptr {
		fmt.Println("result必须是一个指针")
		return
	}
	//result必须是一个结构体,并且结构名与配置文件中分段[]名要一样.
	if tElem.Kind() != reflect.Struct {
		fmt.Println("result必须是一个结构体")
		return
	}

	//1.打开文件
	var index int = 0

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		fmt.Printf("打开配置文件%s失败:%w\n", fileName, err)
		return err
	}
	defer f.Close()
	fileScanner := bufio.NewScanner(f)
	group := ""
	//2.将读取的文件数据按照行读取
	for fileScanner.Scan() {
		index++
		line := fileScanner.Text()
		//去除字符串首尾的空白
		line = strings.TrimSpace(line)
		// 以#或;开头视为注释,空行和注释不读取
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, ";") {
			continue
		}
		//检查是否前缀是[,后缀是]的分组,并取出group组名称.
		if len(line) > 2 && line[0:1] == "[" && line[len(line)-1:] == "]" {
			group = line[1 : len(line)-1]
			group = strings.TrimSpace(group)
			continue
		}
		//判断与传进来的结构体名称相等.
		//fmt.Println(tElem.Name())
		group = strings.ToUpper(group)
		tElemName := strings.ToUpper(tElem.Name())
		if group != tElemName {
			continue
		}
		//判断是不是具体配置项,判断是不是有等号.
		index := strings.Index(line, "=")
		if index == -1 {
			fmt.Printf("第%行语法错误:%w\n", index, err)
			return
		}
		//按照等号=分割,左边是KEY,右边是VALUE
		key := line[:index]
		value := line[index+1:]
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if len(key) == 0 {
			fmt.Printf("第%行语法错误:%w\n", index, err)
			return
		}
		//利用反射给result赋值
		//遍历结构体的每一个字段和KEY比较,匹配上就赋值.
		for i := 0; i < tElem.NumField(); i++ {
			field := tElem.Field(i)      //取得结构体字段
			tag := field.Tag.Get("conf") //到得该字段的Tag
			//如果配置文件中的Key等于该结构体字段的Tag,就把value值赋给结构体对应字段.
			if key == tag {
				fieldType := field.Type // 拿到每个字段的类型
				//根据字段的类型,对应赋值
				switch fieldType.Kind() {
				case reflect.String:
					//根据(reflect.ValueOf)中用字段名找到对应的值对象.
					//fieldValue := vElem.FieldByName(field.Name)
					////将配置文件中的value值,赋值给对应的结构体字段
					//fieldValue.SetString(value)

					//以上也可这样赋值
					vElem.Field(i).SetString(value)
				case reflect.Int64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
					valueInt64, _ := strconv.ParseInt(value, 10, 64)
					vElem.Field(i).SetInt(valueInt64)
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					valueUint64, _ := strconv.ParseUint(value, 10, 64)
					vElem.Field(i).SetUint(valueUint64)
				case reflect.Float32, reflect.Float64:
					valueFloat64, _ := strconv.ParseFloat(value, 64)
					vElem.Field(i).SetFloat(valueFloat64)
				}
			}
		}
	}
	return
}
