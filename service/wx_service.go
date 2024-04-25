package service

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

const token = "hello112233"

// 用于接收微信消息的结构体
type TextMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
	MsgId        int64    `xml:"MsgId"`
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		validateToken(w, r)
	case "POST":
		processMessage(w, r)
	default:
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
}

func validateToken(w http.ResponseWriter, r *http.Request) {
	signature := r.URL.Query().Get("signature")
	timestamp := r.URL.Query().Get("timestamp")
	nonce := r.URL.Query().Get("nonce")
	echostr := r.URL.Query().Get("echostr")

	tmpStrs := sort.StringSlice{token, timestamp, nonce}
	sort.Sort(tmpStrs)
	tmpStr := strings.Join(tmpStrs, "")

	hasher := sha1.New()
	io.WriteString(hasher, tmpStr)
	hashcode := fmt.Sprintf("%x", hasher.Sum(nil))

	if hashcode == signature {
		w.Write([]byte(echostr))
	} else {
		w.Write([]byte("verification failed"))
	}
}

func processMessage(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var msg TextMessage
	if err := xml.Unmarshal(body, &msg); err != nil {
		http.Error(w, "Failed to parse message", http.StatusInternalServerError)
		return
	}

	// 处理消息，例如回复用户
	fmt.Fprintf(w, "Received message: %+v", msg)
	w.Write([]byte("hello"))
}
