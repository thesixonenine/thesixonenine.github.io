package types

import (
    "fmt"
    "gosrc/internal/constant"
    "strconv"
    "strings"
)

// EndfieldGacha 角色/武器对象, 示例:
// {"poolId":"beginner","poolName":"启程寻访","charId":"chr_0029_pograni","charName":"骏卫","rarity":6,"isFree":false,"isNew":true,"gachaTs":"1769095519018","seqId":"40"}
// {"poolId":"standard","poolName":"基础寻访","charId":"chr_0023_antal","charName":"安塔尔","rarity":4,"isFree":false,"isNew":false,"gachaTs":"1769616274849","seqId":"207"}
// {"poolId":"special_1_0_1","poolName":"熔火灼痕","charId":"chr_0016_laevat","charName":"莱万汀","rarity":6,"isFree":false,"isNew":true,"gachaTs":"1769310873324","seqId":"192"}
// {"poolId":"weponbox_1_0_1","poolName":"熔铸申领","weaponId":"wpn_sword_0006","weaponName":"熔铸火焰","weaponType":"E_WeaponType_Sword","rarity":6,"isNew":true,"gachaTs":"1769832455301","seqId":"80"}
type EndfieldGacha struct {
	PoolID     string `json:"poolId"`
	PoolName   string `json:"poolName"`
	CharID     string `json:"charId"`
	CharName   string `json:"charName"`
	WeaponID   string `json:"weaponId"`
	WeaponName string `json:"weaponName"`
	WeaponType string `json:"weaponType"`
	Rarity     int    `json:"rarity"`
	IsNew      bool   `json:"isNew"`
	IsFree     bool   `json:"isFree"`
	GachaTs    string `json:"gachaTs"`
	SeqId      string `json:"seqId"`
}

func (receiver EndfieldGacha) String() string {
    if receiver.PoolName == "" {
        return ""
    }
    if receiver.IsWeapon() {
        return fmt.Sprintf("PoolName[%s]Weapon[%s][%s]", receiver.PoolName, receiver.WeaponName, strings.Repeat("⭐", receiver.Rarity))
    }
    return fmt.Sprintf("PoolName[%s]Char[%s][%s]", receiver.PoolName, receiver.CharName, strings.Repeat("⭐", receiver.Rarity))
}
func (receiver EndfieldGacha) Equal(other EndfieldGacha) bool {
	return receiver.PoolID == other.PoolID && receiver.GachaTs == other.GachaTs && receiver.SeqId == other.SeqId
}
func (receiver EndfieldGacha) TimeLt(other EndfieldGacha) bool {
    if receiver.GachaTs == other.GachaTs {
        is, _ := strconv.Atoi(receiver.SeqId)
        js, _ := strconv.Atoi(other.SeqId)
        return is < js
    }

    its, _ := strconv.ParseInt(receiver.GachaTs, 10, 64)
    jts, _ := strconv.ParseInt(other.GachaTs, 10, 64)
    return its < jts
}
func (receiver EndfieldGacha) TimeGt(other EndfieldGacha) bool {
    if receiver.GachaTs == other.GachaTs {
        is, _ := strconv.Atoi(receiver.SeqId)
        js, _ := strconv.Atoi(other.SeqId)
        return is > js
    }

    its, _ := strconv.ParseInt(receiver.GachaTs, 10, 64)
    jts, _ := strconv.ParseInt(other.GachaTs, 10, 64)
    return its > jts
}
func (receiver EndfieldGacha) GachaId() string {
	return receiver.PoolID + receiver.GachaTs + receiver.SeqId
}
func (receiver EndfieldGacha) MatchGachaType(gachaType string) bool {
    if receiver.IsWeapon() {
        return false
    }
    return strings.HasPrefix(receiver.PoolID, constant.EndfieldCharGachaTypeMap[gachaType])
}
func (receiver EndfieldGacha) IsBeginner() bool {
    return receiver.PoolID == "beginner"
}
func (receiver EndfieldGacha) IsStandard() bool {
    return receiver.PoolID == "standard"
}
func (receiver EndfieldGacha) IsSpecial() bool {
    return strings.HasPrefix(receiver.PoolID, "special")
}
func (receiver EndfieldGacha) IsWeapon() bool {
    return receiver.WeaponID != ""
}