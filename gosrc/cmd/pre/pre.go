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
	content = content + "\n\n" + buildArkNights()
	content = content + "\n\n" + buildZZZ()
	content = content + "\n\n" + buildSR()
	filePath := "../content/post/genshin-impact/index.md"
	_ = utils.KeepHeadAndAppend(filePath, 9, content)
}

func buildArkNights() string {
	f := utils.ReadJSONFile[[]types.Data]("../static/data/arknights.json")
	sort.Slice(f, func(i, j int) bool {
		return f[i].Ts < f[j].Ts
	})

	type Pool struct {
		Name      string
		Count     int
		SixStars  []string // 存储六星记录的切片
		PityCount int      // 保底计数器
	}

	// 使用map提高查找效率, 同时保持顺序
	poolMap := make(map[string]*Pool)
	// 记录池子出现的顺序
	poolOrder := []string{}

	for _, it := range f {
		pool, exists := poolMap[it.Pool]
		if !exists {
			// 新池子初始化
			pool = &Pool{Name: it.Pool}
			poolMap[it.Pool] = pool
			poolOrder = append(poolOrder, it.Pool)
		}

		pool.Count += len(it.Chars)
		for _, char := range it.Chars {
			pool.PityCount++
			if char.Rarity == 5 {
				// 格式化六星记录: 干员名(抽取序号)
				record := fmt.Sprintf("%s(%d)", char.Name, pool.PityCount)
				pool.SixStars = append(pool.SixStars, record)
				// 重置保底计数器
				pool.PityCount = 0
			}
		}
	}

	var tableBuilder strings.Builder
	tableBuilder.WriteString("## 明日方舟\n\n|池子|总抽取数量|六星|已抽|\n|---|---|---|---|\n")

	// 按出现顺序输出池子
	for _, name := range poolOrder {
		p := poolMap[name]
		// 将六星记录连接为逗号分隔的字符串
		sixStarStr := strings.Join(p.SixStars, ",")
		tableBuilder.WriteString(fmt.Sprintf("|%s|%d|%s|%d|\n", p.Name, p.Count, sixStarStr, p.PityCount))
	}

	return tableBuilder.String()
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
	hsr := `## 崩坏：星穹铁道

|池子|总抽取数量|五星|已抽|
|---|---|---|---|
`
	hsr = hsr + s
	return hsr
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
