package model

import "database/sql"

// user_id INT PRIMARY KEY ,
// username VARCHAR(50) UNIQUE NOT NULL,
// `password` VARCHAR(50) NOT NULL,
// RealInfo VARCHAR(100)
// dynamic_attributes JSON (for storing region, data_level, etc)

type MysqlUser struct {
	UserID            string         `json:"user_id"`
	Username          string         `json:"username"`
	Password          string         `json:"password"`
	RealInfo          string         `json:"real_info"`
	UserType          string         `json:"user_type"`          // 从区块链获取的用户类型
	DynamicAttributes sql.NullString `json:"dynamic_attributes"` // JSON string like {"region":"Sichuan","data_level":"Internal"}
}
