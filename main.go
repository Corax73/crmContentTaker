package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type State struct {
	Method string
}

type Sender struct {
	State
	InputRequest, InputKey, InputDomain, Display *widget.Entry
	ScrollContainer                              *container.Scroll
	SendBtn                                      *widget.Button
	SelectMethod                                 *widget.Select
}

func main() {
	app := app.New()

	btnExit := widget.NewButton("Exit", func() {
		app.Quit()
	})
	sender := Sender{
		InputRequest: widget.NewEntry(),
		InputKey:     widget.NewEntry(),
		InputDomain:  widget.NewEntry(),
		Display:      widget.NewEntry(),
	}

	sender.SelectMethod = sender.GetSelectMethod()
	sender.ScrollContainer = sender.GetScrollDisplay()
	sender.SendBtn = sender.SendBtnHandler()

	content := container.NewGridWithRows(2,
		container.NewGridWithColumns(
			2,
			container.NewGridWithRows(5, sender.InputRequest, sender.InputKey, sender.InputDomain, sender.SelectMethod, sender.SendBtn),
			container.NewGridWithRows(1, sender.ScrollContainer),
		),
		btnExit,
	)
	window := app.NewWindow("Crm sender")
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
		if sender.InputRequest.Text != "" && sender.Method != "" {
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
		resp, err = http.Get(sender.InputRequest.Text)
	default:
		return resp, err
	}
	return resp, err
}

func (httpSender *Sender) showResp(data string) {
	var strBuilder strings.Builder
	strBuilder.WriteString("[")
	strBuilder.WriteString("{")
	httpSender.Display.SetText(data)
	strBuilder.WriteString("}")
	strBuilder.WriteString("]")
}
