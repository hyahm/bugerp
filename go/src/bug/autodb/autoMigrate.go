package autodb

import (
	"fmt"
	"gaconfig"
	"gadb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func InitDb() {
	connstring := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4",
		gaconfig.ReadString("mysqluser"),
		gaconfig.ReadString("mysqlpwd"),
		gaconfig.ReadString("mysqlhost"),
		gaconfig.ReadInt("mysqlport"),
		gaconfig.ReadString("mysqldb"),
	)
	db, err := gorm.Open(gaconfig.ReadString("sqldriver"), connstring)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// 忽略s
	db.SingularTable(true)
	// 自动迁移模式
	err = db.AutoMigrate(&Apilist{},
		&Status{},
		&Projectname{},
		&Jobs{},
		&User{},
		&Header{},
		&Version{},
		&Rolegroup{},
		&Restfulname{},
		&Apiproject{},
		&Usergroup{},
		&Email{},
		&Statusgroup{},
		&Types{},
		&Level{},
		&Bugs{},
		&Options{},
		&Defaultvalue{},
		&Importants{},
		&Sharefile{},
		&Log{},
		&Headerlist{},
		&Informations{},
		&Environment{},
		&Roles{}).Error
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("================= sync database success ========================")
	mysqldb, err := gadb.NewSqlConfig().ConnDB()
	if err != nil {
		log.Fatal(err)
	}
	defer mysqldb.Db.Close()
	var count int
	err = mysqldb.GetOne("select count(id) from user").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		_, err = mysqldb.Insert("insert into user(nickname,password,email,createtime,realname,level) values(?,?,?,?,?,?)",
			"admin", "69ad5117e7553ecfa7f918a223426dd8da08a57f", "admin@qq.com", time.Now().Unix(), "admin", 0)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = mysqldb.GetOne("select count(status) from defaultvalue").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count != 1 {
		//清空表
		_, err = mysqldb.Update("truncate defaultvalue")
		if err != nil {
			log.Fatal(err)
		}
		_, err = mysqldb.Insert("insert into defaultvalue(status,important,level) values(0,0,0)")
		if err != nil {
			log.Fatal(err)
		}
	}

	// 角色表
	//清空表
	
	_, err = mysqldb.Update("truncate roles")
	if err != nil {
		log.Fatal(err)
	}
	_, err = mysqldb.Insert("insert into roles(role) values(?),(?),(?),(?),(?),(?),(?),(?),(?),(?),(?)",
		"version", "project", "env", "status", "log", "statusgroup", "rolegroup", "important", "level",
		"position", "usergroup")
	if err != nil {
		log.Fatal(err)
	}

	// 增加类型
	err = mysqldb.GetOne("select count(id) from  types ").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		_, err = mysqldb.Insert("insert into types(`name`,`default`) values(?,?),(?,?),(?,?),(?,?),(?,?),(?,?),(?,?),(?,?),(?,?),(?,?),(?,?),(?,?),(?,?)",
			"int", "0", "string", "", "int64", "0", "double", "0.00", "bool", "false", "int8", "0", "int16", "0", "uint8", "0",
			"uint16", "0", "uint32", "0", "uint64", "0", "float32", "0", "float64", "0")
		if err != nil {
			log.Fatal(err)
		}
	}
	if gaconfig.ReadBool("apihelp") {
		createapi(mysqldb)
	}
	fmt.Println("================= check tables success ========================")
}
