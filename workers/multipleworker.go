package workers

import "luzhibo/api"
import "luzhibo/api/getters"
import "errors"
import "fmt"
import "time"

//循环模式

type multipleworker struct {
	dirPath string
	index   int64
	cb      WorkCompletedCallBack
	run     bool
	ch      chan bool
	ch2     chan bool
	API     *api.LuzhiboAPI
	sw      *singleworker
}

//NewMultipleWorker 创建对象
func NewMultipleWorker(oa *api.LuzhiboAPI, dirpath string, callbcak WorkCompletedCallBack) (r *multipleworker, err error) {
	if oa != nil {
		r = &multipleworker{}
		_, _, err = oa.GetRoomInfo()
		if err != nil {
			err = errors.New("没有这个房间")
			return
		}
		r.cb = callbcak
		r.dirPath = dirpath
		r.API = oa
		return
	}
	err = errors.New("-1") //参数错误
	return
}

//Start 实现接口
func (i *multipleworker) Start() {
	if i.run {
		return
	}
	i.run = true
	i.ch = make(chan bool, 0)
	go i.do()
}

//Stop 实现接口
func (i *multipleworker) Stop() {
	if i.run {
		i.run = false
		if _, r, _, _, _ := i.sw.GetTaskInfo(false); r {
			i.sw.Stop()
		}
		<-i.ch
		close(i.ch)
	}
}

//Restart 实现接口
func (i *multipleworker) Restart() (Worker, error) {
	if i.run {
		i.Stop()
	}
	r, e := NewMultipleWorker(i.API, i.dirPath, i.cb)
	if e == nil {
		i = r
		i.Start()
	}
	return i, e
}

//GetTaskInfo 实现接口
func (i *multipleworker) GetTaskInfo(g bool) (int64, bool, int64, string, *getters.LiveInfo) {
	if i.sw != nil {
		_, _, _, _, r := i.sw.GetTaskInfo(g)
		return 2, i.run, i.index, i.dirPath, r
	}
	return 2, i.run, i.index, i.dirPath, nil
}

func (i *multipleworker) do() {
	var ec int64
	for i.run {
		i.ch2 = make(chan bool, 0)
		i.index++
		fn := fmt.Sprintf("%s/%d.flv", i.dirPath, i.index)
		r, err := NewSingleWorker(i.API, fn, func(x int64) {
			ec = x
			i.ch2 <- true
		})
		b := false
		if err == nil {
			i.sw = r
			b = true
		} else {
			i.index--
		}
		if b {
			i.sw.Start()
			<-i.ch2
		}
		if ec == 5 {
			break
		}
		if i.run {
			time.Sleep(5 * time.Minute)
		}
	}
	if !i.run {
		i.ch <- true
	}
	i.run = false
	if i.cb != nil {
		i.cb(ec)
	}
}
