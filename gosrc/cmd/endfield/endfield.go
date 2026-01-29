package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gosrc/internal/constant"
    "gosrc/internal/types"
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
                FetchCharGacha(uid, m, "")
			}
		} else {
			fmt.Println(fmt.Sprintf("开始获取[%s]池", v))
            FetchWeaponGacha(uid, "")
		}
	}
}

// FetchCharGacha 拉取角色数据
func FetchCharGacha(uid string, poolType string, seqId string) types.EndfieldCharGacha {
    client := &http.Client{}
    posStr := ""
    if seqId != "" {
        posStr = "&seq_id=" + seqId
    }
    ul := fmt.Sprintf("https://ef-webview.hypergryph.com/api/record/char?lang=zh-cn&token=%s&server_id=1&pool_type=%s%s", url.QueryEscape(uid), poolType, posStr)
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
    history := types.EndfieldCharGacha{}
    _ = json.Unmarshal(bodyText, &history)
    return history
}

// FetchWeaponGacha 拉取武器数据
func FetchWeaponGacha(uid string, seqId string) types.EndfieldWeaponGacha {
    client := &http.Client{}
    posStr := ""
    if seqId != "" {
        posStr = "&seq_id=" + seqId
    }
    ul := fmt.Sprintf("https://ef-webview.hypergryph.com/api/record/weapon?lang=zh-cn&token=%s&server_id=1%s", url.QueryEscape(uid), posStr)
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
    history := types.EndfieldWeaponGacha{}
    _ = json.Unmarshal(bodyText, &history)
    return history
}
