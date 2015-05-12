package main

import (
	"bytes"
	"compress/zlib"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
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

type ByInt []Entry

func (entries ByInt) Len() int           { return len(entries) }
func (entries ByInt) Swap(i, j int)      { entries[i], entries[j] = entries[j], entries[i] }
func (entries ByInt) Less(i, j int) bool { return entries[i].Int < entries[j].Int }

func main() {
	/*buf, err := ioutil.ReadFile("data.z")
	if err != nil {
		log.Fatalln("failed to read file 'data.z'")
	}*/

	response, err := http.Get("http://cmts.iqiyi.com/bullet/63/00/365426300_300_1.z?rn=0.6115525625646114")
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

	var result Result
	err = xml.Unmarshal(str_buf, &result)
	// log.Println(result)
	log.Printf("code : %s\n", result.Code)
	num_of_scales := len(result.Dat.Entries)
	log.Printf("number of scales : %d\n", num_of_scales)
	sort.Sort(ByInt(result.Dat.Entries))
	for _, entry := range result.Dat.Entries {
		log.Printf("\tscale : %d, number of bullets : %d\n", entry.Int, len(entry.BList.Bullets))
	}
}
