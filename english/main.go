package main

import (
	"bufio"
	"english/conf"
	"english/word"
	"fmt"
	"github.com/hajimehoshi/oto"
	"github.com/tosone/minimp3"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	wait        = make(chan bool)
	again       = make(chan bool)
	next        = make(chan bool)
	inputStr    string
	lastNum     int
	kernel32    *syscall.LazyDLL  = syscall.NewLazyDLL(`kernel32.dll`)
	proc        *syscall.LazyProc = kernel32.NewProc(`SetConsoleTextAttribute`)
	CloseHandle *syscall.LazyProc = kernel32.NewProc(`CloseHandle`)
	// 给字体颜色对象赋值
	FontColor Color = Color{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
)

type Color struct {
	black        int // 黑色
	blue         int // 蓝色
	green        int // 绿色
	cyan         int // 青色
	red          int // 红色
	purple       int // 紫色
	yellow       int // 黄色
	light_gray   int // 淡灰色（系统默认值）
	gray         int // 灰色
	light_blue   int // 亮蓝色
	light_green  int // 亮绿色
	light_cyan   int // 亮青色
	light_red    int // 亮红色
	light_purple int // 亮紫色
	light_yellow int // 亮黄色
	white        int // 白色
}

func (c *Color) getColor(color string) int {
	color = strings.ToLower(color)
	switch color {
	case "black":
		return 0
	case "blue":
		return 1
	case "green":
		return 2
	case "cyan":
		return 3
	case "red":
		return 4
	case "purple":
		return 5
	case "yellow":
		return 6
	case "light_gray":
		return 7
	case "gray":
		return 8
	case "light_blue":
		return 9
	case "light_green":
		return 10
	case "light_cyan":
		return 11
	case "light_red":
		return 12
	case "light_purple":
		return 13
	case "light_yellow":
		return 14
	case "white":
		return 15
	default:
		return 7
	}
}

// 输出有颜色的字体
func ColorPrint(s string, i int) {
	handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(i))
	print(s)
	CloseHandle.Call(handle)
}

func main() {
	//初始化语音文件
	err := word.InitVoice()
	if err != nil {
		fmt.Println(err)
	}
	//初始化单词文本文件内容.
	word.InitWord()

	//另启动一个协程,读取键盘输入,未直接捕获键盘输入,杀毒软件容易误判拦截.
	go readKeyboard()
	//配置文件中,lastNum上一次读到多少行的单词
	lastNum = conf.StudyConf.GetLastNum()

	fmt.Println("  r+回车---继续学习本单词")
	fmt.Println("  n+回车---学习下一个单词")
	fmt.Println("  回车键---暂停/继续播放")
	fmt.Println("  q+回车---退出")


	for {
		if len(word.Words) == 0 {
			fmt.Println("单词文本为空.或者文本保存格式错误请检查。(格式为:ANSI编码)")
			break
		}
		//读取单词
		for i, v := range word.Words {
			//配置文件中,lastNum上一次读到多少行的单词,接着上一次最后读的单词往后读.
			if i < lastNum {
				continue
			}
			lastNum = i
			color := conf.StudyConf.GetfontColor()

		againLocation: //继续学习本单词

			//终端显示单词信息
			ColorPrint(strconv.Itoa(i)+"            "+v.Word+"  ", FontColor.getColor(color))
			ColorPrint(" "+v.Chinese+"  ", FontColor.getColor(color))
			fmt.Println()
			time.Sleep(time.Second * 1)
			//获取重读次数
			repeatNum := conf.StudyConf.GetRepeatNum()
			for i := 0; i < repeatNum; i++ {
				//go的mp3解码这个mp3解不了码,没找到合适的mp3解码,暂时只能调用这个程序播放语音.执行了,再kill掉.只能暂时用.
				cmd := exec.Command("mplayer.exe", v.VoiceFilePath)
				err := cmd.Start()
				if err != nil {
					fmt.Println(err)
					return
				}
				time.Sleep(time.Second * 1)  //读完一个单词,再kill掉mplayer.exe程序.
				killEXE("mplayer.exe")
				cmd.Wait()

				//通过管道,判断键盘输入.
				select {
				case <-next:
					goto nextLocation
				default:
				}
				time.Sleep(time.Second * 2) //读一个单词,暂停一下.
			}
			time.Sleep(time.Second * 3)//读下一个单词,暂停一下.

			select {
			case <-wait:  //判断是否键盘输入:回车键暂停.
				scanner := bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					t := scanner.Text()
					if t == "q" || t == "Q" {
						killEXE("mplayer.exe")
						err := conf.SetConf("study", "last_num", strconv.Itoa(lastNum))
						if err != nil {
							fmt.Println(err)
						}
						os.Exit(0)
					}
					if t == "" {
						break
					}
				}
			case <-again: //判断是否键盘输入: r+回车,继续学习本单词
				goto againLocation
			case <-next:  //判断是否键盘输入:n+回车,学习下一个单词
				goto nextLocation
			default:
				continue
			}
		nextLocation: //跳到学习下一个单词的标签.
		}

		//每读一个单词,把读到的单词行号保存在配置文件中.
		err = conf.SetConf("study", "last_num", strconv.Itoa(0))
		if err != nil {
			fmt.Println(err)
		}
	}
}

//读取键盘输入,未直接捕获键盘输入,杀毒软件容易误判拦截.
func readKeyboard() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		t := scanner.Text()
		if t == "r" || t == "R" || t == "." {
			again <- true
		}
		if t == "n" || t == "N" || t == "0" {
			next <- true
		}
		if t == "q" || t == "Q" {
			killEXE("mplayer.exe")
			err := conf.SetConf("study", "last_num", strconv.Itoa(lastNum))
			if err != nil {
				fmt.Println(err)
			}
			os.Exit(0)
		}
		if t == "" {
			killEXE("mplayer.exe")
			fmt.Println("请输入回车继续(q键加回车退出)...")
			fmt.Println()
			wait <- true
		}
	}
}

func run(fileName string) {
	var file, _ = ioutil.ReadFile(fileName)
	dec, data, _ := minimp3.DecodeFull(file)
	player, _ := oto.NewPlayer(dec.SampleRate, dec.Channels, 2, 1024)
	player.Write(data)
}

func killEXE(proName string) bool {
	arg := []string{"/im", proName, "/f"}
	cmd := exec.Command("taskkill", arg...)
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}
