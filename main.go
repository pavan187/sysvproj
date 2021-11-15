package main

import (
        "fmt"
        "io"
        "io/ioutil"
        "log"
        "net/http"
        "strings"
        "time"
)

type Message struct {
        UrlFrom string `json:"urlfrom"`
}

var urlinfo1 = []string{"example.com", "www.notrust.com", "www.suspicious.com"}

func urlinfo(w http.ResponseWriter, r *http.Request) {

        fmt.Println(r.URL.RawQuery)
        getUrl := strings.Split(r.URL.RawQuery, "=")
        Url := getUrl[1]
        fmt.Println(Url)
        for _, malurl := range urlinfo1 {
                if strings.Contains(malurl, Url) {
                        fmt.Println("this is malware URL, cannot be accessed")
                        return
                }
        }

        fmt.Fprintf(w, "This is the response from urlinfo Server")

}


func hello(w http.ResponseWriter, req *http.Request) {
        var bytes io.Reader
        fmt.Println("remore addr:", req.RemoteAddr)
        s := strings.Split(req.RemoteAddr, ":")
        clientInfo := s[0]
        fmt.Println("Client Info: ", clientInfo)

        client := &http.Client{}
        t := time.Now()
        hostname := strings.Split(req.RemoteAddr, ":")
        st:= fmt.Sprint("http://localhost:12345/urlinfo?urlfrom=%s", hostname[0])
        req, err := http.NewRequest(http.MethodGet, st, bytes)
        //req.Header.Set("Content-Type", "application/json; charset=utf-8")
        if err != nil {
                fmt.Println("http request sending failed !!")
        }
        resp, err := client.Do(req)
        if err != nil {
                fmt.Println("Sending HTTP request failed: %v", err)
        }
        fmt.Println("%s %s %v in %v", req.Method, req.URL, resp.Status, time.Since(t))
        contents, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                fmt.Println("HTTP Response reading has failed: %v", err)
        }

        w.Write(contents)
        defer resp.Body.Close()



}
func main()  {

        malwareInfo := http.NewServeMux()
        malwareInfo.HandleFunc("/urlinfo", urlinfo)


        //Running Second Server
        go func() {
                log.Println("Server started on: http://localhost:12345")
                http.ListenAndServe(":12345", malwareInfo)
        }()

        log.Println("Server started on: http://localhost:8080")
        http.HandleFunc("/hello", hello)
        http.ListenAndServe(":8080", nil)
}
