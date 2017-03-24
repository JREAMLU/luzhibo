package workers

import (
	"bytes"
	"io"
	"luzhibo/api/getters"
	"net/http"
	"crypto/tls"
)

//下载器

type downloader struct {
	url      string
	filePath string
	cb       WorkCompletedCallBack
	run      bool
	ch       chan bool
}

func newDownloader(url, filepath string, callbcak WorkCompletedCallBack) *downloader {
	if url != "" && filepath != "" {
		r := &downloader{}
		r.url = url
		r.filePath = filepath
		r.cb = callbcak
		return r
	}
	return nil
}

//Start 实现接口
func (i *downloader) Start() {
	if i.run {
		return
	}
	i.run = true
	i.ch = make(chan bool, 0)
	go i.download(i.url, i.filePath)
}

//Stop 实现接口
func (i *downloader) Stop() {
	if i.run {
		i.run = false
		<-i.ch
		close(i.ch)
	}
}

//Restart 实现接口
func (i *downloader) Restart() (Worker, error) {
	if i.run {
		i.Stop()
	}
	i.Start()
	return i, nil
}

//GetTaskInfo 实现接口
func (i *downloader) GetTaskInfo(g bool) (int64, bool, int64, string, *getters.LiveInfo) {
	return 0, i.run, 0, i.filePath, nil
}

func (i *downloader) download(url, filepath string) {
	ec := int64(0) //正常停止
	defer func() {
		if !i.run {
			i.ch <- true
		}
		if !i.run {
			ec = 1 //主动停止
		}
		i.run = false
		if i.cb != nil {
			i.cb(ec)
		}
	}()
	resp, err :=httpGetResp(url)
	if err != nil || resp.StatusCode != 200 {
		ec = 2 //请求时错误
		return
	}
	defer resp.Body.Close()
	f, err := createFile(filepath)
	if err != nil {
		ec = 3 //创建文件错误
		return
	}
	defer f.Close()
	buf := make([]byte, bytes.MinRead)
	for i.run {
		t, err := resp.Body.Read(buf)
		if err != nil {
			if err == io.EOF {
				f.Write(buf[:t])
			} else {
				ec = 4 //下载数据错误
			}
			return
		}
		_, err = f.Write(buf[:t])
		if err != nil {
			ec = 5 //写入文件错误
			return
		}
	}
}

func httpGetResp(url string) (resp *http.Response, err error) {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	var req *http.Request
	client := http.Client{Transport:tr}
	req, err = http.NewRequest("GET", url, nil)
	if err == nil {
		resp, err = client.Do(req)
	}
	return
}
