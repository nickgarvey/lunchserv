package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type lunchHandler struct {
	mutex sync.Mutex
	url   []byte
	orderExpireTime time.Time
}

func (handler *lunchHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Println(request.RemoteAddr, request.Header)
	switch request.Method {
	case "POST":
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Println("Error: ", err)
			return
		}
		matched, err := regexp.Match("^https://drd\\.sh/cart/[a-zA-Z0-9]+/?$", body)
		if !matched {
			errStr := "Invalid URL: " + string(body)
			http.Error(writer, errStr, 400)
			return
		}
		handler.mutex.Lock()
		defer handler.mutex.Unlock()
		handler.url = body
		handler.orderExpireTime = time.Now().Add(time.Hour * 12)
		log.Println(string(handler.url), handler.orderExpireTime)
	case "GET":
		handler.mutex.Lock()
		defer handler.mutex.Unlock()

		// invalidate the url if it's been too long
		if handler.orderExpireTime.Before(time.Now()) {
			handler.url = nil
		}

		if handler.url == nil {
			writer.Write([]byte("Order hasn't been created yet\n"))
			return
		}
		http.Redirect(writer, request, string(handler.url), http.StatusTemporaryRedirect)
		writer.Write(handler.url)
	}
}

func main() {
	handler := new(lunchHandler)
	handler.url = nil
	http.Handle("/", handler)
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s port", os.Args[0])
	}
	port, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	listenStr := ":" + strconv.FormatInt(port, 10)
	log.Println("Listening on " + listenStr)
	log.Fatal(http.ListenAndServe(listenStr, nil))
}
