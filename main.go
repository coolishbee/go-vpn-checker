package main

import (
	"fmt"
	"os/exec"
	"strings"
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
var onOff = false

func main() {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}
	kb.SetKeys(keybd_event.VK_ESC)

	quit := make(chan bool)
	//sendCh := make(chan bool)
	//recvCh := make(chan bool)

	mw := new(sMainWindow)
	var outTE *walk.TextEdit
	var protocolType *walk.ComboBox

	MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "VPN Checker",
		Size:     declarative.Size{Width: 120, Height: 240},
		Layout:   VBox{},
		Children: []Widget{

			TextEdit{AssignTo: &outTE, ReadOnly: true},
			ComboBox{
				AssignTo: &protocolType,
				Editable: false,
				Model:    protocol,
			},
			PushButton{
				MaxSize: Size{Width: 100, Height: 50},
				MinSize: Size{Width: 100, Height: 50},
				Text:    "시작",
				OnClicked: func() {
					if !onOff {
						onOff = true
						if protocolType.CurrentIndex() == 0 {
							go pptpCheck(quit, quit, outTE, kb)
						} else {
							go openVPNCheck(quit, quit, outTE, "연결 끊김", kb)
						}
					} else {
						outTE.SetText("이미 실행중입니다")
					}
				},
			},
			PushButton{
				MaxSize: Size{Width: 100, Height: 50},
				MinSize: Size{Width: 100, Height: 50},
				Text:    "중지",
				OnClicked: func() {
					if onOff {
						onOff = false
						quit <- true
					} else {
						outTE.SetText("이미 중지되었습니다")
					}
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

func pptpCheck(
	recvCh <-chan bool,
	sendCh chan<- bool,
	logText *walk.TextEdit,
	keybd keybd_event.KeyBonding) {

	for {
		select {
		case <-recvCh:
			logText.SetText("중지 되었습니다.")
			onOff = false
			return
		default:
			out, err := exec.Command("ipconfig", "/all").Output()
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
