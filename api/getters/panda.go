package getters

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

//panda 熊猫直播
type panda struct{}

//Site 实现接口
func (i *panda) Site() string { return "熊猫直播" }

//实现接口
func (i *panda) SiteURL() string {
	return "http://www.panda.tv"
}

//GetRoomInfo 实现接口
func (i *panda) GetRoomInfo(url string) (id string, live bool, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	url = strings.ToLower(url)
	reg, _ := regexp.Compile("panda\\.tv/(\\d+)")
	id = reg.FindStringSubmatch(url)[1]
	url = "http://www.panda.tv/ajax_search?roomid=" + id
	tmp, err := httpGet(url)
	json := pruseJSON(tmp).jToken("data").jTokens("items")
	var r interface{}
	r, err = forEachOne(json, func(v interface{}) bool { return (*v.(*jObject))["roomid"] == id })
	live = (*r.(*jObject))["status"] == "2"
	if id == "" {
		err = errors.New("fail get data")
	}
	return
}

//GetLiveInfo 实现接口
func (i *panda) GetLiveInfo(id string) (live LiveInfo, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	live = LiveInfo{RoomID: id}
	url := "http://www.panda.tv/api_room_v2?__plat=pc_web&roomid="
	url = fmt.Sprintf("%s%s&_=%d", url, id, getUnixTimesTamp())
	tmp, err := httpGet(url)
	json := pruseJSON(tmp).jToken("data")
	roomInfo, videoInfo, hostInfo := *(json.jToken("roominfo")), *(json.jToken("videoinfo")), *(json.jToken("hostinfo"))
	nick := hostInfo["name"].(string)
	title := roomInfo["name"].(string)
	details := roomInfo["bulletin"].(string)
	img := (*(roomInfo.jToken("pictures")))["img"].(string)
	key := videoInfo["room_key"].(string)
	flag := videoInfo["plflag"].(string)
	flag = strings.Split(flag, "_")[1]
	video := fmt.Sprintf("http://pl%s.live.panda.tv/live_panda/%s.flv", flag, key)
	live.LiveNick = nick
	live.RoomTitle = title
	live.RoomDetails = details
	live.LivingIMG = img
	live.VideoURL = video
	if video == "" {
		err = errors.New("fail get data")
	}
	return
}
