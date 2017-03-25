package getters

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//chushou 触手直播
type chushou struct{}

//SiteURL 实现接口
func (i *chushou) SiteURL() string {
	return "https://chushou.tv"
}

//Site 实现接口
func (i chushou) Site() string { return "触手直播" }

//GetRoomInfo 实现接口
func (i *chushou) GetRoomInfo(url string) (id string, live bool, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	url = strings.ToLower(url)
	reg, _ := regexp.Compile("chushou\\.tv/room/(\\d+)\\.htm")
	id = reg.FindStringSubmatch(url)[1]
	url = fmt.Sprintf("https://chushou.tv/room/m-%s.htm",id)
	html, err := httpGetWithUA(url, ipadUA)
	if !strings.Contains(html, "访问失败啦") {
		live = !strings.Contains(html, "playUrl=\"\"")
	} else {
		err = errors.New("fail get id")
	}
	if id == "" {
		err = errors.New("fail get data")
	}
	return
}

//GetLiveInfo 实现接口
func (i *chushou) GetLiveInfo(id string) (live LiveInfo, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	live = LiveInfo{RoomID: id}
	url := fmt.Sprintf("https://chushou.tv/room/m-%s.htm",id)
	resp, _ := httpGetResp(url, ipadUA)
	doc, _ := goquery.NewDocumentFromResponse(resp)
	nick := doc.Find("span.mzb_nickname").Text()
	details := doc.Find("span.announcement_text").Text()
	img, _ := doc.Find("video.videoBlock").Attr("poster")
	t, _ := doc.Find("video.videoBlock").Attr("src")
	reg, _ := regexp.Compile("[a-f0-9]{32}")
	t = reg.FindString(t)
	if t != "" {
		url=fmt.Sprintf("https://chushou.tv/room/%s.htm",id)
		resp, _=httpGetResp(url, "")
		doc, _ = goquery.NewDocumentFromResponse(resp)
		title :=doc.Find("p.zb_player_gamedesc").Text()
		video := fmt.Sprintf("http://hdl6.kascend.com/chushou_live/%s.flv", t)
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
