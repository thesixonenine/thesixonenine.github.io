package main

import (
	"encoding/json"
	"fmt"
	"gosrc/internal/types"
	"gosrc/internal/utils"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	parseCurl "github.com/thesixonenine/parse-curl"
)

const JSONFilePath = "../assets/data/arknightsV2.json"

var roleToken = ""
var accountToken = ""
var cookie = ""
var cates []types.ArkNightsCategory

// 收集每次从官网获取的抽卡数据, 然后合并到本地JSON文件中.
func main() {
	cUrl, _ := ExtractCUrlBash()
	if cUrl == nil {
		return
	}
	c, err := url.ParseRequestURI(cUrl.Url)
	if err != nil {
		fmt.Println("未获取到uid, 请重新粘贴cURL")
		return
	}
	roleToken = cUrl.Header["x-role-token"]
	accountToken = cUrl.Header["x-account-token"]
	cookie = cUrl.Header["Cookie"]
	cates = FetchCate(c.Query().Get("uid")).List
	if len(cates) == 0 {
		fmt.Println("无抽卡记录")
		return
	}
	UpdateGacha(c.Query().Get("uid"))
}

func ExtractCUrlBash() (*parseCurl.Request, string) {
	multiLine := utils.ReadContinuedLinesStdin("请粘贴cURL命令并按Enter结束:")
	curl, err := parseCurl.Parse(multiLine)
	if err != nil {
		log.Printf("\n[%v]解析错误: %v\n", multiLine, err.Error())
		return nil, multiLine
	}
	return curl, multiLine
}

func UpdateGacha(uid string) {
	history := LocalHistory()
	var maxWish = types.ArkNightsChar{}
	if len(history) > 0 {
		maxWish = history[0]
		for _, wish := range history {
			if types.WishCompare(maxWish, wish) < 0 {
				maxWish = wish
			}
		}
		log.Printf("本地最新抽卡记录: [%s]%s\n", maxWish.PoolName, maxWish.String())
	}

	hasNewGacha := false
	for _, cate := range cates {
		gachaTs := ""
		pos := -1
		page := 1
		for {
			cateName := strings.ReplaceAll(cate.Name, "\n", " ")
			pageStr := "[" + cateName + "]" + "第" + strconv.Itoa(page) + "页"
			curGacha := FetchGacha(uid, cate.ID, gachaTs, pos)
			if curGacha.Code != 0 {
				log.Printf(pageStr+"查询失败,code[%d][%s]\n", curGacha.Code, curGacha.Msg)
				break
			}
			listLen := len(curGacha.Data.List)
			if listLen == 0 {
				log.Println(pageStr + "无数据,结束查询")
				break
			} else {
				log.Printf(pageStr+"共%d条数据\n", listLen)
			}
			hasNew := false
			for _, wish := range curGacha.Data.List {
				gachaTs = wish.GachaTs
				pos = wish.Pos
				if types.WishCompare(maxWish, wish) < 0 {
					log.Printf("查询到新的抽卡记录,%s\n", wish.String())
					hasNew = true
					hasNewGacha = true
					history = append(history, wish)
				}
			}
			if !hasNew {
				log.Println("[" + cateName + "]" + "卡池已获取所有最新的抽卡记录")
				break
			}
			if listLen < 10 {
				break
			}
			time.Sleep(1 * time.Second)
			page++
		}
	}
	if hasNewGacha {
		log.Println("开始同步抽卡记录")
		StoreWishes(history)
	}
}

func LocalHistory() []types.ArkNightsChar {
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
	var history []types.ArkNightsChar
	_ = json.Unmarshal(historyFileContent, &history)
	return history
}

// FetchCate 拉取分类
func FetchCate(uid string) types.ArkNightsCate {
	client := &http.Client{}
	parse, parseErr := url.Parse(fmt.Sprintf("https://ak.hypergryph.com/user/api/inquiry/gacha/cate?uid=%s", uid))
	if parseErr != nil {
		log.Fatal(parseErr)
	}
	req, err := http.NewRequest("GET", parse.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Cookie", cookie)
	req.Header.Set("X-Account-Token", accountToken)
	req.Header.Set("X-Role-Token", roleToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	cate := types.ArkNightsCate{}
	_ = json.Unmarshal(bodyText, &cate)
	return cate
}

// FetchGacha 拉取数据
func FetchGacha(uid string, category string, gachaTs string, pos int) types.ArkNightsGacha {
	client := &http.Client{}
	gachaTsStr := ""
	if gachaTs != "" {
		gachaTsStr = "&gachaTs=" + gachaTs
	}
	posStr := ""
	if pos != -1 {
		posStr = "&pos=" + strconv.Itoa(pos)
	}
	// https://ak.hypergryph.com/user/api/inquiry/gacha/cate?uid=
	parse, parseErr := url.Parse(fmt.Sprintf("https://ak.hypergryph.com/user/api/inquiry/gacha/history?uid=%s&category=%s%s%s&size=10", uid, category, posStr, gachaTsStr))
	if parseErr != nil {
		log.Fatal(parseErr)
	}
	req, err := http.NewRequest("GET", parse.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Cookie", cookie)
	req.Header.Set("X-Account-Token", accountToken)
	req.Header.Set("X-Role-Token", roleToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	history := types.ArkNightsGacha{}
	_ = json.Unmarshal(bodyText, &history)
	return history
}

// StoreWishes 存储抽卡历史
func StoreWishes(wishMap []types.ArkNightsChar) {
	wishMap = SortWishMap(wishMap)
	marshal, err := json.Marshal(wishMap)
	if err != nil {
		log.Fatalf("JSON序列化异常[%s]\n", err.Error())
		return
	}
	_ = os.WriteFile(JSONFilePath, utils.JSONIndent(marshal), 0600)
}

// SortWishMap 排序
func SortWishMap(ist []types.ArkNightsChar) []types.ArkNightsChar {
	sort.Slice(ist, func(i, j int) bool {
		its := ist[i].GachaTs
		jts := ist[j].GachaTs
		if strings.EqualFold(its, jts) {
			return ist[i].Pos < ist[j].Pos
		}
		ii, _ := strconv.ParseInt(its[:len(its)-3], 10, 64)
		ji, _ := strconv.ParseInt(jts[:len(jts)-3], 10, 64)
		return ii < ji
	})
	return ist
}
