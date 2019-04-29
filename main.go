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
	driver.SetDriver() // 自动获取浏览器驱动
	chromeDriver := webdriver.NewChromeDriver(fmt.Sprintf("%s/driver/%s", cf.Environ.Root, cf.Environ.DriverName))
	browser := page.NewSession(chromeDriver)
	browser.OpenBrowser()
	lg := page.NewLoginPage(browser)
	lg.Login()
	//lg.SendCode()
	reply := make(map[string]bool)
	msgList := make(map[string]map[string]string)
	msg := page.NewMessagePage(msgList, reply, browser)
	msg.Listen()
	defer page.TearDown(browser)
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
