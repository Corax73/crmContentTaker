package sender

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
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type State struct {
	Method, Data string
	Window       *fyne.Window
}

type Sender struct {
	State
	InputRequest, InputKey, InputDomain, ErrDisplay, Params *widget.Entry
	SendBtn                                              *widget.Button
	SelectMethod                                         *widget.Select
}

func (sender *Sender) GetSelectMethod() *widget.Select {
	resp := widget.NewSelect([]string{"GET"}, func(value string) {
		sender.Method = value
	})
	resp.PlaceHolder = "Select method"
	return resp
}

func (sender *Sender) SendBtnHandler() *widget.Button {
	return widget.NewButton("Send", func() {
		if sender.InputRequest.Text != "" && sender.Method != "" && sender.InputKey.Text != "" && sender.InputDomain.Text != "" {
			sender.Data = ""
			resp, err := sender.SendByMethod()
			if err == nil {
				body, err := io.ReadAll(resp.Body)
				if err == nil {
					defer resp.Body.Close()
					var prettyJSON bytes.Buffer
					if err := json.Indent(&prettyJSON, []byte(body), "", "    "); err == nil {
						sender.Data = prettyJSON.String()
						sender.RunSavingDialog(*sender.Window)
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
		url := sender.InputRequest.Text + "?domain=" + sender.InputDomain.Text + "&time=" + strconv.FormatInt(time.Now().Unix(), 10) + "&token=" + str
		params := sender.GetParams()
		strParams := sender.MakeGetParams(params)
		resp, err = http.Get(url + "&" + strParams)
	default:
		return resp, err
	}
	return resp, err
}

func (sender *Sender) showResp(data string) {
	var strBuilder strings.Builder
	strBuilder.WriteString("[")
	strBuilder.WriteString("{")
	sender.ErrDisplay.SetText(data)
	strBuilder.WriteString("}")
	strBuilder.WriteString("]")
}

func (httpSender *Sender) GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func (sender *Sender) RunSavingDialog(appWindow fyne.Window) {
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err == nil && writer != nil {
			_, err := writer.Write([]byte(sender.Data))
			if err != nil {
				dialog.ShowError(err, appWindow)
			}
		}
	}, appWindow)
}

func (sender *Sender) GetParams() *bytes.Buffer {
	data := make(map[string]interface{})
	str := sender.Params.Text
	if str == "" {
		str = "{}"
	}
	err := json.Unmarshal([]byte(str), &data)
	if err != nil {
		sender.showResp(err.Error())
	}
	postBody, _ := json.Marshal(data)
	return bytes.NewBuffer(postBody)
}

func (sender *Sender) MakeGetParams(data *bytes.Buffer) string {
	strParams := data.String()
	strParams = strings.Trim(strParams, "{")
	strParams = strings.Trim(strParams, "}")
	strParams = strings.Replace(strParams, "\"", "", -1)
	strParams = strings.Replace(strParams, ":", "=", -1)
	strParams = strings.Replace(strParams, ",", "&", -1)
	return strParams
}
