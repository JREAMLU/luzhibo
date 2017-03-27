package workers

import (
	"bytes"
	"io"
	"net/http"
	"crypto/tls"
	"github.com/Baozisoftware/luzhibo/api/getters"
	"os/exec"
	"strings"
)

//下载器

type downloader struct {
	url      string
	filePath string
	cb       WorkCompletedCallBack
	run      bool
	ch       chan bool
	ch2      chan bool
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
	i.ch2 = make(chan bool, 1)
	if strings.Contains(i.url, "rtmp://") || strings.Contains(i.url, ".m3u8") {
		go i.ffmpeg(i.url, i.filePath)
	} else {
		go i.http(i.url, i.filePath)
	}
}

//Stop 实现接口
func (i *downloader) Stop() {
	if i.run {
		i.ch2 <- true
		i.run = false
		<-i.ch
		close(i.ch)
		close(i.ch2)
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

func (i *downloader) http(url, filepath string) {
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
	resp, err := httpGetResp(url)
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
	client := http.Client{Transport: tr}
	req, err = http.NewRequest("GET", url, nil)
	if err == nil {
		resp, err = client.Do(req)
	}
	return
}

func (i *downloader) ffmpeg(url, filepath string) {
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
	cmd := exec.Command("ffmpeg", "-y", "-i", i.url, "-vcodec", "copy", "-acodec", "copy", i.filePath)
	go func() {
		if err := cmd.Start(); err != nil {
			ec = 2 //ffmpeg启动失败
			i.ch2 <- true
		}
		cmd.Wait()
		if i.run {
			i.ch2 <- true
		}
	}()
	<-i.ch2
	if cmd.Process != nil {
		cmd.Process.Kill()
	}
}
