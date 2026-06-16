package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)
type Config struct{ Server ServerConfig; Database DatabaseConfig; Redis RedisConfig }
type ServerConfig struct{ Host string; Port int; ReadTimeout time.Duration; WriteTimeout time.Duration }
type DatabaseConfig struct{ Host string; Port int; User string; Password string; DBName string; SSLMode string; MaxOpenConns int; MaxIdleConns int; ConnMaxLifetime time.Duration }
type RedisConfig struct{ Host string; Port int; Password string; DB int }
func(d*DatabaseConfig)DSN()string{return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",d.Host,d.Port,d.User,d.Password,d.DBName,d.SSLMode)}
func Load(p string)(*Config,error){
v:=viper.New();v.SetConfigName("config");v.SetConfigType("yaml");if p!=""{v.AddConfigPath(p)};v.AddConfigPath(".");v.AddConfigPath("./config")
v.SetEnvPrefix("LEAP_CONFIG");v.SetEnvKeyReplacer(strings.NewReplacer(".","_"));v.AutomaticEnv()
v.SetDefault("server.host","0.0.0.0");v.SetDefault("server.port",8014);v.SetDefault("server.read_timeout",30*time.Second);v.SetDefault("server.write_timeout",30*time.Second)
v.SetDefault("database.host","localhost");v.SetDefault("database.port",5432);v.SetDefault("database.user","postgres");v.SetDefault("database.password","postgres");v.SetDefault("database.db_name","config_db");v.SetDefault("database.sslmode","disable");v.SetDefault("database.max_open_conns",50);v.SetDefault("database.max_idle_conns",10);v.SetDefault("database.conn_max_lifetime",30*time.Minute)
v.SetDefault("redis.host","localhost");v.SetDefault("redis.port",6379);v.SetDefault("redis.password","");v.SetDefault("redis.db",0)
if e:=v.ReadInConfig();e!=nil{if _,ok:=e.(viper.ConfigFileNotFoundError);!ok{return nil,fmt.Errorf("读取配置失败:%w",e)}}
var cfg Config;if e:=v.Unmarshal(&cfg);e!=nil{return nil,fmt.Errorf("解析配置失败:%w",e)};return &cfg,nil
}
