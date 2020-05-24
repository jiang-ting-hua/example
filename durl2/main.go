package main

import (
	"encoding/json"
	"fmt"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
	"log"
	"path/filepath"
	"strconv"
	"sync"
)

var (
	maxGo      int = 15        //最大并发数
	imgSize    int = 1024 * 30 //图片大于30K,才下载
	MaxLayer   int = 3         //查找网页的最大层数,用于图片下载
	waitGroup  sync.WaitGroup
	m3u8Url    string //m3u8下载的URL
	imgUrl     string //img下载的URL
	DowloadMgr Dowload
	isdowload  bool
)
//解析网页传的json数据.
type Dowload struct {
	Urltype  int    `json:"urltype"`
	MaxGo    string `json:"maxgo"`
	ImgSize  string `json:"imgsize"`
	MaxLayer string `json:"maxlayer"`
	UrlAddr  string `json:"urladdr"` //下载的URL
}

func main() {

	w, err := window.New(sciter.SW_TITLEBAR|
		sciter.SW_RESIZEABLE|
		sciter.SW_CONTROLS|
		sciter.SW_MAIN|
		sciter.SW_ENABLE_DEBUG,
		//给窗口设置个大小
		&sciter.Rect{Left: 200, Top: 100, Right: 800, Bottom: 400})
	if err != nil {
		log.Fatal(err)
	}

	fp, err := filepath.Abs("durl.html")
	if err != nil {
		log.Fatal(err)
		return
	}
	w.LoadFile(fp)
	w.SetTitle("M3U8视频和网页图片下载")
	defFunc(w)
	w.Show()
	w.Run()
}

func defFunc(w *window.Window) {
	//注册reg函数，用于处理注册逻辑
	w.DefineFunction("reg", func(args ...*sciter.Value) *sciter.Value {
		if isdowload==true{
			return sciter.NullValue()
		}
		DowloadMgr = Dowload{}
		for _, v := range args {
			err := json.Unmarshal([]byte(v.String()), &DowloadMgr)
			if err != nil {
				fmt.Println(err)
				return sciter.NullValue()
			}
			maxGo, err = strconv.Atoi(DowloadMgr.MaxGo)
			if err != nil {
				maxGo = 15
			}
			imgSize, err = strconv.Atoi(DowloadMgr.ImgSize)
			if err != nil {
				imgSize = 1024 * 30
			}
			MaxLayer, err = strconv.Atoi(DowloadMgr.MaxLayer)
			if err != nil {
				MaxLayer = 3
			}
			fmt.Println(DowloadMgr)
			switch DowloadMgr.Urltype {
			case 0:
				isdowload=true
				m3u8Url = DowloadMgr.UrlAddr
				err := dowloadM3u8(m3u8Url)
				if err != nil {
					fmt.Println(err)
					return sciter.NullValue()
				}
				isdowload=false
			case 1:
				isdowload=true
				imgUrl = DowloadMgr.UrlAddr
				DownloadImg(imgUrl)
				isdowload=false
			default:
				fmt.Println(`选择类型无效`)
				return sciter.NullValue()
			}
		}
		return sciter.NullValue()
	})
}
