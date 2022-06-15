package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"

	iconv "github.com/djimenez/iconv-go"
	"github.com/go-vgo/robotgo"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	. "github.com/lxn/walk/declarative"
	"github.com/micmonay/keybd_event"
)

type sMainWindow struct {
	*walk.MainWindow
}

var protocol = []string{"PPTP", "OpenVPN"}

func main() {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}
	kb.SetKeys(keybd_event.VK_ESC)

	quit := make(chan bool, 1)

	mw := new(sMainWindow)
	var outTE *walk.TextEdit
	var protocolCB *walk.ComboBox

	MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "VPN Checker",
		Size:     declarative.Size{Width: 120, Height: 240},
		Layout:   VBox{},
		Children: []Widget{
			TextEdit{AssignTo: &outTE, ReadOnly: true},
			ComboBox{
				AssignTo: &protocolCB,
				Editable: false,
				Model:    protocol,
			},
			PushButton{
				MaxSize: Size{Width: 100, Height: 50},
				MinSize: Size{Width: 100, Height: 50},
				Text:    "시작",
				OnClicked: func() {
					if protocolCB.CurrentIndex() == -1 {
						protocolCB.SetCurrentIndex(1)
					}
					//fmt.Println(runtime.NumGoroutine())
					if runtime.NumGoroutine() < 3 {
						if protocolCB.CurrentIndex() == 0 {
							go pptpCheck(quit, quit, outTE, kb)
						} else {
							go openVPNCheck(quit, quit, outTE, "연결 끊김", kb)
						}
					}
				},
			},
			PushButton{
				MaxSize: Size{Width: 100, Height: 50},
				MinSize: Size{Width: 100, Height: 50},
				Text:    "중지",
				OnClicked: func() {
					go func(ch chan bool) {
						if safeCheck(ch) {
							quit <- true
						}
					}(quit)
				},
			},
		},
		Bounds: Rectangle{X: 1289, Y: 803},
	}.Run()

}

func (mw *sMainWindow) showMessageBox(msg string) {
	walk.MsgBox(mw,
		"Message",
		msg,
		walk.MsgBoxOK|walk.MsgBoxIconInformation)
}

func (mw *sMainWindow) showMessageError(msg string) {
	walk.MsgBox(mw,
		"Error",
		msg,
		walk.MsgBoxOK|walk.MsgBoxIconError)
}

func safeCheck(ch <-chan bool) bool {
	select {
	case <-ch:
		return false
	default:
	}
	return true
}

func pptpCheck(
	recvCh <-chan bool,
	sendCh chan<- bool,
	logText *walk.TextEdit,
	keybd keybd_event.KeyBonding) {

	for {
		select {
		case <-recvCh:
			logText.SetText("중지 되었습니다.")
			return
		default:
			cmd := exec.Command("ipconfig", "/all")
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			out, err := cmd.Output()
			if err != nil {
				fmt.Println(err)
				return
			}
			result, _ := iconv.ConvertString(string(out), "euc-kr", "utf-8")
			fmt.Println(result)

			if strings.Contains(result, "PPP") {
				logText.SetText("vpnChecking...")
				fmt.Println("vpnChecking...")
			} else {
				logText.SetText("중지 되었습니다.")
				robotgo.MoveClick(1348, 887, "left", true)
				robotgo.MoveClick(1348, 887, "left", true)
				robotgo.MoveClick(1348, 887, "left", true)
				robotgo.MoveClick(1348, 887, "left", true)
				keybd.Launching()

				sendCh <- true
			}
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func openVPNCheck(
	recvCh <-chan bool,
	sendCh chan<- bool,
	logText *walk.TextEdit,
	checkType string,
	keybd keybd_event.KeyBonding) {

	for {
		select {
		case <-recvCh:
			logText.SetText("중지 되었습니다.")
			return
		default:
			out, err := exec.Command("ipconfig", "/all").Output()
			if err != nil {
				fmt.Println(err)
				return
			}
			result, _ := iconv.ConvertString(string(out), "euc-kr", "utf-8")
			fmt.Println(result)

			if !strings.Contains(result, checkType) {
				logText.SetText("vpnChecking...")
				fmt.Println("vpnChecking...")
			} else {
				logText.SetText("중지 되었습니다.")
				robotgo.MoveClick(1348, 887, "left", true)
				robotgo.MoveClick(1348, 887, "left", true)
				robotgo.MoveClick(1348, 887, "left", true)
				robotgo.MoveClick(1348, 887, "left", true)
				keybd.Launching()

				sendCh <- true
			}
		}
		time.Sleep(1000 * time.Millisecond)
	}
}
