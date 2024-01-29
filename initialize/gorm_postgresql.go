package initialize

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"main.go/global"
)

func GormPostgresql() *gorm.DB {
	// 获取全局配置中的 PostgreSQL 配置
	p := global.GVA_CONFIG.Postgresql
	if p.Database == "" {
		fmt.Println("找不到数据库....")
		return nil
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		p.Host,
		p.Port,
		p.Username,
		p.Database,
		p.Password,
	)
	fmt.Println(dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("连接字段错误....")
		return nil
	} else {
		fmt.Println("数据库postgresql连接成功！！！")
		return db
	}
}
