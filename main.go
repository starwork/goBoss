package main

import (
	"fmt"
	cf "goBoss/config"
	"goBoss/page"
	"log"
	"os"
	"goBoss/driver"
	"github.com/fedesog/webdriver"
)

func main() {
	setLog()
	driver.SetDriver()  // 自动获取浏览器驱动
	chromeDriver := webdriver.NewChromeDriver(fmt.Sprintf("%s/driver/%s", cf.Environ.Root, cf.Environ.DriverName))
	engine := &page.Engine{Dr: chromeDriver}
	lg := &page.Login{Eg: engine}
	lg.Eg.Start()
	lg.Eg.OpenBrowser()
	lg.Login()
	reply := make(map[string]bool)
	msgList := make(map[string]map[string]string)
	msg := &page.Message{
		Eg: engine,
		ReplyList: reply, MsgList: msgList,
	}
	msg.Listen()
	defer page.TearDown(engine)
}

func setLog() {
	//set logfile Stdout
	logFile, logErr := os.OpenFile(fmt.Sprintf("%s/boss.log", cf.Environ.Root), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		fmt.Println("Fail to find", logFile, "cServer start Failed")
		os.Exit(1)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
