package getters

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//huya 虎牙直播
type huya struct{}

//SiteURL 实现接口
func (i *huya) SiteURL() string {
	return "http://www.huya.com"
}

//Site 实现接口
func (i huya) Site() string { return "虎牙直播" }


//GetRoomInfo 实现接口
func (i *huya) GetRoomInfo(url string) (id string, live bool, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	url = strings.ToLower(url)
	reg, _ := regexp.Compile("huya\\.com/(\\w+)")
	id = reg.FindStringSubmatch(url)[1]
	url = "http://m.huya.com/" + id
	html, err := httpGetWithUA(url, ipadUA)
	if !strings.Contains(html, "找不到此页面") {
		live = strings.Contains(html, "ISLIVE = true")
	} else {
		err = errors.New("fail get id")
	}
	if id == "" {
		err = errors.New("fail get data")
	}
	return
}

//GetLiveInfo 实现接口
func (i *huya) GetLiveInfo(id string) (live LiveInfo, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	live = LiveInfo{RoomID: id}
	url := "http://m.huya.com/" + id
	resp, err := httpGetResp(url, ipadUA)
	doc, err := goquery.NewDocumentFromResponse(resp)
	n := doc.Find("div.live-info-desc")
	nick := n.Find("h2").Text()
	title := n.Find("h1").Text()
	details := doc.Find("div.notice_content").Text()
	details = strings.TrimSpace(details)
	n = doc.Find("video#html5player-video")
	img, _ := n.Attr("poster")
	t, _ := doc.Find("source").Attr("src")
	reg, _ := regexp.Compile("\\d+-\\d+")
	t = reg.FindString(t)
	t = strings.Replace(t, "-", "_", -1)
	if t != "" {
		video := fmt.Sprintf("http://hls.yy.com/%s_100571200.flv", t)
		live.LiveNick = nick
		live.LivingIMG = img
		live.RoomDetails = details
		live.RoomTitle = title
		live.VideoURL = video
	} else {
		err = errors.New("faild get data")
	}
	if live.VideoURL == "" {
		err = errors.New("fail get data")
	}
	return
}
