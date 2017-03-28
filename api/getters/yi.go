package getters

import (
	"errors"
	"regexp"
)

//yi 一直播
type yi struct{}

//Site 实现接口
func (i *yi) Site() string { return "一直播" }

//SiteURL 实现接口
func (i *yi) SiteURL() string {
	return "http://www.yizhibo.com"
}

//SiteIcon 实现接口
func (i *yi) SiteIcon() string {
	return i.SiteURL() + "/favicon.ico"
}

//FileExt 实现接口
func (i *yi) FileExt() string {
	return "flv"
}

//NeedFFMpeg 实现接口
func (i *yi) NeedFFMpeg() bool {
	return false
}

//GetRoomInfo 实现接口
func (i *yi) GetRoomInfo(url string) (id string, live bool, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	reg, _ := regexp.Compile("yizhibo\\.com/l/(\\S+)\\.html")
	id = reg.FindStringSubmatch(url)[1]
	if id != "" {
		url = "http://api.xiaoka.tv/live/web/get_play_live?scid=" + id
		var tmp string
		tmp, err = httpGet(url)
		json := pruseJSON(tmp)
		if (*json)["result"].(float64) == 1 {
			live = (*json.jToken("data"))["status"].(float64) == 10
		} else {
			id = ""
		}
	}
	if id == "" {
		err = errors.New("fail get data")
	}
	return
}

//GetLiveInfo 实现接口
func (i *yi) GetLiveInfo(id string) (live LiveInfo, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	live = LiveInfo{RoomID: id}
	url := "http://api.xiaoka.tv/live/web/get_play_live?scid=" + id
	tmp, err := httpGet(url)
	json := *(pruseJSON(tmp).jToken("data"))
	nick := json["nickname"].(string)
	title := json["title"].(string)
	video := json["linkurl"].(string)
	img := json["cover"].(string)
	img = "http://alcdn.img.xiaoka.tv/" + img
	live.LiveNick = nick
	live.RoomTitle = title
	live.RoomDetails = ""
	live.LivingIMG = img
	live.VideoURL = video
	if video == "" {
		err = errors.New("fail get data")
	}
	return
}
