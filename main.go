package main

import (
	"crmSender/customTheme"
	"crmSender/sender"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	app := app.New()
	app.Settings().SetTheme(customTheme.NewCustomTheme())
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
	sender.ClearParametersBtn = sender.ClearParametersBtnHandler()

	content := container.NewGridWithRows(
		6,
		container.NewGridWithColumns(1,sender.InputRequest),
		container.NewGridWithColumns(2, sender.InputKey, sender.InputDomain),
		container.NewGridWithColumns(1,sender.Params),
		container.NewGridWithColumns(3, sender.SelectMethod, sender.SendBtn, sender.ClearParametersBtn),
		container.NewGridWithColumns(1,sender.ErrDisplay),
		btnExit,
	)
	window.SetContent(content)
	window.CenterOnScreen()
	window.Resize(fyne.NewSize(500, 400))
	window.ShowAndRun()
}
