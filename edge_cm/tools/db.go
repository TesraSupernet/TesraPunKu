package tools

import (
	"fmt"
	"github.com/DOSNetwork/core/edge_cm/conf"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var dbErr error

// set table name
func (EdgeCmJob) TableName() string {
	return "edge_cm"
}

// set table name
func (CWA) TableName() string {
	return "cwa"
}

type Db struct {
	gormDB *gorm.DB
}

var db = Db{}

// init database
func (Db) LoadConf(conf *conf.Config) error {
	dbPath := fmt.Sprintf("%s%s", MakePath(conf.DB.Path), "edge_cm.db:locked.sqlite?cache=shared")
	db.gormDB, dbErr = gorm.Open("sqlite3", dbPath)
	if dbErr != nil {
		fmt.Println("db path: " + conf.DB.Path)
		panic(dbErr)
	}
	db.gormDB.DB().SetMaxOpenConns(1)
	db.gormDB.AutoMigrate(&EdgeCmJob{}, &CWA{}, &FailedMessage{})
	return nil
}

// create row
func Insert(value interface{}) *gorm.DB {
	return db.gormDB.Create(value)
}

// delete by example
func Delete(value interface{}, where ...interface{}) *gorm.DB {
	return db.gormDB.Delete(value, where...)
}

// update by example
func Updates(value, where, updateValue interface{}) *gorm.DB {
	return db.gormDB.Model(value).Where(where).UpdateColumns(updateValue)
}

// update by example
func Update(value interface{}) *gorm.DB {
	return db.gormDB.Save(value)
}

// select by example
func First(out interface{}, where ...interface{}) *gorm.DB {
	return db.gormDB.First(out, where...)

}

// select by example
func Find(out interface{}, where ...interface{}) *gorm.DB {
	return db.gormDB.Find(out, where...)
}

// aggregate search
// search table bb into AA struct e.g:
// Aggregate(&AA{}, &BB{}, "id <= ?", "Max(a) as max_a, Max(b) as max_b, Min(c) as min_c", 4)
// select Max(a) as max_a, Max(b) as max_b, Min(c) as min_c from bb where id <= 4;
func Aggregate(out, model, query interface{}, cond string, args ...interface{}) *gorm.DB {
	return db.gormDB.Model(model).Where(query, args...).Select(cond).Scan(out)
}

func Count(model, out, query interface{}, args ...interface{}) *gorm.DB {
	return db.gormDB.Model(model).Where(query, args).Count(out)
}
