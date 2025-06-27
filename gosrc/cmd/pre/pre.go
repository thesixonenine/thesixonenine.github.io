package main

import (
	"fmt"
	"gosrc/internal/constant"
	"gosrc/internal/types"
	"gosrc/internal/utils"
	"sort"
	"strings"
)

func main() {
	content := "\n"
	content = content + "\n\n" + buildZZZ()
	content = content + "\n\n" + buildSR()
	filePath := "../content/post/genshin-impact/index.md"
	_ = utils.KeepHeadAndAppend(filePath, 9, content)
}

func buildArkNights() string {

	return ""
}

func buildSR() string {
	f := utils.ReadJSONFile[map[string][]types.MiHoYoWish]("../static/data/star-rail-wish.json")
	s := ""
	for k, v := range constant.SRGachaType {
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
			if wish.RankType == "5" {
				fiveStr = fiveStr + fmt.Sprintf("%s(%d),", wish.Name, cnt)
				cnt = 0
			}
		}
		fiveStr = strings.TrimRight(fiveStr, ",")
		s = s + fmt.Sprintf("|%s|%d|%s|%d|\n", v, l, fiveStr, cnt)
	}
	zzz := `## 崩坏：星穹铁道

|池子|总抽取数量|五星|已抽|
|---|---|---|---|
`
	zzz = zzz + s
	return zzz
}

func buildZZZ() string {
	f := utils.ReadJSONFile[map[string][]types.MiHoYoWish]("../static/data/zzz.json")
	s := ""
	for k, v := range constant.ZZZGachaType {
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
