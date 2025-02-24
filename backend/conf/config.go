package conf

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	SQLite *SQLite `mapstructure:"SQLite"`
	App    *App    `mapstructure:"App"`
	MySQL  *MySQL  `mapstructure:"MySQL"`
}

type App struct {
	HttpHost string `json:"http_host" mapstructure:"HttpHost"`
	HttpPort int64  `json:"http_port" mapstructure:"HttpPort"`
}

type SQLite struct {
	// db文件的路径
	DBPath string `json:"db" mapstructure:"DBPath"`

	// 绑定一个单例对象，缓存一个对象
	lock sync.Mutex
	conn *gorm.DB
}

type MySQL struct {
	Host     string `json:"host" mapstructure:"Host"`
	Port     int64  `json:"port" mapstructure:"Port"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
	Database string `json:"database" mapstructure:"database"`

	// 绑定一个单例对象，缓存一个对象
	lock sync.Mutex
	conn *gorm.DB
}

// DefaultConfig 返回默认的配置
func DefaultConfig() *Config {
	return &Config{
		SQLite: &SQLite{
			DBPath: "template.db",
		},
		App: &App{
			HttpHost: "0.0.0.0",
			HttpPort: 80,
		},
	}
}

// HttpAddr 返回一个监听地址
func (a *App) HttpAddr() string {
	return fmt.Sprintf("%s:%d", a.HttpHost, a.HttpPort)
}

// GetConn 返回一个数据库连接对象
func (s *SQLite) GetConn() *gorm.DB {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.checkDBPath() // 校验配置文件数据库db文件是否存在

	// 如果conn为空，则进行获取
	if s.conn == nil {
		conn, err := gorm.Open(sqlite.Open(s.DBPath), &gorm.Config{})
		if err != nil {
			log.Printf("open sqlite db failed")
			panic(err)
		}
		s.conn = conn
	}

	return s.conn
}

// checkDBPath 检查sqlite db文件是否存在
func (s *SQLite) checkDBPath() {
	// 首先判断配置的 db 是绝对路径还是相对路径
	isAbs := filepath.IsAbs(s.DBPath)
	var path string
	if isAbs {
		path = s.DBPath
	} else {
		baseDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		path = filepath.Join(baseDir, s.DBPath)
	}
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		panic(err)
	}
}

func (m *MySQL) dsn() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.Username, m.Password, m.Host, m.Port, m.Database)
	return dsn

}

// GetConn 获取一个MySQL连接对象
func (m *MySQL) GetConn() *gorm.DB {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.conn == nil {
		conn, err := gorm.Open(mysql.Open(m.dsn()), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		m.conn = conn
	}

	return m.conn
}
