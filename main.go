package main

import (
	"fmt"

	"github.com/ques0942/sampleapp/model"
)

func hyoushow(receiver string, sender string, dateStr string, body string) (string, error) {
	pageModel, err := model.HyoushowPageOpen()
	if err != nil {
		return "", err
	}
	defer pageModel.Close()

	if err := pageModel.SetReceiver(receiver); err != nil {
		return "", err
	}
	if err := pageModel.SetSender(sender); err != nil {
		return "", err
	}
	if err := pageModel.SetDateStr(dateStr); err != nil {
		return "", err
	}
	if err := pageModel.SetCertificateBody(body); err != nil {
		return "", err
	}
	return pageModel.Generate()
}

func main() {
	resultURL, err := hyoushow("testさん", "from test", "12/01", "試しに\nあなたを\n表彰します")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resultURL)
}
