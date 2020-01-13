package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/fliropp/gokube/pkg/grpcclient"
	"github.com/fliropp/gokube/pkg/httpclient"
	"github.com/sirupsen/logrus"
)

// inspired by: https://medium.com/statuscode/how-i-write-go-http-services-after-seven-years-37c208122831

type WebServer struct {
	router *http.ServeMux
	log    *logrus.Logger
}

func NewWebServer(log *logrus.Logger) *WebServer {
	s := &WebServer{}
	s.log = log
	s.router = http.NewServeMux()
	s.routes()
	return s
}

func (s *WebServer) Start() {
	go func() {
		err := http.ListenAndServe(":8080", s.router)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error starting web server: %s", err.Error()))
		}
	}()
}

func (s *WebServer) routes() {
	prefix := "/gokube/"
	router := http.NewServeMux()
	router.HandleFunc("/ping", s.handlePing())
	router.HandleFunc("/whoami", s.handleWhoAmI())
	router.HandleFunc("/getdata", s.handleGoPyKube())
	router.HandleFunc("/grpcrequest", s.handleGrpcClient())
	s.AddHandle(prefix, router)
}

func (s *WebServer) handlePing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong!"))
	}
}

func (s *WebServer) handleGrpcClient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		grpcclient.RunGrpcClient()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("gRPC request sent"))
	}
}

func (s *WebServer) handleGoPyKube() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := httpclient.GetHttpClient()
		var response interface{}
		getURL := "http://192.168.64.2/pykube"
		//getURL := "pykube-service"

		req, err := http.NewRequest("GET", getURL, nil)
		if err != nil {
			fmt.Println(fmt.Errorf("ERROR: create reqeuest failed (%s)", err.Error()))
		}
		req.Header.Add("Accept", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(fmt.Errorf("ERROR: request failed (%s)", err.Error()))

		}

		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			fmt.Println(fmt.Errorf("ERROR: decode failed (%s)", err.Error()))

		}
		result, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result))
	}
}

func (s *WebServer) handleWhoAmI() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("I am Are"))
		return
	}
}

func (s *WebServer) AddHandle(prefix string, router *http.ServeMux) {
	s.router.Handle(prefix, http.StripPrefix(strings.TrimSuffix(prefix, "/"), router))
}
