package main

import (
	"fmt"
	"gosrc/internal/constant"
	"gosrc/internal/types"
	"gosrc/internal/utils"
    "slices"
    "sort"
	"strings"
)

func main() {
	fillGenshinImpact()
	fillHouse()
}

func fillHouse() {
	content := ""
	content = content + "\n\n" + buildHouseCost()
	content = content + "\n\n"
	filePath := "../content/post/house/index.md"
	_ = utils.KeepHeadAndAppendWithEndLine(filePath, 9, "## 网线布线", content)
}

func buildHouseCost() string {
	f := utils.ReadJSONFile[[]types.House]("../assets/data/house.json")
	var tableBuilder strings.Builder
	tableBuilder.WriteString("## 支出流水\n\n|费用名称|金额|状态|支出时间|Qing|Yang|\n|---|---|---|---|---|---|\n")
	qingTotal := 0
	yangTotal := 0
	for _, it := range f {
		tableBuilder.WriteString(fmt.Sprintf("|%s|%s|%s|%s|%s|%s|\n", it.Name, it.Amount, it.Status, it.PayTime, it.Qing, it.Yang))
		// 总计
		fen, err := utils.YuanToFen(string(it.Qing))
		if err != nil {
			fmt.Printf("YuanToFen[%s]Error\n", string(it.Qing))
			break
		}
		qingTotal += fen

		fen, err = utils.YuanToFen(string(it.Yang))
		if err != nil {
			fmt.Printf("YuanToFen[%s]Error\n", string(it.Yang))
			break
		}
		yangTotal += fen
	}

	tableBuilder.WriteString(fmt.Sprintf("|%s|%s|%s|%s|%s|%s|\n", "总计", utils.DivideHundred(qingTotal+yangTotal), "", "", utils.DivideHundred(qingTotal), utils.DivideHundred(yangTotal)))

	return tableBuilder.String()
}

func fillGenshinImpact() {
    content := "\n"
    content = content + "\n\n" + buildHK4E()
    content = content + "\n\n" + buildZZZ()
    content = content + "\n\n" + buildSR()
    content = content + "\n\n" + buildEndfield()
    content = content + "\n\n" + buildArkNightsV2()
    content = content + "\n\n" + buildArkNightsV1()
    filePath := "../content/post/genshin-impact/index.md"
    _ = utils.KeepHeadAndAppend(filePath, 9, content)
}

func buildArkNightsV2() string {
	f := utils.ReadJSONFile[[]types.ArkNightsChar]("../assets/data/arknightsV2.json")
	sort.Slice(f, func(i, j int) bool {
		return types.WishCompare(f[i], f[j]) < 0
	})

	type Pool struct {
		Id        string
		Name      string
		Count     int
		SixStars  []string // 存储六星记录的切片
		PityCount int      // 保底计数器
	}

	// 使用map提高查找效率
	poolMap := make(map[string]*Pool)
	// 记录池子出现的顺序
	poolOrder := []string{}

	for _, it := range f {
		pool, exists := poolMap[it.PoolID]
		if !exists {
			// 新池子初始化
			pool = &Pool{Id: it.PoolID, Name: it.PoolName}
			poolMap[it.PoolID] = pool
			poolOrder = append(poolOrder, it.PoolID)
		}

		pool.Count += 1

		pool.PityCount++
		if it.Rarity == 5 {
			// 格式化六星记录: 干员名(抽取序号)
			record := fmt.Sprintf("%s(%d)", it.CharName, pool.PityCount)
			pool.SixStars = append(pool.SixStars, record)
			// 重置保底计数器
			pool.PityCount = 0
		}
	}

	var tableBuilder strings.Builder
	tableBuilder.WriteString("## 明日方舟\n\n|池子|总抽取数量|六星|已抽|\n|---|---|---|---|\n")

    // 按出现顺序的倒序输出池子
    slices.Reverse(poolOrder)
	for _, name := range poolOrder {
		p := poolMap[name]
		// 将六星记录连接为逗号分隔的字符串
		sixStarStr := strings.Join(p.SixStars, ",")
		tableBuilder.WriteString(fmt.Sprintf("|%s|%d|%s|%d|\n", p.Name, p.Count, sixStarStr, p.PityCount))
	}

	return tableBuilder.String()
}

func buildArkNightsV1() string {
	f := utils.ReadJSONFile[[]types.ArkNightsData]("../assets/data/arknights.json")
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
	tableBuilder.WriteString("## 明日方舟V1\n\n|池子|总抽取数量|六星|已抽|\n|---|---|---|---|\n")

	// 按出现顺序输出池子
	for _, name := range poolOrder {
		p := poolMap[name]
		// 将六星记录连接为逗号分隔的字符串
		sixStarStr := strings.Join(p.SixStars, ",")
		tableBuilder.WriteString(fmt.Sprintf("|%s|%d|%s|%d|\n", p.Name, p.Count, sixStarStr, p.PityCount))
	}

	return tableBuilder.String()
}
func buildHK4E() string {
    f := utils.ReadJSONFile[map[string][]types.MiHoYoWish]("../assets/data/genshin-impact.json")
    s := ""
    ks := []string{}
    for k := range constant.HK4EGachaType {
        ks = append(ks, k)
    }
    sort.Strings(ks)
    for _, k := range ks {
        v := constant.HK4EGachaType[k]
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
    hsr := `## 原神

|池子|总抽取数量|五星|已抽|
|---|---|---|---|
`
    hsr = hsr + s
    return hsr
}
func buildSR() string {
	f := utils.ReadJSONFile[map[string][]types.MiHoYoWish]("../assets/data/star-rail.json")
	s := ""
	ks := []string{}
	for k := range constant.SRGachaType {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	for _, k := range ks {
		v := constant.SRGachaType[k]
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
	f := utils.ReadJSONFile[map[string][]types.MiHoYoWish]("../assets/data/zzz.json")
	s := ""
	ks := []string{}
	for k := range constant.ZZZGachaType {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	for _, k := range ks {
		v := constant.ZZZGachaType[k]
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
func buildEndfield() string {
    f := utils.ReadJSONFile[[]types.EndfieldGacha]("../assets/data/endfield.json")
    sort.Slice(f, func(i, j int) bool {
        return !f[i].TimeGt(f[j])
    })

    type Pool struct {
        Id        string
        Name      string
        Count     int
        SixStars  []string // 存储六星记录的切片
        PityCount int      // 保底计数器
    }

    // 使用map提高查找效率
    poolMap := make(map[string]*Pool)
    // 记录池子出现的顺序
    poolOrder := []string{}

    for _, it := range f {
        pool, exists := poolMap[it.PoolID]
        if !exists {
            // 新池子初始化
            pool = &Pool{Id: it.PoolID, Name: it.PoolName}
            poolMap[it.PoolID] = pool
            poolOrder = append(poolOrder, it.PoolID)
        }

        pool.Count += 1

        pool.PityCount++
        if it.Rarity == 6 {
            // 格式化六星记录: 干员名(抽取序号)
            record := fmt.Sprintf("%s(%d)", it.Name(), pool.PityCount)
            pool.SixStars = append(pool.SixStars, record)
            // 重置保底计数器
            pool.PityCount = 0
        }
    }

    var tableBuilder strings.Builder
    tableBuilder.WriteString("## 明日方舟：终末地\n\n|池子|总抽取数量|六星|已抽|\n|---|---|---|---|\n")

    // 按出现顺序的倒序输出池子
    slices.Reverse(poolOrder)
    for _, name := range poolOrder {
        p := poolMap[name]
        // 将六星记录连接为逗号分隔的字符串
        sixStarStr := strings.Join(p.SixStars, ",")
        tableBuilder.WriteString(fmt.Sprintf("|%s|%d|%s|%d|\n", p.Name, p.Count, sixStarStr, p.PityCount))
    }

    return tableBuilder.String()
}