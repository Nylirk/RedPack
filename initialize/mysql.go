package initialize

import (
	"RedPack/global"
	"RedPack/initialize/internal"
	"RedPack/model/system"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"os"
)

type Mysql struct {
	Username string
	Password string
	Path     string
	Port     string
	Dbname   string
	Config   string
	Prefix   string
	Engine   string
	Singular bool
}

func (m *Mysql) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ":" + m.Port + ")/" + m.Dbname + "?" + m.Config
}

func (m *Mysql) SqlDsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ":" + m.Port + ")/"
}

func InitMysql() *gorm.DB {
	m := Mysql{
		Username: "root",
		Password: "123456",
		Port:     "3306",
		Dbname:   "redpack",
		Config:   "charset=utf8mb4&parseTime=true&loc=Asia%2fShanghai",
		Prefix:   "",
		Engine:   "",
	}
	if m.Dbname == "" {
		return nil
	}
	sources := m
	sources.Path = "1.14.205.162"
	replicas := m
	replicas.Path = "121.40.25.175"
	dbm, err := sql.Open("mysql", sources.SqlDsn())
	if err != nil {
		fmt.Println("连接主库失败!")
		return nil
	}
	dbs, err := sql.Open("mysql", replicas.SqlDsn())
	if err != nil {
		fmt.Println("连接从库失败!")
		return nil
	}
	createSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", m.Dbname)
	_, err = dbm.Exec(createSql)
	if err != nil {
		fmt.Println("创建主数据库失败!", err)
		return nil
	}
	_, err = dbs.Exec(createSql)
	if err != nil {
		fmt.Println("创建从数据库失败!", err)
		return nil
	}
	sourcesConfig := mysql.Config{
		DSN:                       sources.Dsn(), // DSN data source name
		DefaultStringSize:         191,           // string 类型字段的默认长度
		SkipInitializeWithVersion: false,         // 根据版本自动配置
	}
	replicasConfig := mysql.Config{
		DSN:                       replicas.Dsn(), // DSN data source name
		DefaultStringSize:         191,            // string 类型字段的默认长度
		SkipInitializeWithVersion: false,          // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(sourcesConfig), internal.Gorm.Config(m.Prefix, m.Singular)); err != nil {
		return nil
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE="+m.Engine)
		if err := db.Use(dbresolver.Register(dbresolver.Config{
			Replicas:          []gorm.Dialector{mysql.New(replicasConfig)},
			Policy:            dbresolver.RandomPolicy{},
			TraceResolverMode: true,
		})); err != nil {
			fmt.Println("注册读写分离失败!", err)
			return nil
		}
		return db
	}
}

func CreateTables() {
	db := global.DB
	err := db.AutoMigrate(
		system.RedPack{},
		system.RedPackRecord{},
	)
	if err != nil {
		fmt.Println("创建主库表失败!")
		os.Exit(0)
	}
}
