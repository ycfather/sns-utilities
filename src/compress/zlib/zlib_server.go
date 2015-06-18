package main

import (
	"bytes"
	"compress/zlib"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"
)

type Result struct {
	Code string `xml:"code"`
	Dat  Data   `xml:"data"`
}

type Data struct {
	Entries []Entry `xml:"entry"`
}

type Entry struct {
	Int   int32      `xml:"int"`
	BList BulletList `xml:"list"`
}

type BulletList struct {
	Bullets []BulletInfo `xml:"bulletInfo"`
}

type BulletInfo struct {
	Id          string `xml:"contentId"`
	Content     string `xml:"content"`
	ShowTime    int32  `xml:"showTime"`
	AddTime     int32  `xml:"addTime"`
	Likes       int32  `xml:"likes"`
	Font        int32  `xml:"font"`
	Color       string `xml:"color"`
	Opacity     int32  `xml:"opacity"`
	Position    int32  `xml:"position"`
	Background  int32  `xml:"background"`
	ReplyUid    int32  `xml:"replyUid"`
	ContentType int32  `xml:"contentType"`
}

type EntrySorter []Entry

func (entries EntrySorter) Len() int           { return len(entries) }
func (entries EntrySorter) Swap(i, j int)      { entries[i], entries[j] = entries[j], entries[i] }
func (entries EntrySorter) Less(i, j int) bool { return entries[i].Int < entries[j].Int }

func Parse(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		log.Fatalf("Unsupported http method : %s\n", request.Method)
		return
	}
	data_url := request.FormValue("data-url")
	if len(data_url) == 0 {
		log.Fatalf("the url [%s] is invalid\n", data_url)
		writer.WriteHeader(500)
	} else {
		response, err := http.Get(data_url)
		defer response.Body.Close()
		if err != nil {
			log.Fatalln("failed to get content from the specified url")
			return
		}

		log.Printf("Content-Type : %s, Content-Length : %s\n", response.Header.Get("Content-Type"), response.Header.Get("Content-Length"))

		buf, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatalln("failed to read response")
			return
		}
		log.Printf("read %d bytes from the url\n", len(buf))

		zlib_buf := bytes.NewReader(buf)
		r, err := zlib.NewReader(zlib_buf)
		defer r.Close()
		if err != nil {
			panic(err)
		}

		// io.Copy(os.Stdout, r)
		str_buf, err := ioutil.ReadAll(r)
		// fmt.Println(string(str_buf))
		writer.Write(str_buf)

		var result Result
		err = xml.Unmarshal(str_buf, &result)
		// log.Println(result)
		log.Printf("code : %s\n", result.Code)
		num_of_scales := len(result.Dat.Entries)
		log.Printf("number of scales : %d\n", num_of_scales)
		sort.Sort(EntrySorter(result.Dat.Entries))
		for _, entry := range result.Dat.Entries {
			log.Printf("\tscale : %d, number of bullets : %d\n", entry.Int, len(entry.BList.Bullets))
		}
	}
}

func main() {
	http.HandleFunc("/parse", Parse)
	server := &http.Server{
		Addr:           ":8080",
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}
