//go:build mysql || alldb
// +build mysql alldb

package database

import (
	"gorm.io/driver/mysql"
)

func init() {
	dialectors["mysql"] = mysql.Open
}
