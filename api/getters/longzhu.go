package getters

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

//longzhu 龙珠直播
type longzhu struct{}

//Site 实现接口
func (i *longzhu) Site() string { return "龙珠直播" }

//SiteURL 实现接口
func (i *longzhu) SiteURL() string {
	return "http://www.longzhu.com"
}

//SiteIcon 实现接口
func (i *longzhu) SiteIcon() string {
	return i.SiteURL() + "/favicon.ico"
}

//FileExt 实现接口
func (i *longzhu) FileExt() string {
	return "flv"
}

//NeedFFMpeg 实现接口
func (i *longzhu) NeedFFMpeg() bool {
	return false
}

//GetRoomInfo 实现接口
func (i *longzhu) GetRoomInfo(url string) (id string, live bool, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	url = strings.ToLower(url)
	reg, _ := regexp.Compile("longzhu\\.com/(\\w+)")
	id = reg.FindStringSubmatch(url)[1]
	if id != "" {
		url = "http://searchapi.plu.cn/api/search/room?title=" + id
		var tmp string
		tmp, err = httpGet(url)
		json := pruseJSON(tmp).jTokens("items")
		if len(json) > 0 {
			live = (*json[0].jToken("live"))["isLive"].(bool)
		}
	}
	if id == "" {
		err = errors.New("fail get data")
	}
	return
}

//GetLiveInfo 实现接口
func (i *longzhu) GetLiveInfo(id string) (live LiveInfo, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	live = LiveInfo{}
	url := "http://roomapicdn.plu.cn/room/RoomAppStatusV2?domain=" + id
	tmp, err := httpGet(url)
	json := *(pruseJSON(tmp).jToken("BaseRoomInfo"))
	nick := json["Name"].(string)
	title := json["BoardCastTitle"].(string)
	details := json["Desc"].(string)
	_id := json["Id"]
	live.RoomID = fmt.Sprintf("%.f", _id)
	url = "http://livestream.plu.cn/live/getlivePlayurl?roomId=" + live.RoomID
	tmp, err = httpGet(url)
	json = *(pruseJSON(tmp).jTokens("playLines")[0].jTokens("urls")[0])
	video := json["securityUrl"].(string)
	live.LiveNick = nick
	live.RoomTitle = title
	live.RoomDetails = details
	live.LivingIMG = ""
	live.VideoURL = video
	if video == "" {
		err = errors.New("fail get data")
	}
	return
}
