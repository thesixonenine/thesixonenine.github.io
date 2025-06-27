package types

import (
	"fmt"
)

type MiHoYoWish struct {
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

func (wish MiHoYoWish) String() string {
	return fmt.Sprintf("%s %s %s", wish.Time, wish.ItemType, wish.Name)
}
func (wish MiHoYoWish) Fmt(m map[string]string) string {
	return fmt.Sprintf("%s %s %s", wish.Time, m[wish.GachaType], wish.Name)
}

type Page[T any] struct {
	Retcode int    `json:"retcode"`
	Message string `json:"message"`
	Data    struct {
		Page           int    `json:"page"`
		Size           int    `json:"size"`
		List           []T    `json:"list"`
		Region         string `json:"region"`
		RegionTimeZone int    `json:"region_time_zone"`
	} `json:"data"`
}
