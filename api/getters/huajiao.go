package getters

import (
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//huajiao 花椒直播
type huajiao struct{}

//Site 实现接口
func (i *huajiao) Site() string { return "花椒直播" }

//SiteURL 实现接口
func (i *huajiao) SiteURL() string {
	return "http://www.huajiao.com"
}

//SiteIcon 实现接口
func (i *huajiao) SiteIcon() string {
	return i.SiteURL() + "/favicon.ico"
}

//FileExt 实现接口
func (i *huajiao) FileExt() string {
	return "flv"
}

//NeedFFMpeg 实现接口
func (i *huajiao) NeedFFMpeg() bool {
	return false
}

//GetRoomInfo 实现接口
func (i *huajiao) GetRoomInfo(url string) (id string, live bool, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	url = strings.ToLower(url)
	reg, _ := regexp.Compile("huajiao\\.com/l/(\\d+)")
	id = reg.FindStringSubmatch(url)[1]
	url = "http://h.huajiao.com/l/index?liveid=" + id
	tmp, err := httpGet(url)
	if !strings.Contains(tmp, "err-d4bcf8ad0d.png") {
		live = !strings.Contains(tmp, "直播已结束")
	} else {
		err = errors.New("fail get data")
	}
	if id == "" {
		err = errors.New("fail get data")
	}
	return
}

//GetLiveInfo 实现接口
func (i *huajiao) GetLiveInfo(id string) (live LiveInfo, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	live = LiveInfo{RoomID: id}
	url := "http://h.huajiao.com/l/index?liveid=" + id

	resp, err := httpGetResp(url, "")
	doc, err := goquery.NewDocumentFromResponse(resp)
	n := doc.Find("script[type=\"text/javascript\"]")
	tmp := n.Text()
	x, y := strings.Index(tmp, "{")+1, strings.LastIndex(tmp, "}")
	tmp = tmp[x:y]
	tmp = strings.Split(tmp, "\n")[1]
	x, y = strings.Index(tmp, "{"), len(tmp)-1
	tmp = tmp[x:y]
	json := pruseJSON(tmp)
	author, feed := *(json.jToken("author")), *(json.jToken("feed"))
	sn := feed["sn"]
	nick := author["nickname"].(string)
	title := feed["title"].(string)
	url = fmt.Sprintf("http://g2.live.360.cn/liveplay?stype=flv&channel=live_huajiao_v2&bid=huajiao&sn=%s&sid=null&_rate=null&ts=null", sn)
	tmp, err = httpGet(url)
	tmp = tmp[0:3] + tmp[6:]
	bytes, err := base64.StdEncoding.DecodeString(tmp)
	tmp = string(bytes)
	json = pruseJSON(tmp)
	video := (*json)["main"].(string)
	live.LiveNick = nick
	live.LivingIMG = ""
	live.RoomDetails = ""
	live.RoomTitle = title
	live.VideoURL = video
	if video == "" {
		err = errors.New("fail get data")
	}
	return
}
