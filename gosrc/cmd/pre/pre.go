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
        // 追加已有的数据
        if k == "100" {
            l += 20
            cnt += 20
        } else if k == "200" {
            l += 803
            fiveStr = "琴(48),狼的末路(77),天空之傲(9),刻晴(77),狼的末路(78),七七(76),狼的末路(74),莫娜(77),风鹰剑(42),迪希雅(79),天空之傲(76),迪希雅(79)," + strings.Replace(fiveStr, "天空之翼(67)", "天空之翼(78)", 1)
        } else if k == "301" {
            l += 4280
            fiveStr = "刻晴(79),琴(75),钟离(75),莫娜(84),神里绫华(55),雷电将军(63),珊瑚宫心海(30),胡桃(81),迪卢克(77),胡桃(76),莫娜(80),胡桃(77),七七(76),胡桃(19),胡桃(4),胡桃(74),优菈(77),迪卢克(75),甘雨(13),莫娜(77),八重神子(76),刻晴(79),神里绫华(76),七七(76),神里绫华(5),神里绫华(28),神里绫华(76),神里绫华(80),神里绫华(77),琴(78),妮露(74),阿贝多(75),七七(77),纳西妲(76),提纳里(79),雷电将军(56),雷电将军(76),迪卢克(28),夜兰(78),夜兰(76),夜兰(64),迪卢克(49),夜兰(75),夜兰(79),申鹤(80),迪卢克(75),枫原万叶(75),温迪(76),刻晴(78),芙宁娜(76),迪希雅(55),芙宁娜(77),娜维娅(77),琴(77),纳西妲(2),那维莱特(80),阿蕾奇诺(75),流浪者(30),希格雯(80),芙宁娜(9),刻晴(81),纳西妲(79),恰斯卡(77),那维莱特(81)," + strings.Replace(fiveStr, "茜特菈莉(20)", "茜特菈莉(80)", 1)
        } else if k == "302" {
            l += 673
            fiveStr = "天空之脊(14),天空之脊(67),雾切之回光(68),天空之脊(67),尘世之锁(14),护摩之杖(39),护摩之杖(6),天空之卷(67),薙草之稻光(69),护摩之杖(5),若水(52),四风原典(65),裁断(67),星鹫赤羽(66)," + fiveStr
            cnt += 7
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