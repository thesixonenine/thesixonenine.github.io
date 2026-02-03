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
    "time"
)

const JSONFilePath = "../assets/data/endfield.json"

func main() {
	// 先登录 https://user.hypergryph.com/bindCharacters?game=endfield
	// 拿到 https://as.hypergryph.com/user/oauth2/v2/grant 中的 data.token, 可以用浏览器 network 搜索 grant 然后 copy response
	history := LocalHistoryJSONFile()
	if newest, b := findNewest(history); b {
		fmt.Println("最新抽卡记录: " + newest.String())
	}
	StoreGacha(history, Fetch(history, findU8Token("271646071", findOauthToken())))
}

// Fetch 根据 本地数据及token来查询所有新的数据
func Fetch(local []types.EndfieldGacha, token string) []types.EndfieldGacha {
    poolNewest := FindPoolNewest(local)
    poolNew := []types.EndfieldGacha{}
    for k, v := range map[string]string{"char": "角色", "weapon": "武器"} {
        if k == "char" {
            for m, n := range constant.EndfieldCharGachaType {
                newest := poolNewest[constant.EndfieldCharGachaTypeMap[m]]
                if newest.Exist() {
                    fmt.Println(fmt.Sprintf("本地[%s%s]池的最新抽卡记录[%s]", n, v, newest.String()))
                }
                fmt.Println(fmt.Sprintf("开始获取[%s%s]池", n, v))
                gachas := FetchWishes([]types.EndfieldGacha{}, newest, token, m, "")
                poolNew = append(poolNew, gachas...)
            }
        } else {
            newest := poolNewest[constant.WEPONBOX]
            if newest.Exist() {
                fmt.Println(fmt.Sprintf("本地[%s]池的最新抽卡记录[%s]", v, newest.String()))
            }
            fmt.Println(fmt.Sprintf("开始获取[%s]池", v))
            gachas := FetchWishes([]types.EndfieldGacha{}, newest, token, "", "")
            poolNew = append(poolNew, gachas...)
        }
    }
    return poolNew
}

// FindPoolNewest 查询每个池子中最新的记录
func FindPoolNewest(local []types.EndfieldGacha) map[string]types.EndfieldGacha{
    var EndfieldGachaPool = map[string]types.EndfieldGacha{
        constant.BEGINNER: {},
        constant.STANDARD: {},
        constant.SPECIAL:  {},
        constant.WEPONBOX: {},
    }
    found := 4
    SortReversedGacha(local)
    for _, gacha := range local {
        if found == 0 {
            break
        }
        if gacha.IsBeginner() && !EndfieldGachaPool[constant.BEGINNER].Exist() {
            EndfieldGachaPool[constant.BEGINNER] = gacha
            found--
            continue
        }
        if gacha.IsStandard() && !EndfieldGachaPool[constant.STANDARD].Exist() {
            EndfieldGachaPool[constant.STANDARD] = gacha
            found--
            continue
        }
        if gacha.IsSpecial() && !EndfieldGachaPool[constant.SPECIAL].Exist() {
            EndfieldGachaPool[constant.SPECIAL] = gacha
            found--
            continue
        }
        if gacha.IsWeapon() && !EndfieldGachaPool[constant.WEPONBOX].Exist() {
            EndfieldGachaPool[constant.WEPONBOX] = gacha
            found--
            continue
        }
    }
    return EndfieldGachaPool
}

// FetchWishes 通过递归调用自己来完成对四个池子中指定池子的分页查询
func FetchWishes(local []types.EndfieldGacha, newest types.EndfieldGacha, token string, poolType string, seqId string) []types.EndfieldGacha {
	endpoint := "https://ef-webview.hypergryph.com/api/record/%s?lang=zh-cn&server_id=1"
    gachaResp := types.EndfieldGachaResp{}
	if poolType != "" {
		// 角色
		gachaResp = FetchGacha(fmt.Sprintf(endpoint, "char"), token, poolType, seqId)
		if gachaResp.Code != 0 {
			log.Fatalf("角色查询失败,Code[%d]Msg[%s]", gachaResp.Code, gachaResp.Msg)
		}
	} else {
		// 武器
		gachaResp = FetchGacha(fmt.Sprintf(endpoint, "weapon"), token, "", seqId)
		if gachaResp.Code != 0 {
			log.Fatalf("武器查询失败,Code[%d]Msg[%s]", gachaResp.Code, gachaResp.Msg)
		}
	}
    for _, gacha := range gachaResp.Data.List {
        if gacha.GachaId() == newest.GachaId() {
            return local
        }
        local = append(local, gacha)
        seqId = gacha.SeqId
    }
    if gachaResp.Data.HasMore {
        return FetchWishes(local, newest, token, poolType, seqId)
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
    time.Sleep(1 * time.Second)
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
    l := len(newest)
    if l == 0 {
        fmt.Println("本次没有查询到新记录")
        return
    }
    fmt.Printf("开始合并本次查询到的新记录(%d条)并存储到文件\n", l)
    newest = append(newest, ist...)
	SortReversedGacha(newest)
	marshal, err := json.Marshal(newest)
	if err != nil {
		log.Fatalf("JSON序列化异常[%s]\n", err.Error())
		return
	}
	_ = os.WriteFile(JSONFilePath, utils.JSONIndent(marshal), 0600)
}

// 查询 oauthToken
func findOauthToken() string {
	multiLine := utils.ReadContinuedLinesStdin("请先对[https://as.hypergryph.com/user/oauth2/v2/grant]copy response然后粘贴并按Enter结束:")
	m := map[string]map[string]string{}
	_ = json.Unmarshal([]byte(multiLine), &m)
	return m["data"]["token"]
}

// 查询 u8Token
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
