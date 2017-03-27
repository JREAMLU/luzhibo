package getters

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

//afreeca AfreecaTV
type afreeca struct{}

//Site 实现接口
func (i *afreeca) Site() string { return "AfreecaTV" }

//SiteURL 实现接口
func (i *afreeca) SiteURL() string {
	return "http://www.afreecatv.com"
}

//GetRoomInfo 实现接口
func (i *afreeca) GetRoomInfo(url string) (id string, live bool, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	url = strings.ToLower(url)
	reg, _ := regexp.Compile("play\\.afreecatv\\.com/(\\w+)/\\d+")
	id = reg.FindStringSubmatch(url)[1]
	tmp, err := httpGet(url)
	if !strings.Contains(tmp, fmt.Sprintf("id : '%s'", id)) {
		id = ""
	} else {
		live = strings.Contains(tmp, "\"og:title\" content=\"[생]")
	}
	if id == "" {
		err = errors.New("fail get data")
	}
	return
}

//GetLiveInfo 实现接口
func (i *afreeca) GetLiveInfo(id string) (live LiveInfo, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	live = LiveInfo{RoomID: id}
	url := "http://play.afreecatv.com/" + id
	tmp, err := httpGet(url)
	reg, _ := regexp.Compile("no : '(\\d+)'")
	rid := reg.FindStringSubmatch(tmp)[1]
	tmp, err = httpPost("http://live.afreecatv.com:8057/afreeca/player_live_api.php", "bno="+rid)
	json := *(pruseJSON(tmp).jToken("CHANNEL"))
	nick := fmt.Sprint(json["BJNICK"])
	title := fmt.Sprint(json["TITLE"])
	img := fmt.Sprintf("http://liveimg.afreecatv.com/%s.gif", rid)
	url = fmt.Sprintf("http://sessionmanager01.afreeca.tv:6060/broad_stream_assign.html?broad_key=%s-flash-hd-rtmp", rid)
	tmp, err = httpGet(url)
	json = *pruseJSON(tmp)
	video := fmt.Sprint(json["view_url"])
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
