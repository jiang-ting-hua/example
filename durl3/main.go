package main

import (
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"sync"
)

var (
	maxGo         int = 15        //最大并发数
	imgSize       int = 1024 * 30 //图片大于30K,才下载
	MaxLayer      int = 3         //查找网页的最大层数,用于图片下载
	waitGroup     sync.WaitGroup
	m3u8Url       string //m3u8下载的URL
	imgUrl        string //img下载的URL
	isdowload     bool
	Window        *ui.Window
	Processbar    *ui.ProgressBar
	Greeting      *ui.Label
	ProcessbarNum int
	isClose       bool
	mux           sync.Mutex = sync.Mutex{}
)

func main() {
	err := ui.Main(windowInit)
	if err != nil {
		panic(err)
	}
}

//窗口初始化
func windowInit() {
	//生成:单选
	radioLabel := ui.NewLabel(`请选择要下载的类型：`)
	radio := ui.NewRadioButtons()
	radio.Append("M3U8视频")
	radio.Append("网页图片")
	radio.SetSelected(0)
	// 生成：maxgo
	maxGoLabel := ui.NewLabel(`并发下载数量(默认15)：`)
	maxGoSpinbox := ui.NewSpinbox(1, 100)
	maxGoSpinbox.SetValue(15)
	// 生成：imgSize
	imgSizeLabel := ui.NewLabel(`图片下载大小(默认大于30KB)：`)
	imgSizeSpinbox := ui.NewSpinbox(1, 102400)
	imgSizeSpinbox.SetValue(30)
	// 生成：maxLayer
	maxLayerLabel := ui.NewLabel(`图片下载网页深度(默认向下3层)：`)
	maxLayerSpinbox := ui.NewSpinbox(1, 100)
	maxLayerSpinbox.SetValue(3)
	// 生成：urlAddr
	urlAddrTextbox := ui.NewEntry()
	// 生成：按钮
	button := ui.NewButton(`开始下载`)
	//生成：标签
	Greeting = ui.NewLabel(``)
	// 生成：进度条
	Processbar = ui.NewProgressBar()
	//Processbar.SetValue(50)

	// 设置：按钮点击事件
	button.OnClicked(func(*ui.Button) {
		url := urlAddrTextbox.Text()
		if len(url) == 0 {
			ui.MsgBoxError(Window, "提示", "请输入网址.")
			return
		}

		switch radio.Selected() {
		case 0:

			button.Disable()
			m3u8Url = url
			go func() {
				err := DowloadM3u8(m3u8Url)
				if err != nil {
					button.Enable()
					Processbar.SetValue(0)
					ui.MsgBoxError(Window, "提示", err.Error())
					return
				}
				button.Enable()
				Processbar.SetValue(0)
				ui.MsgBox(Window, "提示", "下载完成.")
			}()
		case 1:
			button.Disable()
			imgUrl = url
			go func() {
				DownloadImg(imgUrl)
				button.Enable()
				ui.MsgBox(Window, "提示", "下载完成.")
			}()
		default:
			ui.MsgBoxError(Window, "提示", "选择类型无效.")
			return
		}
	})
	container1 := ui.NewGroup("")
	container1.SetChild(radioLabel)
	container2 := ui.NewGroup("")
	container2.SetChild(radio)

	container3 := ui.NewGroup("")
	container3.SetChild(maxGoLabel)
	container4 := ui.NewGroup("")
	container4.SetChild(maxGoSpinbox)

	container5 := ui.NewGroup("")
	container5.SetChild(imgSizeLabel)
	container6 := ui.NewGroup("")
	container6.SetChild(imgSizeSpinbox)

	container7 := ui.NewGroup("")
	container7.SetChild(maxLayerLabel)
	container8 := ui.NewGroup("")
	container8.SetChild(maxLayerSpinbox)

	//------水平排列的容器
	boxs_1 := ui.NewHorizontalBox()
	boxs_1.Append(container1, true)
	boxs_1.Append(container2, true)

	boxs_2 := ui.NewHorizontalBox()
	boxs_2.Append(container3, true)
	boxs_2.Append(container4, true)

	boxs_3 := ui.NewHorizontalBox()
	boxs_3.Append(container5, true)
	boxs_3.Append(container6, true)

	boxs_4 := ui.NewHorizontalBox()
	boxs_4.Append(container7, true)
	boxs_4.Append(container8, true)
	// 生成：垂直容器
	box := ui.NewVerticalBox()

	// 往 垂直容器 中添加 控件

	box.Append(boxs_1, true)
	box.Append(boxs_2, true)
	box.Append(boxs_3, true)
	box.Append(boxs_4, true)
	//box.Append(radio,false)
	//box.Append(ui.NewLabel(`并发下载数量(默认15)：`), false)
	//box.Append(maxGoSpinbox, false)
	//box.Append(ui.NewLabel(`图片下载大小(默认大于30KB)：`), false)
	//box.Append(imgSizeSpinbox, false)
	//box.Append(ui.NewLabel(`图片下载网页深度(默认向下3层)：`), false)
	//box.Append(maxLayerSpinbox, false)
	box.Append(ui.NewLabel(`下载网址：`), false)
	box.Append(urlAddrTextbox, false)
	box.Append(button, false)
	box.Append(Greeting, false)
	box.Append(Processbar, false)

	// 生成：窗口（标题，宽度，高度，是否有 菜单 控件）
	Window = ui.NewWindow(`M3U8视频和网页图片下载`, 400, 300, false)

	// 窗口容器绑定
	Window.SetChild(box)

	// 设置：窗口关闭时
	Window.OnClosing(func(*ui.Window) bool {
		// 窗体关闭
		isClose = true
		ui.Quit()
		return true
	})

	// 窗体显示
	Window.Show()
}
func addProcessbar() {
	mux.Lock()
	if Processbar.Value() < 90 {
		ProcessbarNum++
		Processbar.SetValue(ProcessbarNum)
	}
	mux.Unlock()
}
func setGreeting(str string) {
	mux.Lock()
	Greeting.SetText(str)
	mux.Unlock()
}
