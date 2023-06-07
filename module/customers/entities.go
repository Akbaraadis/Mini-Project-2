package customers

import "gorm.io/gorm"

type Actors struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	TokenKey string `json:"token_key"`
	RoleID   string `json:"role_id"`
	FlagAct  string `json:"flag_act"`
	FlagVer  string `json:"flag_ver"`
}
