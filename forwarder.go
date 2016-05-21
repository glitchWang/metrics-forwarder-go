package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Config struct {
	ListenAddr  *string
	SinkAddr    *string
	AllowOrigin *string
	logPath     *string
}

var appConfig *Config

func ParseCommandLine() {
	var listenAddr string
	var sinkAddr string
	var allowOrigin string
	var logPath string
	appConfig = new(Config)
	flag.StringVar(&listenAddr, "listenAddr", ":80", "Listen address, default to :80")
	flag.StringVar(&sinkAddr, "sinkAddr", "dd-agent:8125", "Sink address, default to dd-agent:8125")
	flag.StringVar(&allowOrigin, "allowOrigin", "*", "CORS setting, default to *")
	flag.StringVar(&logPath, "logPath", "/var/log", "Log file path, default to /var/log")
	flag.Parse()
	appConfig.ListenAddr = &listenAddr
	appConfig.SinkAddr = &sinkAddr
	appConfig.AllowOrigin = &allowOrigin
	appConfig.logPath = &logPath
}

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Http to UDP forwarder\n")
}

func HealthCheckHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "OK\n")
}

func CORS(w http.ResponseWriter) {
	w.Header().Add("Access-Control-Allow-Origin", *appConfig.AllowOrigin)
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST")
	w.Header().Add("Access-Control-Max-Age", "604800")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
}

func Forward(content string) {

}

func MetricsReceivedHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		CORS(w)
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func ListenAndServe() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/v1/health-check", HealthCheckHandler)
	http.HandleFunc("/v1/send", MetricsReceivedHandler)
	fmt.Printf("Http forwarder from %s to %s\n", *appConfig.ListenAddr, *appConfig.SinkAddr)
	log.Fatal(http.ListenAndServe(*appConfig.ListenAddr, nil))
}

func main() {
	ParseCommandLine()
	ListenAndServe()
}
