package config

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Clay294/forum/flog"
	_ "github.com/go-sql-driver/mysql"
	ormMySQL "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	DefaultConfigFile = "./etc/config.toml"
)

var globalConf *Config

func GlobalConf() *Config {
	return globalConf
}

type Config struct {
	*MySQL      `toml:"mysql"`
	*HttpServer `toml:"http_server"`
}

type MySQL struct {
	*MySQLForumBase `toml:"mysql_forum_base"`
	*MySQLKeysBase  `toml:"mysql_keys_base"`
	*MySQLAdvanced  `toml:"mysql_advanced"`
}

type MySQLForumBase struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `env:"MYSQL_PASSWORD"`
	Database string `toml:"database"`
}

type MySQLKeysBase struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `env:"MYSQL_PASSWORD"`
	Database string `toml:"database"`
}

type MySQLAdvanced struct {
	MaxIdleConns    int           `toml:"max_idle_conns"`
	MaxOpenConns    int           `toml:"max_open_conns"`
	ConnMaxLifetime time.Duration `toml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `toml:"conn_max_idle_time"`
}

type HttpServer struct {
	*HttpServerBase     `toml:"http_server_base"`
	*HttpServerAdvanced `toml:"http_server_advanced"`
}

type HttpServerBase struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type HttpServerAdvanced struct {
	ReadTimeout  int `toml:"readtimeout"`
	WriteTimeout int `toml:"writetimeout"`
}

func (m *MySQL) CreateConnPools() (map[string]*sql.DB, error) {
	pools := make(map[string]*sql.DB)

	dsnForum := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		m.MySQLForumBase.Username,
		m.MySQLForumBase.Password,
		m.MySQLForumBase.Host,
		m.MySQLForumBase.Port,
		m.MySQLForumBase.Database,
	)

	dsnKeys := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		m.MySQLKeysBase.Username,
		m.MySQLKeysBase.Password,
		m.MySQLKeysBase.Host,
		m.MySQLKeysBase.Port,
		m.MySQLKeysBase.Database,
	)

	poolForum, err := sql.Open("mysql", dsnForum)
	if err != nil {
		flog.Flogger().Error().Msgf("creating mysql connection pool of forum failed: %s", err)
		return nil, fmt.Errorf("creating mysql connection pool of forum failed: %s", err)
	}

	poolKeys, err := sql.Open("mysql", dsnKeys)
	if err != nil {
		flog.Flogger().Error().Msgf("creating mysql connection pool of rsakeys failed: %s", err)
		return nil, fmt.Errorf("creating mysql connection pool of rsakeys failed: %s", err)
	}

	poolForum.SetMaxIdleConns(m.MaxIdleConns)
	poolForum.SetMaxOpenConns(m.MaxOpenConns)
	poolForum.SetConnMaxLifetime(time.Hour * m.ConnMaxLifetime)
	poolForum.SetConnMaxIdleTime(time.Minute * m.ConnMaxIdleTime)

	poolKeys.SetMaxIdleConns(m.MaxIdleConns)
	poolKeys.SetMaxOpenConns(m.MaxOpenConns)
	poolKeys.SetConnMaxLifetime(time.Hour * m.ConnMaxLifetime)
	poolKeys.SetConnMaxIdleTime(time.Minute * m.ConnMaxIdleTime)

	pools[m.MySQLForumBase.Database] = poolForum
	pools[m.MySQLKeysBase.Database] = poolKeys

	return pools, nil
}

func (m *MySQL) CreateConnByORM() (map[string]*gorm.DB, error) {
	var gdbs = make(map[string]*gorm.DB)

	pools, err := m.CreateConnPools()
	if err != nil {
		return nil, err
	}

	for name, pool := range pools {
		gdb, err := gorm.Open(
			ormMySQL.New(
				ormMySQL.Config{
					Conn: pool,
				},
			),
			&gorm.Config{
				PrepareStmt:    true,
				TranslateError: true,
			},
		)
		if err != nil {
			flog.Flogger().Error().Msgf("creating mysql connection to %s database failed: %s", name, err)
			return nil, fmt.Errorf("creating mysql connection to %s database failed: %s", name, err)
		}

		gdbs[name] = gdb
	}

	return gdbs, nil
}

func (hs *HttpServer) CreateAddr() string {
	return fmt.Sprintf("%s:%d", hs.Host, hs.Port)
}
