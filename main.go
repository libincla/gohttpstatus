package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
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
	fileobj, err := os.Open("./urls.txt")
	if err != nil {
		panic(err)
	}
	defer fileobj.Close()
	scanner := bufio.NewScanner(fileobj)

	for scanner.Scan() {
		var js jsonArray

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
		fmt.Println(resp.StatusCode)

	}

}