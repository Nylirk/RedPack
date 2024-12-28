package internal

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type _gorm struct{}

var Gorm = new(_gorm)

func (g *_gorm) Config(prefix string, singular bool) *gorm.Config {
	config := &gorm.Config{
		//设定建表规则 表前缀 表名是否复数
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   prefix,
			SingularTable: singular,
		},
		//禁止创建外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	return config
}
