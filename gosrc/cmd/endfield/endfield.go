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
    "os"
    "sort"
    "syscall"
)

const JSONFilePath = "../assets/data/endfield.json"

func main() {
	// 先登录 https://user.hypergryph.com/bindCharacters?game=endfield
	// 拿到 https://as.hypergryph.com/user/oauth2/v2/grant 中的 data.token, 可以用浏览器 network 搜索 grant 然后 copy response
	history := LocalHistoryJSONFile()
	if newest, b := findNewest(history); b {
		fmt.Println("Found newest endfield: " + newest.String())
	}
	StoreGacha(history, Fetch(history, findU8Token("271646071", findOauthToken())))
}

// Fetch 根据 本地数据及token来查询所有新的数据
func Fetch(local []types.EndfieldGacha, token string) []types.EndfieldGacha {
	SortReversedGacha(local)
	for k, v := range map[string]string{"char": "角色", "weapon": "武器"} {
		if k == "char" {
			for m, n := range constant.EndfieldCharGachaType {
				newest := types.EndfieldGacha{}
				for _, gacha := range local {
					if gacha.MatchGachaType(m) {
						newest = gacha
						break
					}
				}
				if newest.CharName != "" {
					fmt.Println(fmt.Sprintf("本地[%s%s]池的最新抽卡记录[%s]", n, v, newest.String()))
				}
				fmt.Println(fmt.Sprintf("开始获取[%s%s]池", n, v))
				FetchWishes([]types.EndfieldGacha{}, newest, token, m, "")
			}
		} else {
			newest := types.EndfieldGacha{}
			for _, gacha := range local {
				if gacha.IsWeapon() {
					newest = gacha
					break
				}
			}
			if newest.CharName != "" {
				fmt.Println(fmt.Sprintf("本地[%s]池的最新抽卡记录[%s]", v, newest.String()))
			}
			fmt.Println(fmt.Sprintf("开始获取[%s]池", v))
			FetchWishes([]types.EndfieldGacha{}, newest, token, "", "")
		}
	}
	fmt.Println(len(local))
	return local
}

// FetchWishes 通过递归调用自己来完成对四个池子中指定池子的分页查询
func FetchWishes(local []types.EndfieldGacha, newest types.EndfieldGacha, token string, poolType string, seqId string) []types.EndfieldGacha {
	endpoint := "https://ef-webview.hypergryph.com/api/record/%s?lang=zh-cn&server_id=1"
	if poolType != "" {
		// 角色
		gachaResp := FetchGacha(fmt.Sprintf(endpoint, "char"), token, poolType, seqId)
		if gachaResp.Code != 0 {
			log.Fatalf("角色查询失败,Code[%d]Msg[%s]", gachaResp.Code, gachaResp.Msg)
		}
		if len(local) > 0 {
			//
		}
		if gachaResp.Data.HasMore {
			// sid, _ := findMaxSeqId(gachaResp.Data.List)
			local = append(local, gachaResp.Data.List...)
			// FetchWishes(local, maxSeqId, token, poolType, sid)
		}
	} else {
		// 武器
		gachaResp := FetchGacha(fmt.Sprintf(endpoint, "weapon"), token, "", seqId)
		if gachaResp.Code != 0 {
			log.Fatalf("武器查询失败,Code[%d]Msg[%s]", gachaResp.Code, gachaResp.Msg)
		}
		if gachaResp.Data.HasMore {
			// FetchWishes(local, maxSeqId, token, poolType, seqId)
		}
	}
	return local
}

// findNewest 根据时间查询最新的记录
func findNewest(local []types.EndfieldGacha) (types.EndfieldGacha, bool) {
	if len(local) == 0 {
		return types.EndfieldGacha{}, false
	}
	SortReversedGacha(local)
	return local[0], true
}

// findOldest 根据时间查询最早的记录
func findOldest(local []types.EndfieldGacha) (types.EndfieldGacha, bool) {
	if len(local) == 0 {
		return types.EndfieldGacha{}, false
	}
	SortGacha(local)
	return local[0], true
}

// SortGacha 根据抽卡 ID (时间)进行正序
func SortGacha(ist []types.EndfieldGacha) {
	sort.Slice(ist, func(i, j int) bool {
		return ist[i].TimeLt(ist[j])
	})
}

// SortReversedGacha 根据抽卡 ID (时间)进行倒序
func SortReversedGacha(ist []types.EndfieldGacha) {
	sort.Slice(ist, func(i, j int) bool {
		return ist[i].TimeGt(ist[j])
	})
}

// FetchGacha 拉取数据
func FetchGacha(ul string, token string, poolType string, seqId string) types.EndfieldGachaResp {
	client := &http.Client{}
	ul += "&token=" + url.QueryEscape(token)
	if seqId != "" {
		ul += "&seq_id=" + seqId
	}
	if poolType != "" {
		ul += "&pool_type=" + poolType
	}
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
	history := types.EndfieldGachaResp{}
	_ = json.Unmarshal(bodyText, &history)
	return history
}

// LocalHistoryJSONFile 将本地抽卡历史 JSON 文件转为 [] 对象, 如果文件不存在则创建
func LocalHistoryJSONFile() []types.EndfieldGacha {
	historyFile, err := os.OpenFile(JSONFilePath, syscall.O_RDWR|syscall.O_CREAT, os.ModePerm)
	if err != nil {
		log.Fatalf("打开文件[%s]异常,err[%s]\n", JSONFilePath, err.Error())
	}
	defer func() {
		if err := historyFile.Close(); err != nil {
			log.Fatalf("关闭文件[%s]异常,err[%s]\n", JSONFilePath, err.Error())
		}
	}()
	historyFileContent, err := io.ReadAll(historyFile)
	if err != nil {
		log.Fatalf("读取文件[%s]异常,err[%s]\n", JSONFilePath, err.Error())
	}
	history := []types.EndfieldGacha{}
	_ = json.Unmarshal(historyFileContent, &history)
	return history
}

// StoreGacha 存储抽卡历史
func StoreGacha(ist []types.EndfieldGacha, newest []types.EndfieldGacha) {
	SortReversedGacha(ist)
	marshal, err := json.Marshal(ist)
	if err != nil {
		log.Fatalf("JSON序列化异常[%s]\n", err.Error())
		return
	}
	_ = os.WriteFile(JSONFilePath, utils.JSONIndent(marshal), 0600)
}
func findOauthToken() string {
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
