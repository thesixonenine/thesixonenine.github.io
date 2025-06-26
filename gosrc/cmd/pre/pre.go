package main

import (
	"bufio"
	"fmt"
	"gosrc/internal/constant"
	"gosrc/internal/utils"
	"os"
	"sort"
	"strings"
)

type HKRPGWish struct {
	Uid       string `json:"uid"`
	GachaId   string `json:"gacha_id"`
	GachaType string `json:"gacha_type"`
	ItemId    string `json:"item_id"`
	Count     string `json:"count"`
	Time      string `json:"time"`
	Name      string `json:"name"`
	Lang      string `json:"lang"`
	ItemType  string `json:"item_type"`
	RankType  string `json:"rank_type"`
	Id        string `json:"id"`
}

func main() {
	zzz := buildZZZ()
	_ = fixGenshinImpact(9, zzz)
}

func buildZZZ() string {
	f := utils.ReadJSONFile[map[string][]HKRPGWish]("../static/data/zzz.json")
	s := ""
	for k, v := range constant.GachaTypeMap {
		wishes := f[k]
		l := len(wishes)
		// 排序
		sort.Slice(wishes, func(i, j int) bool {
			return wishes[i].Id < wishes[j].Id
		})
		// 统计
		fiveStr := ""
		cnt := 0
		for _, wish := range wishes {
			cnt++
			if wish.RankType == "4" {
				fiveStr = fiveStr + fmt.Sprintf("%s(%d),", wish.Name, cnt)
				cnt = 0
			}
		}
		fiveStr = strings.TrimRight(fiveStr, ",")
		s = s + fmt.Sprintf("|%s|%d|%s|%d|\n", v, l, fiveStr, cnt)
	}
	zzz := `## 绝区零

|池子|总抽取数量|五星|已抽|
|---|---|---|---|
`
	zzz = zzz + s
	return zzz
}

// fixGenshinImpact 保留开头几行并追加指定内容
func fixGenshinImpact(keepHeadLine int, appendContent string) error {
	filePath := "../content/post/genshin-impact/index.md"
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	var preservedLines []string
	scanner := bufio.NewScanner(file)
	for i := 0; i < keepHeadLine && scanner.Scan(); i++ {
		preservedLines = append(preservedLines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for _, line := range preservedLines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	if _, err := writer.WriteString(appendContent); err != nil {
		return err
	}
	return writer.Flush()
}
