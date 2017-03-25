package getters

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

//bilibili Bilibili直播
type bilibili struct{}

//Site 实现接口
func (i *bilibili) Site() string { return "Bilibili直播" }

//SiteURL 实现接口
func (i *bilibili) SiteURL() string {
	return "http://live.bilibili.com"
}

//GetRoomInfo 实现接口
func (i *bilibili) GetRoomInfo(url string) (id string, live bool, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	tmp, err := httpGet(url)
	reg, _ := regexp.Compile("ROOMID = (\\d+)")
	id = reg.FindStringSubmatch(tmp)[1]
	url = "http://live.bilibili.com/live/getInfo?roomid=" + id
	tmp, err = httpGet(url)
	if !strings.Contains(tmp, "\\u623f\\u95f4\\u4e0d\\u5b58\\u5728") {
		live = strings.Contains(tmp, "\"_status\":\"on\"")
	} else {
		err = errors.New("fild get data")
	}
	if id == "" {
		err = errors.New("fail get data")
	}
	return
}

//GetLiveInfo 实现接口
func (i *bilibili) GetLiveInfo(id string) (live LiveInfo, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	live = LiveInfo{RoomID: id}
	api := "room_info"
	args := "room_id=" + id
	url := getAPIURL(api, args)
	tmp, err := httpGet(url)
	json := *(pruseJSON(tmp).jToken("data"))
	title := json["title"].(string)
	nick := json["uname"].(string)
	img := json["cover"].(string)
	cid := (*json.jToken("schedule"))["cid"].(float64)
	if cid > 10000000 {
		cid -= 10000000
	}
	api = "playurl"
	args = fmt.Sprintf("cid=%.f&rnd=%d", cid, randInt64(100, 9999))
	url = getAPIURL(api, args)
	tmp, err = httpGet(url)
	x, y := strings.Index(tmp, "<url><![CDATA[")+14, strings.LastIndex(tmp, "]]></url>")
	video := tmp[x:y]
	live.LiveNick = nick
	live.RoomTitle = title
	live.LivingIMG = img
	live.RoomDetails = ""
	live.VideoURL = video
	if video == "" {
		err = errors.New("fail get data")
	}
	return
}

func getAPIURL(api, args string) string {
	t1 := fmt.Sprintf("http://live.bilibili.com/api/%s?", api)
	t2 := "appkey=422fd9d7289a1dd9&" + args
	t3 := t2 + "ba3a4e554e9a6e15dc4d1d70c2b154e3"
	t4 := "&sign=" + getMD5String(t3)
	r := t1 + t2 + t4
	return r
}
