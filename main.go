package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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

var (
	wholeArray = make([]*jsonArray,0)
)
func main() {
	getvalue()
	//推数据
	d, err := json.Marshal(wholeArray)
	if err != nil {
		fmt.Println(err.Error())
	}
	//输出json到标准输出
	fmt.Println(string(d))

}
func getattribute() jsonArray {
	//函数获取Endpoint,Metric,Type,Step,timestep
	var js jsonArray
	//判断endpoint
	ip, err := exec.Command("bash","-c","ifconfig").Output()
	if err != nil {
		fmt.Println(err.Error())
		js.Endpoint = ""
	}
	js.Endpoint = string(ip)

	//获取Metric
	js.Metric = "http.status.code"
	js.Type = "GAUGE"
	//取得step
	step := strings.Split(os.Args[0],"_")[0]
	intstep, err := strconv.Atoi(step)
	if err != nil {
		intstep = 60
	}
	js.Step = intstep

	//取得timestamp值
	js.Timestamp = time.Now().Unix()

	return js
}
func getvalue()  {
	js := getattribute()
	fileobj, err := os.Open("./urls.txt")
	if err != nil {
		panic(err)
	}
	defer fileobj.Close()
	scanner := bufio.NewScanner(fileobj)

	for scanner.Scan() {
		x := scanner.Text()

		//判断每一条文本是否以http
		if ! strings.HasPrefix(x, "http") {
			var build strings.Builder
			build.WriteString("http://")
			build.WriteString(x)
			x = build.String()

		}
		//取得tags值
		domain := strings.Split(x,"/")[2]
		js.Tags = fmt.Sprintf("domain=%s", domain)
		resp, err := http.Get(x)
		if err != nil {
			fmt.Println(err.Error())
		}
		js.Value = resp.StatusCode

		wholeArray = append(wholeArray, &js)
	}

}
