package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ques0942/hyoushow-api/model"

	"github.com/throttled/throttled"
	"github.com/throttled/throttled/store/memstore"
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

func getParam(r *http.Request, key string) (string, error) {
	if param, ok := r.URL.Query()[key]; ok && len(param) > 0 {
		return param[0], nil
	} else {
		return "", fmt.Errorf("%s not found", key)
	}
}

func determineListenAddress() string {
	port := os.Getenv("PORT")
	if port == "" {
		return ":8080"
	} else {
		return ":" + port
	}
}

func main() {

	store, err := memstore.New(65536)
	if err != nil {
		log.Fatal(err)
	}

	quota := throttled.RateQuota{throttled.PerMin(20), 5}
	rateLimiter, err := throttled.NewGCRARateLimiter(store, quota)
	if err != nil {
		log.Fatal(err)
	}

	httpRateLimiter := throttled.HTTPRateLimiter{
		RateLimiter: rateLimiter,
		VaryBy:      &throttled.VaryBy{Path: true},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		receiver, err := getParam(r, "receiver")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("receiver not found"))
			return
		}
		sender, err := getParam(r, "sender")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("sender not found"))
			return
		}
		dateStr, err := getParam(r, "dateStr")
		body, err := getParam(r, "body")
		if err != nil || body == "" {
			log.Println("body not found")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("body not found"))
			return
		}
		resultURL, err := hyoushow(receiver, sender, dateStr, body)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(resultURL)
		http.Redirect(w, r, resultURL, http.StatusMovedPermanently)
	}
	limitedFunc := httpRateLimiter.RateLimit(http.HandlerFunc(handler))
	http.Handle("/", limitedFunc)

	log.Fatal(http.ListenAndServe(determineListenAddress(), nil))
}
