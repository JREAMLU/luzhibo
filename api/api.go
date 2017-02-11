package api

import (
	"errors"
	"luzhibo/api/getters"
	"regexp"
	"strings"
)

//LuzhiboAPI API object
type LuzhiboAPI struct {
	id      string
	URL     string
	g       getters.Getter
	Site    string
	SiteURL string
	Icon    string
}

//New 使用网址创建一个实例
func New(url string) *LuzhiboAPI {
	g := getGetter(url)
	if g != nil {
		i := &LuzhiboAPI{}
		i.g = g
		i.URL = url
		i.Site = g.Site()
		i.SiteURL = g.SiteURL()
		i.Icon = i.SiteURL + "/favicon.ico"
		return i
	}
	return nil
}

//GetRoomInfo 取直播间信息
func (i *LuzhiboAPI) GetRoomInfo() (id string, live bool, err error) {
	if i.URL == "" || i.g == nil {
		err = errors.New("not has url or not found getter")
		return
	}
	id, live, err = i.g.GetRoomInfo(i.URL)
	i.id = id
	return
}

//GetLiveInfo 取直播信息
func (i *LuzhiboAPI) GetLiveInfo() (live getters.LiveInfo, err error) {
	if i.id == "" || i.g == nil {
		err = errors.New("not has id or not found getter")
		return
	}
	live, err = i.g.GetLiveInfo(i.id)
	return
}

func getGetter(url string) getters.Getter {
	url = strings.ToLower(url)
	regs := []string{"(douyu\\.tv)|((douyu)|(douyutv)\\.com)",
		"panda\\.tv",
		"zhanqi\\.tv",
		"longzhu\\.com",
		"huya\\.com",
		"live\\.qq\\.com",
		"live\\.bilibili\\.com",
		"quanmin\\.tv",
		"huajiao\\.com",
		"huomao\\.com"}
	for i := 0; i < len(getters.Getters); i++ {
		if ok, _ := regexp.MatchString(regs[i], url); ok {
			return getters.Getters[i]
		}
	}
	return nil
}
