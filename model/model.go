package model

import (
	"errors"
	"log"
	"time"

	"github.com/sclevine/agouti"
)

/*
HyoushowPage Hyoshowのページモデル
*/
type HyoushowPage interface {
	SetReceiver(receiver string) error
	SetSender(sender string) error
	SetDateStr(dateStr string) error
	SetCertificateBody(certificateBody string) error
	Generate() (string, error)
	Close() error
}

type hyoushowPage struct {
	driver *agouti.WebDriver
	page   *agouti.Page
}

/*
HyoushowPageOpen hyou.showのページモデルを生成する
*/
func HyoushowPageOpen() (HyoushowPage, error) {
	h := &hyoushowPage{}
	h.driver = agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			"--headless",
		}),
	)
	if err := h.driver.Start(); err != nil {
		return nil, err
	}
	if page, err := h.driver.NewPage(agouti.Browser("chrome")); err == nil {
		h.page = page
	} else {
		return nil, err
	}

	if err := h.page.Navigate("https://hyou.show"); err != nil {
		log.Fatalln(err)
	}
	return h, nil
}

func (h *hyoushowPage) Close() error {
	return h.driver.Stop()
}

func (h *hyoushowPage) SetReceiver(receiver string) error {
	receiverElem := h.page.Find("#receiver")
	if receiverElem == nil {
		return errors.New("receiver not found")
	}
	return receiverElem.Fill(receiver)
}

func (h *hyoushowPage) SetSender(sender string) error {
	senderElem := h.page.Find("#sender")
	if senderElem == nil {
		return errors.New("sender not found")
	}
	return senderElem.Fill(sender)
}

func (h *hyoushowPage) SetDateStr(dateStr string) error {
	dateElem := h.page.Find("#date")
	if dateElem == nil {
		return errors.New("date not found")
	}
	return dateElem.Fill(dateStr)
}

func (h *hyoushowPage) SetCertificateBody(certificateBody string) error {
	return h.page.Find("#certificate_body").Fill(certificateBody)
}

func (h *hyoushowPage) Generate() (string, error) {
	if err := h.page.Find(".share_link").Click(); err != nil {
		return "", err
	}
	if resultURL, err := h.page.Find("#result_cert").Attribute("src"); err == nil {
		for resultURL == "https://hyou.show/image/cert.png" {
			if resultURL, err = h.page.Find("#result_cert").Attribute("src"); err != nil {
				return "", err
			}
			time.Sleep(1 * time.Second)
		}
		return resultURL, nil
	} else {
		return "", err
	}
}
