package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type State struct {
	Method string
}

type Sender struct {
	State
	InputRequest, InputKey, InputDomain, Display *widget.Entry
	ScrollContainer                              *container.Scroll
	SendBtn, SaveResultBtn                       *widget.Button
	SelectMethod                                 *widget.Select
}

func main() {
	app := app.New()
	window := app.NewWindow("Crm sender")

	btnExit := widget.NewButton("Exit", func() {
		app.Quit()
	})
	sender := Sender{
		InputRequest: widget.NewEntry(),
		InputKey:     widget.NewEntry(),
		InputDomain:  widget.NewEntry(),
		Display:      widget.NewEntry(),
	}

	sender.InputRequest.SetPlaceHolder("Enter the address bar for the request")
	sender.InputKey.SetPlaceHolder("Enter key")
	sender.InputDomain.SetPlaceHolder("Enter domain")

	sender.SelectMethod = sender.GetSelectMethod()
	sender.ScrollContainer = sender.GetScrollDisplay()
	sender.SendBtn = sender.SendBtnHandler()
	sender.SaveResultBtn = sender.SaveResultBtnHandler(window)

	content := container.NewGridWithRows(3,
		container.NewGridWithColumns(
			2,
			container.NewGridWithRows(5, sender.InputRequest, sender.InputKey, sender.InputDomain, sender.SelectMethod, sender.SendBtn),
			container.NewGridWithRows(1, sender.ScrollContainer),
		),
		sender.SaveResultBtn,
		btnExit,
	)
	window.SetContent(content)
	window.CenterOnScreen()
	window.Resize(fyne.NewSize(500, 400))
	window.ShowAndRun()
}

func (sender *Sender) GetSelectMethod() *widget.Select {
	resp := widget.NewSelect([]string{"GET"}, func(value string) {
		sender.Method = value
	})
	resp.PlaceHolder = "Select method"
	return resp
}

func (sender *Sender) GetScrollDisplay() *container.Scroll {
	return container.NewVScroll(container.NewGridWithRows(
		1,
		sender.Display,
	))
}

func (sender *Sender) SendBtnHandler() *widget.Button {
	return widget.NewButton("Send", func() {
		if sender.InputRequest.Text != "" && sender.Method != "" && sender.InputKey.Text != "" && sender.InputDomain.Text != "" {
			sender.Display.SetText("")
			resp, err := sender.SendByMethod()
			if err == nil {
				body, err := io.ReadAll(resp.Body)
				if err == nil {
					defer resp.Body.Close()
					var prettyJSON bytes.Buffer
					if err := json.Indent(&prettyJSON, []byte(body), "", "    "); err == nil {
						sender.showResp(prettyJSON.String())
					} else {
						sender.showResp(err.Error())
					}
				} else {
					sender.showResp(err.Error())
				}
			} else {
				sender.showResp(err.Error())
			}
		} else {
			sender.showResp("Enter the request string")
		}
	})
}

func (sender *Sender) SendByMethod() (*http.Response, error) {
	var resp *http.Response
	var err error
	switch sender.Method {
	case "GET":
		str := sender.GetMD5Hash(sender.InputDomain.Text + strconv.FormatInt(time.Now().Unix(), 10) + sender.InputKey.Text)
		resp, err = http.Get(sender.InputRequest.Text + "?domain=" + sender.InputDomain.Text + "&time=" + strconv.FormatInt(time.Now().Unix(), 10) + "&token=" + str)
	default:
		return resp, err
	}
	return resp, err
}

func (sender *Sender) showResp(data string) {
	var strBuilder strings.Builder
	strBuilder.WriteString("[")
	strBuilder.WriteString("{")
	sender.Display.SetText(data)
	strBuilder.WriteString("}")
	strBuilder.WriteString("]")
}

func (httpSender *Sender) GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func (sender *Sender) SaveResultBtnHandler(appWindow fyne.Window) *widget.Button {
	return widget.NewButton("Save result to file", func() {
		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err == nil && writer != nil {
				_, err := writer.Write([]byte(sender.Display.Text))
				if err != nil {
					dialog.ShowError(err, appWindow)
				}
			}
		}, appWindow)
	})
}
