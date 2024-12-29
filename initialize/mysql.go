package initialize

import (
	"RedPack/global"
	"RedPack/initialize/internal"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"os"
)

type Mysql struct {
	Username     string
	Password     string
	Path         string
	Port         string
	Dbname       string
	Config       string
	Prefix       string
	Engine       string
	Singular     bool
	MaxIdleConns int
	MaxOpenConns int
}

func (m *Mysql) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ":" + m.Port + ")/" + m.Dbname + "?" + m.Config
}

func (m *Mysql) SqlDsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ":" + m.Port + ")/"
}

func InitMysql() *gorm.DB {
	m := Mysql{
		Username:     "root",
		Password:     "123456",
		Port:         "3306",
		Dbname:       "redpack",
		Config:       "charset=utf8mb4&parseTime=true&loc=Asia%2fShanghai",
		Prefix:       "",
		Engine:       "",
		MaxIdleConns: 20,
		MaxOpenConns: 1000,
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
		if err = db.Use(dbresolver.Register(dbresolver.Config{
			Replicas:          []gorm.Dialector{mysql.New(replicasConfig)},
			Policy:            dbresolver.RandomPolicy{},
			TraceResolverMode: true,
		})); err != nil {
			fmt.Println("注册读写分离失败!")
			return nil
		}
		return db
	}
}

func CreateTables() {
	db := global.DB
	for i := 0; i < 10; i++ {
		redPackTableName := fmt.Sprintf("red_packs_%d", i)
		redPackRecordTableName := fmt.Sprintf("red_pack_records_%d", i)

		// 创建 red_packs 分片表
		err := db.Exec(fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				id bigint(20) NOT NULL AUTO_INCREMENT,
  				created_at datetime(3) DEFAULT NULL,
  				user_id varchar(191) DEFAULT NULL COMMENT '用户UUID',
    			total_amount double DEFAULT NULL COMMENT '总金额',
    			surplus_amount double DEFAULT NULL COMMENT '剩余金额',
    			total bigint(20) DEFAULT NULL COMMENT '红包总数',
    			surplus_total bigint(20) DEFAULT NULL COMMENT '剩余红包总数',
			    PRIMARY KEY (id),
  				KEY idx_red_pack_user_id (user_id)
		)`, redPackTableName)).Error
		if err != nil {
			fmt.Println("创建表出错!")
			os.Exit(0)
		}
		// 创建 red_pack_records 分片表
		err = db.Exec(fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				id bigint(20) NOT NULL AUTO_INCREMENT,
  				created_at datetime(3) DEFAULT NULL,
  				red_pack_id bigint(20) NOT NULL,
  				user_id varchar(191) NOT NULL,
  				amount double DEFAULT NULL,
				FOREIGN KEY (red_pack_id) REFERENCES %s(id),
			    PRIMARY KEY (id),
  				KEY idx_red_pack_record_red_pack_id (red_pack_id),
  				KEY idx_red_pack_record_user_id (user_id)
			)`, redPackRecordTableName, redPackTableName)).Error
		if err != nil {
			fmt.Println("创建表出错!")
			os.Exit(0)
		}
	}
}
