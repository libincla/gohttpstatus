package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/710leo/urlooker/dataobj"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type jsonArray struct {
	 Endpoint string `json:"endpoint"`
	 Metric string `json: "metric"`
	 Tags   string `json: "tags"`
	 Value  interface{} `json: "value"`
	 Timestamp int64 `json: "timestamp"`
	 Type string `json: "type"`
	 Step int `json: "step"`
}
func main() {
	wholeArray = make([]*jsonArray,0)

}
func getdata() jsonArray {
	var js jsonArray

	fileobj, err := os.Open("./urls.txt")
	if err != nil {
		panic(err)
	}
	defer fileobj.Close()
	scanner := bufio.NewScanner(fileobj)

	for scanner.Scan() {


		//判断endpoint
		ip, err := exec.Command("bash","-c","ifconfig").Output()
		if err != nil {
			fmt.Println(err.Error())
			js.Endpoint = ""
		}
		js.Endpoint = string(ip)
		js.Metric = "http.status.code"
		js.Type = "GAUGE"
		//取得step值
		filename := os.Args[0]
		step := strings.Split(filename,"_")[0]
		intstep, _ := strconv.Atoi(step)
		js.Step = intstep

		//取得timestamp值
		js.Timestamp = time.Now().Unix()

		x := scanner.Text()

		//判断每一条文本是否以http
		if ! strings.HasPrefix(x, "http") {
			var build strings.Builder
			build.WriteString("http://")
			build.WriteString(x)
			x = build.String()

		}
		resp, err := http.Get(x)
		if err != nil {
			fmt.Println(err.Error())
		}
		js.Value = resp.StatusCode

	}
	return js
}
//type MetricValue struct {
//	Endpoint  string      `json:"endpoint"`
//	Metric    string      `json:"metric"`
//	Tags      string      `json:"tags"`
//	Value     interface{} `json:"value"`
//	Timestamp int64       `json:"timestamp"`
//	Type      string      `json:"counterType"`
//	Step      int64       `json:"step"`
//}
//
//func PushFalcon(addr string, itemCheckedArray []*dataobj.CheckResult, ip string) {
//	pushDatas := make([]*MetricValue, 0)
//	for _, itemChecked := range itemCheckedArray {
//		tags := fmt.Sprintf("domain=%s,creator=%s,from=%s", itemChecked.Domain, itemChecked.Creator, ip)
//		if len(itemChecked.Tag) > 0 { //补充用户自定义tag
//			tags += "," + itemChecked.Tag
//		}
//
//		//url 状态
//		data := getMetric(itemChecked, "url_status", tags, itemChecked.Status)
//		pushDatas = append(pushDatas, &data)
//
//		//url 响应时间
//		data2 := getMetric(itemChecked, "url_resp_time", tags, int64(itemChecked.RespTime))
//		pushDatas = append(pushDatas, &data2)
//	}
//
//	err := pushData(addr, pushDatas)
//	if err != nil {
//		log.Println("push error", err)
//	}
//}
//
//func getMetric(item *dataobj.CheckResult, metric, tags string, value int64) MetricValue {
//	var data MetricValue
//	data.Endpoint = fmt.Sprintf("api_%d_%s", item.Sid, item.Domain)
//	if item.Endpoint != "" {
//		data.Endpoint = item.Endpoint
//	}
//
//	data.Timestamp = item.PushTime
//	data.Metric = metric
//	data.Type = "GAUGE"
//	data.Step = item.Step
//	data.Tags = tags
//	data.Value = value
//	return data
//}
//
//func pushData(addr string, data []*MetricValue) error {
//	d, err := json.Marshal(data)
//	if err != nil {
//		return err
//	}
//
//	resp, err := httplib.Post(addr).Header("Content-Type", "application/json").Body(d).String()
//	log.Printf("send:%s resp:%s\n", string(d), resp)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}