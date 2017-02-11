package getters

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
	"time"
)

//实现一些通用函数/结构

func httpGetWithUA(url, ua string) (data string, err error) {
	resp, err := httpGetResp(url, ua)
	var body []byte
	if err == nil {
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		if err == nil {
			data = string(body)
		}
	}
	return
}

func httpGet(url string) (data string, err error) {
	return httpGetWithUA(url, "")
}

func httpGetResp(url, ua string) (resp *http.Response, err error) {
	if ua == "" {
		ua = "User-Agent:Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36"
	}
	var req *http.Request
	client := http.Client{}
	req, err = http.NewRequest("GET", url, nil)
	if err == nil {
		req.Header.Set("User-Agent", ua)
		resp, err = client.Do(req)
	}
	return
}

func getUnixTimesTamp() int64 {
	return time.Now().Unix()
}

func getMD5String(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	r := m.Sum(nil)
	return hex.EncodeToString(r)
}

func forEach(list interface{}, f func(interface{}) bool, maxCount int) (result []interface{}, err error) {
	t := 0
	switch reflect.TypeOf(list).Kind() {
	case reflect.Array:
		t = 1
	case reflect.Slice:
		t = 2
	default:
		err = errors.New("list type error")
	}
	if maxCount < 0 {
		err = errors.New("list count error")
	} else if t > 1 {
		maxCount--
		value := reflect.ValueOf(list)
		count := value.Len()
		tmp := make([]interface{}, 0)
		defer func() {
			if recover() != nil {
				err = errors.New("fild to for each")
			}
		}()
		for i := 0; i < count; i++ {
			if v := value.Index(i).Interface(); f(v) {
				tmp = append(tmp, v)
			}
			if i == maxCount {
				break
			}
		}
		if len(tmp) > 0 {
			result = tmp
		} else {
			err = errors.New("not find")
		}
	}
	return
}

func forEachOne(list interface{}, f func(interface{}) bool) (result interface{}, err error) {
	tmp, err := forEach(list, f, 1)
	if err == nil {
		result = tmp[0]
	}
	return
}

type jObject map[string]interface{}

func pruseJSON(data string) *jObject {
	var o interface{}
	if json.Unmarshal([]byte(data), &o) == nil {
		if m, ok := o.(map[string]interface{}); ok {
			j := jObject(m)
			return &j
		}
	}
	return nil
}

func (v *jObject) jToken(key string) *jObject {
	t, ok := (*v)[key]
	if ok {
		if m, ok := t.(map[string]interface{}); ok {
			j := jObject(m)
			return &j
		}
	}
	return nil
}

func (v *jObject) jArray(key string) []interface{} {
	t, ok := (*v)[key]
	if ok {
		if m, ok := t.([]interface{}); ok {
			return m
		}
	}
	return nil
}

func (v *jObject) jTokens(key string) []*jObject {
	t := v.jArray(key)
	if t != nil {
		r := make([]*jObject, 0)
		for _, v := range t {
			if m, ok := v.(map[string]interface{}); ok {
				j := jObject(m)
				r = append(r, &j)
			}
		}
		return r
	}
	return nil
}

func randInt64(min, max int64) int64 {
	rand.Seed(time.Now().Unix())
	return min + rand.Int63n(max-min)
}

//LiveInfo 直播间信息结构
type LiveInfo struct {
	RoomTitle   string
	LivingIMG   string
	VideoURL    string
	RoomDetails string
	RoomID      string
	LiveNick    string
}

//Getter 房间/直播信息获取接口
type Getter interface {
	GetRoomInfo(string) (string, bool, error) //获取房间信息,参数为房间地址,返回房间号,是否开播
	GetLiveInfo(string) (LiveInfo, error)     //获取直播信息,参数为房间号,返回直播信息
	Site() string                             //返回平台名称
	SiteURL() string                          //返回平台首页
}

//Getters 所有获取接口
var Getters = []Getter{&douyu{}, &panda{}, &zhanqi{}, &longzhu{}, &huya{}, &qie{}, &bilibili{}, &quanmin{}, &huajiao{}, &huomao{}}
