package page

import (
	"github.com/fedesog/webdriver"
	"goBoss/utils"
	"fmt"
	"log"
	"errors"
	cf "goBoss/config"
	"time"
	"os"
	"encoding/base64"
)

const defaultUrl = "127.0.0.1"

type session struct {
	d *webdriver.ChromeDriver
	*webdriver.Session
}

func NewSession(d *webdriver.ChromeDriver) *session {
	s := start(d)
	if s == nil {
		return nil
	}
	return &session{d, s}
}

func start(d *webdriver.ChromeDriver) *webdriver.Session {
	err := d.Start()
	if err != nil {
		log.Printf("启动浏览器驱动失败: %s", err.Error())
		return nil
	}
	args := make([]string, 0)
	if cf.Config.Headless {
		args = append(args, "--headless")
	}
	desired := webdriver.Capabilities{
		"Platform":           "Mac",
		"goog:chromeOptions": map[string][]string{"args": args, "extensions": {}, "excludeSwitches": {"enable-automation"}},
		"browserName":        "chrome",
		"version":            "",
		"platform":           "ANY",
	}
	required := webdriver.Capabilities{}
	se, err := d.NewSession(desired, required)
	if err != nil {
		log.Printf("open browser failed: %s", err.Error())
		return nil
	}
	return se
}

type Driver interface {
	MaxWindow() error
	OpenBrowser() error
	SetWindow(width, height int) error
	GetElement(root, name string) *Element
	Screen() ([]byte, error)
	ScreenShot(key string) string
	ScreenAsBs64() string
}

func (w *session) GetElement(root, name string) *Element {
	r, ok := Page[root]
	if !ok || r == nil {
		log.Panicf("page/element.json未找到页面: [%s]", root)
	}
	ele, ok := Page[root][name]
	if !ok {
		log.Panicf("page/element.json未找到root: [%s] key: [%s]", root, name)
	}
	return &ele
}

func (w *session) MaxWindow() error {
	p := fmt.Sprintf(`{"windowHandle": "current", "sessionId": "%s"}`, w.Id)
	req := utils.RequestData{
		JSON: p,
	}
	url := fmt.Sprintf("http://%s:%d/session/%s/window/current/maximize", defaultUrl, w.d.Port, w.Id)
	_, err := utils.HttpPost(url, req)
	if err != nil {
		log.Printf("最大化窗口失败, error: %v", err)
		return err
	}
	return nil
}

func (w *session) OpenBrowser() error {
	w.Url(cf.Config.LoginURL)
	err := w.SetWindow(1600, 900)
	if err != nil {
		log.Panicf("最大化浏览器失败!!!error: %s", err.Error())
	}
	err = w.SetTimeoutsImplicitWait(cf.Config.WebTimeout)
	if err != nil {
		log.Printf("设置脚本超时时间失败, error: %v\n", err)
		return err
	}
	return nil
}

func (w *session) SetWindow(width, height int) error {
	p := fmt.Sprintf(`{"windowHandle": "current", "sessionId": "%s", "height": %d, "width": %d}`, w.Id, height, width)
	url := fmt.Sprintf("http://127.0.0.1:%d/session/%s/window/current/size", w.d.Port, w.Id)
	res, err := utils.HttpPost(url, utils.RequestData{JSON: p})
	if err != nil {
		return errors.New(fmt.Sprintf(`设置浏览器窗口失败, 请检查!%+v`, res.Error))
	}
	return nil
}

func (w *session) Close() error {
	err := w.CloseCurrentWindow()
	if err != nil {
		return err
	}
	return w.d.Stop()
}

func (w *session) Screen() ([]byte, error) {
	return w.Screenshot()
}

func (w *session) ScreenShot(key string) string {
	pic, _ := w.Screen()
	filename := fmt.Sprintf("%s_%s.png", key, time.Now().Format("2006_01_02_15_04_05"))
	filename = fmt.Sprintf("%s/picture/%s", cf.Environ.Root, filename)
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Printf("发送消息后截图失败!Error: %s", err.Error())
	}
	f.Write(pic)
	defer f.Close()
	return filename
}

func (w *session) ScreenAsBs64() string {
	bt, err := w.Screen()
	if err != nil {
		log.Println("截图出错!Error: ", err.Error())
	}
	bs64 := utils.Encode(base64.StdEncoding, bt)
	return bs64
}
