package main

import (
	"crmSender/sender"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	app := app.New()
	window := app.NewWindow("Crm sender")

	btnExit := widget.NewButton("Exit", func() {
		app.Quit()
	})
	sender := sender.Sender{
		InputRequest: widget.NewEntry(),
		InputKey:     widget.NewEntry(),
		InputDomain:  widget.NewEntry(),
		ErrDisplay:   widget.NewEntry(),
		Params:       widget.NewEntry(),
	}

	sender.Window = &window
	sender.InputRequest.SetPlaceHolder("Enter the address bar for the request")
	sender.InputKey.SetPlaceHolder("Enter key")
	sender.InputDomain.SetPlaceHolder("Enter domain")
	sender.Params.MultiLine = true
	sender.Params.SetPlaceHolder("Enter parameters by JSON")
	sender.ErrDisplay.MultiLine = true
	sender.ErrDisplay.SetPlaceHolder("Error output")

	sender.SelectMethod = sender.GetSelectMethod()
	sender.SendBtn = sender.SendBtnHandler()

	content := container.NewGridWithRows(2,
		container.NewGridWithColumns(
			3,
			container.NewGridWithRows(4, sender.InputRequest, sender.InputKey, sender.InputDomain, sender.SelectMethod),
			container.NewGridWithRows(2, sender.Params, sender.SendBtn),
			container.NewGridWithRows(1, sender.ErrDisplay),
		),
		btnExit,
	)
	window.SetContent(content)
	window.CenterOnScreen()
	window.Resize(fyne.NewSize(500, 400))
	window.ShowAndRun()
}
