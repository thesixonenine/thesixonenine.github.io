package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gosrc/internal/constant"
	"gosrc/internal/utils"
	"io"
	"log"
	"net/http"
	"net/url"
)

func main() {
	uid := "271646071"
	// 如何获取 u8_token
	// 先登录 https://user.hypergryph.com/
	// 拿到 https://as.hypergryph.com/user/oauth2/v2/grant 中的 data.token, 可以用浏览器的 copy response
    outhToken := findOuthToken()
	u8Token := findU8Token(uid, outhToken)
	Fetch(u8Token, u8Token)
}
func findOuthToken() string {
	multiLine := utils.ReadContinuedLinesStdin("请先对[https://as.hypergryph.com/user/oauth2/v2/grant]copy response然后粘贴并按Enter结束:")
	m := map[string]map[string]string{}
	_ = json.Unmarshal([]byte(multiLine), &m)
	return m["data"]["token"]
}

func findU8Token(uid string, token string) string {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://binding-api-account-prod.hypergryph.com/account/binding/v1/u8_token_by_uid", bytes.NewBufferString(fmt.Sprintf("{\"token\":\"%s\",\"uid\":\"%s\"}", token, uid)))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	m := make(map[string]map[string]string)
	_ = json.Unmarshal(bodyText, &m)
	return m["data"]["token"]
}

func Fetch(uid string, u8Token string) {
	for k, v := range map[string]string{"char": "角色", "weapon": "武器"} {
		if k == "char" {
			for m, n := range constant.EndfieldCharGachaType {
				fmt.Println(fmt.Sprintf("开始获取[%s%s]池", n, v))
				FetchGacha(uid, k, m, "")
			}
		} else {
			fmt.Println(fmt.Sprintf("开始获取[%s]池", v))
			FetchGacha(uid, k, "", "")
		}
	}
}

// FetchGacha 拉取数据
func FetchGacha(uid string, category string, charPoolType string, seqId string) {
	client := &http.Client{}
	posStr := ""
	if seqId != "" {
		posStr = "&seq_id=" + seqId
	}
	poolType := ""
	if category == "char" {
		poolType = "&pool_type=" + charPoolType
	}
	if category != "char" && category != "weapon" {
		log.Fatal("unsupport category: " + category)
	}
	ul := fmt.Sprintf("https://ef-webview.hypergryph.com/api/record/%s?lang=zh-cn&token=%s&server_id=1%s%s", category, url.QueryEscape(uid), poolType, posStr)
	parse, parseErr := url.Parse(ul)
	if parseErr != nil {
		log.Fatal(parseErr)
	}
	req, err := http.NewRequest("GET", parse.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bodyText))
}
