package getters

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

//quanmin 全民直播
type quanmin struct{}

//Site 实现接口
func (i *quanmin) Site() string { return "全民直播" }

//实现接口
func (i *quanmin) SiteURL() string {
	return "http://www.quanmin.tv"
}

//GetRoomInfo 实现接口
func (i *quanmin) GetRoomInfo(url string) (id string, live bool, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	url = strings.ToLower(url)
	reg, _ := regexp.Compile("quanmin\\.tv/(\\w+)")
	id = reg.FindStringSubmatch(url)[1]
	url = fmt.Sprintf("http://www.quanmin.tv/json/rooms/%s/info1.json", id)
	tmp, err := httpGet(url)
	json := *(pruseJSON(tmp))
	live = json["play_status"].(bool)
	if id == "" {
		err = errors.New("fail get data")
	}
	return
}

//GetLiveInfo 实现接口
func (i *quanmin) GetLiveInfo(id string) (live LiveInfo, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	live = LiveInfo{RoomID: id}
	url := fmt.Sprintf("http://www.quanmin.tv/json/rooms/%s/info1.json", id)
	tmp, err := httpGet(url)
	json := *pruseJSON(tmp)
	nick := json["nick"].(string)
	title := json["title"].(string)
	details := json["intro"].(string)
	img := json["thumb"].(string)
	video := fmt.Sprintf("http://flv.quanmin.tv/live/%s.flv", id)
	live.LiveNick = nick
	live.LivingIMG = img
	live.RoomDetails = details
	live.RoomTitle = title
	live.VideoURL = video
	if video == "" {
		err = errors.New("fail get data")
	}
	return
}
