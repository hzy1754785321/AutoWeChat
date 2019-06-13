package main

import (
	"fmt"
	"go-simplejson"
	"io/ioutil"
	e "itchat4go/enum"
	m "itchat4go/model"
	s "itchat4go/service"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/robfig/cron"
)

//SendMsg 发送消息
func SendMsg() {
	var con conf
	conf := con.getConf()
	peoples := conf.NameList
	len := len(peoples)
	friend := make([]string, len)
	city := make([]string, len)
	var yan = GetEveryYan()


	for i, j := range peoples {
		for k, v := range contactMap {
			if v.NickName == j || v.RemarkName == j {
				friend[i] = k
				city[i] = v.City
				break
			}
		}
		var msg string
		msg = GetWeather(city[i])
		var sendMsg = fmt.Sprintf("%s%s\n", msg, yan)
		wxSendMsg := m.WxSendMsg{}
		wxSendMsg.Type = 1
		wxSendMsg.Content = sendMsg
		wxSendMsg.FromUserName = loginMap.SelfUserName
		if friend[i] != "" {
			wxSendMsg.ToUserName = friend[i]
		}else{
			println("user is null")
		}
		wxSendMsg.LocalID = fmt.Sprintf("%d", time.Now().Unix())
		wxSendMsg.ClientMsgId = wxSendMsg.LocalID
		time.Sleep(time.Second)
		err := s.SendMsg(&loginMap, wxSendMsg)
		panicErr(err)
	}

	println("消息发送成功")
}

//GetEveryYan  获取每日一言
func GetEveryYan() (msg string) {

	client := &http.Client{}

	url := "https://api.ooopn.com/yan/api.php"
	reqest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	response, _ := client.Do(reqest)

	//返回的状态码
	status := response.StatusCode
	yanJSON, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(status)
		panic(err)
	}
	//获取json格式页面数据
	js, err := simplejson.NewJson(yanJSON)
	if err != nil {
		panic(err)
	}
	yiyan := js.Get("hitokoto").MustString()
	everyMsg := fmt.Sprintf("每日一言: %s", yiyan)
	return everyMsg
}

//loginIn 登陆微信
func loginIn() {
	/* 从微信服务器获取UUID */
	uuid, err = s.GetUUIDFromWX()
	panicErr(err)

	/* 根据UUID获取二维码 */
	err = s.DownloadImagIntoDir(e.QRCODE_URL+uuid, "./qrcode")
	panicErr(err)
	cmd := exec.Command(`cmd`, `/c start ./qrcode/qrcode.jpg`)
	err = cmd.Run()
	panicErr(err)

	/* 轮询服务器判断二维码是否扫过暨是否登陆了 */
	for {
		fmt.Println("正在验证登陆... ...")
		status, msg := s.CheckLogin(uuid)

		if status == 200 {
			fmt.Println("登陆成功,处理登陆信息...")
			cmd := exec.Command(`cmd`, `/c close ./qrcode/qrcode.jpg`)
			err = cmd.Run()
			loginMap, err = s.ProcessLoginInfo(msg)
			panicErr(err)

			fmt.Println("登陆信息处理完毕,正在初始化微信...")
			err = s.InitWX(&loginMap)
			panicErr(err)

			fmt.Println("初始化完毕,通知微信服务器登陆状态变更...")
			err = s.NotifyStatus(&loginMap)
			panicErr(err)

			fmt.Println("通知完毕,本次登陆信息：")
			fmt.Println(e.SKey + "\t\t" + loginMap.BaseRequest.SKey)
			fmt.Println(e.PassTicket + "\t\t" + loginMap.PassTicket)
			break
		} else if status == 201 {
			fmt.Println("请在手机上确认")
		} else if status == 408 {
			fmt.Println("请扫描二维码")
		} else {
			fmt.Println(msg)
		}
	}
	contactMap, err = s.GetAllContact(&loginMap)
	panicErr(err)
	// for _,v := range contactMap{
	// 	fmt.Print(v.NickName)
	// 	fmt.Print("=========>")
	// 	fmt.Println(v.UserName)
	// }
}

//GetWeather 获取天气消息
func GetWeather(city string) (msg string) {
	client := &http.Client{}
	var cityCode string
	if _, ok := cityDict[city]; ok {
		cityCode = cityDict[city]
	} else {
		cityCode = "101020100"
	}
	url := fmt.Sprintf("http://t.weather.sojson.com/api/weather/city/%s", cityCode)
	reqest, err := http.NewRequest("GET", url, nil)
	panicErr(err)

	//处理返回结果
	response, _ := client.Do(reqest)

	//返回的状态码
	status := response.StatusCode
	weatherJSON, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(status)
		panic(err)
	}

	//获取json格式页面数据
	js, err := simplejson.NewJson(weatherJSON)
	panicErr(err)

	now := time.Now()
	//今日天气
	todayWeather := js.Get("data").Get("forecast").GetIndex(0)
	//天气类型
	weatherType := todayWeather.Get("type").MustString()
	weatherType = fmt.Sprintf("天气 : %s", weatherType)
	//时间
	todayTime := now.Format("2006年01月02日 15:04:05")
	//天气注意事项
	notice := todayWeather.Get("notice").MustString()
	//温度
	high := todayWeather.Get("high").MustString()
	highSplit := strings.Split(high, " ")
	low := todayWeather.Get("low").MustString()
	lowSplit := strings.Split(low, " ")
	temperature := fmt.Sprintf("温度: %s/%s", lowSplit[1], highSplit[1])
	//风
	fx := todayWeather.Get("fx").MustString()
	fl := todayWeather.Get("fl").MustString()
	wind := fmt.Sprintf("%s : %s", fx, fl)
	//空气指数
	quality := js.Get("data").Get("quality").MustString()
	weatherQuality := fmt.Sprintf("空气质量 : %s", quality)
	if city == ""{
			city = "未知地区"
		}
	cityName := fmt.Sprintf("地区：%s", city)
	lastMsg := fmt.Sprintf("%s\n%s\n%s,%s。\n%s\n%s\n%s\n", todayTime, cityName, weatherType, notice, temperature, wind, weatherQuality)
	return lastMsg
}

func main() {
	loginIn()
	c := cron.New()
	var spec string
	var con conf
	conf := con.getConf()
	planTime := strings.Split(conf.Time, ":")
	hour := planTime[0]
	minute := planTime[1]
	second := planTime[2]
	if conf.Everyday == 1 {
		spec = fmt.Sprintf("%s %s %s * * ?", second, minute, hour)
	} else {
		planDate := strings.Split(conf.Date, "-")
		month := planDate[0]
		day := planDate[1]
		spec = fmt.Sprintf("%s %s %s %s %s ?", second, minute, hour, day, month)
	}
	c.AddFunc(spec, SendMsg)
	c.Start()
	SendMsg()
	time.Sleep(90 * time.Second)
}
