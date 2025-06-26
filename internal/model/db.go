package model

// orm连接参数配置参考
// https://help.aliyun.com/document_detail/281785.html?spm=5176.19908233.help.dexternal.3bd01450k0SUAX

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/rulego/rulego-server/config"

	"github.com/silenceper/log"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	TIMEZONE_CHINA  string = "Asia/Shanghai"
	SSLMODE_DISABLE string = "disable"
)

// LogMode 日志等级
const (
	LogModeSilent = `silent`
	LogModeError  = `error`
	LogModeWarn   = `warn`
	LogModeInfo   = `info`
)

var (
	// log                logger.ILogger
	ErrConfigJSONParse error = errors.New("configParseError")
)

var pgOnce sync.Once
var DBClient *DB

// StartDB 启动并初始化数据库
func StartDB(config config.Config) error {
	var err error
	pgOnce.Do(func() {

		client, err := NewDBWithStruct(&ORMConfig{
			User:            config.PostGre.User,
			Password:        config.PostGre.Password,
			Host:            config.PostGre.Host,
			Port:            config.PostGre.Port,
			DBname:          config.PostGre.DBname,
			MaxIdleConns:    4,
			MaxOpenConns:    4,
			ConnMaxLifeTime: 60 * time.Second,
			LogMode:         "info",
		})
		if err != nil {
			log.Infof("PG_CONGIG_ERR %s", err.Error())
			panic(err)
		}
		DBClient = client
	})
	return err
}

type ORMConfig struct {
	User            string        `json:"user"`            // 必须
	Password        string        `json:"password"`        // 必须
	Host            string        `json:"host"`            // 必须
	Port            int32         `json:"port"`            // 可选，默认5432
	DBname          string        `json:"dbname"`          // 必须
	SSlMode         string        `json:"sslmode"`         // 可选，默认SSLMODE_DISABLE
	MaxIdleConns    int           `json:"maxIdleConns"`    // 可选，默认1
	MaxOpenConns    int           `json:"maxOpenConns"`    // 可选，默认15
	ConnMaxLifeTime time.Duration `json:"connMaxLifeTime"` // 可选，默认time.Hour
	ConnMaxIdleTime time.Duration `json:"connMaxIdleTime"` // 可选，默认10 * time.Minute
	timeZone        string        `json:"-"`               // 固定Asia/Shanghai
	Callback        func(orm *DB) `json:"callback"`        // 初始化后需要执行的方法，比如注册callback之类的
	LogMode         string        `json:"logMode"`         // 可选，默认warn，可传入常量定义例如：LogModeSilent、LogModeError、LogModeWarn、LogModeInfo
}

type DB struct {
	Client *gorm.DB
	config *ORMConfig
}

// func init() {
// 	_log := logger.GetLogger()
// 	log = _log
// }

func LogModeString2GormLogLevel(mode string) gormLogger.LogLevel {
	switch mode {
	case LogModeSilent:
		return gormLogger.Silent
	case LogModeError:
		return gormLogger.Error
	case LogModeWarn:
		return gormLogger.Warn
	case LogModeInfo:
		return gormLogger.Info
	default:
		return gormLogger.Warn
	}
}

func NewDBWithStruct(config *ORMConfig) (*DB, error) {
	c := &DB{}
	defaultConfig := getDefaultConfig()
	if config.SSlMode == "" {
		config.SSlMode = defaultConfig.SSlMode
	}
	config.timeZone = defaultConfig.timeZone
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s timezone=%s", config.Host, config.User, config.Password, config.DBname, config.Port, config.SSlMode, config.timeZone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:         gormLogger.Default.LogMode(LogModeString2GormLogLevel(config.LogMode)),
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		log.Errorf("dbInitFailed%s%v", err.Error(), config)
		return nil, err
	}
	sqlDB, _ := db.DB()
	if config.ConnMaxLifeTime != 0 {
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
		sqlDB.SetConnMaxLifetime(config.ConnMaxLifeTime)
		sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)
	} else {
		sqlDB.SetConnMaxIdleTime(defaultConfig.ConnMaxIdleTime)
		sqlDB.SetConnMaxLifetime(defaultConfig.ConnMaxLifeTime)
		sqlDB.SetMaxIdleConns(defaultConfig.MaxIdleConns)
		sqlDB.SetMaxOpenConns(defaultConfig.MaxOpenConns)
	}
	c.Client = db
	c.config = config
	if config.Callback != nil {
		config.Callback(c)
	}
	return c, nil
}

func NewDBWithJson(c interface{}) (*DB, error) {
	value, err := json.Marshal(c)
	if err != nil {
		return nil, ErrConfigJSONParse
	}
	config := ORMConfig{}
	err = json.Unmarshal(value, &config)
	if err != nil {
		return nil, ErrConfigJSONParse
	}
	return NewDBWithStruct(&config)
}

func getDefaultConfig() *ORMConfig {
	return &ORMConfig{
		timeZone:        TIMEZONE_CHINA,
		SSlMode:         SSLMODE_DISABLE,
		MaxIdleConns:    1,
		MaxOpenConns:    15,
		ConnMaxLifeTime: time.Hour,
		ConnMaxIdleTime: 10 * time.Minute,
	}
}

// 检查链接是否健康
func (d *DB) GetHealthStatus(ctx context.Context) string {
	gormDB := d.Client.WithContext(ctx)
	sqlDB, err := gormDB.DB()
	if err != nil {
		return "unhealth"
	}
	// verifies a connection to the database is still alive
	err = sqlDB.Ping()
	if err != nil {
		return "unhealth"
	}
	err = gormDB.Raw(`select 1`).Error
	if err != nil {
		return "unhealth"
	}
	return "health"
}

// 获取目前数据库状态参数
func (d *DB) GetState() *sql.DBStats {
	db, err := d.Client.DB()
	if err != nil {
		log.Errorf("getSqlDBFailed%v", zap.Any("error", err))
		return nil
	}
	state := db.Stats()
	return &state
}

func (d *DB) Close() {
	db, _ := d.Client.DB()
	db.Close()
}
