package types

import (
	"fmt"
	"strconv"
	"strings"
)

type ArkNightsData struct {
	Chars []ArkNightsChars `json:"chars"`
	Pool  string           `json:"pool"`
	Ts    int64            `json:"ts"`
}
type ArkNightsChars struct {
	IsNew  bool   `json:"isNew"`
	Name   string `json:"name"`
	Rarity int    `json:"rarity"`
}

type ArkNightsChar struct {
	PoolID   string `json:"poolId"`
	PoolName string `json:"poolName"`
	CharID   string `json:"charId"`
	CharName string `json:"charName"`
	Rarity   int    `json:"rarity"`
	IsNew    bool   `json:"isNew"`
	GachaTs  string `json:"gachaTs"`
	Pos      int    `json:"pos"`
}

func (c ArkNightsChar) String() string {
	n := ""
	if c.IsNew {
		n = "(新)"
	}
	return fmt.Sprintf("%s-%s"+n, c.CharName, strconv.Itoa(c.Rarity+1)+"星")
}

func WishCompare(w1 ArkNightsChar, w2 ArkNightsChar) int {
	its := w1.GachaTs
	jts := w2.GachaTs
	if its == "" {
		if jts == "" {
			return w1.Pos - w2.Pos
		} else {
			return -1
		}
	} else {
		if jts == "" {
			return 1
		}
	}
	if strings.EqualFold(its, jts) {
		return w1.Pos - w2.Pos
	}
	ii, _ := strconv.ParseInt(its[:len(its)-3], 10, 64)
	ji, _ := strconv.ParseInt(jts[:len(jts)-3], 10, 64)
	return int(ii - ji)
}

type ArkNightsGacha struct {
	Code int64 `json:"code"`
	Data struct {
		List    []ArkNightsChar `json:"list"`
		HasMore bool            `json:"hasMore"`
	} `json:"data"`
	Msg string `json:"msg"`
}
