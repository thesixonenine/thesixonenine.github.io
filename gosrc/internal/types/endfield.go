package types

// EndfieldChar 角色对象, 示例:
// {"poolId":"standard","poolName":"基础寻访","charId":"chr_0023_antal","charName":"安塔尔","rarity":4,"isFree":false,"isNew":false,"gachaTs":"1769616274849","seqId":"207"}
type EndfieldChar struct {
	PoolID   string `json:"poolId"`
	PoolName string `json:"poolName"`
	CharID   string `json:"charId"`
	CharName string `json:"charName"`
	Rarity   int    `json:"rarity"`
	IsNew    bool   `json:"isNew"`
	IsFree   bool   `json:"isFree"`
	GachaTs  string `json:"gachaTs"`
	SeqId    string `json:"seqId"`
}

type EndfieldWeapon struct {
	PoolID     string `json:"poolId"`
	PoolName   string `json:"poolName"`
	WeaponID   string `json:"weaponId"`
	WeaponName string `json:"weaponName"`
	WeaponType string `json:"weaponType"`
	Rarity     int    `json:"rarity"`
	IsNew      bool   `json:"isNew"`
	GachaTs    string `json:"gachaTs"`
	SeqId      string `json:"seqId"`
}
