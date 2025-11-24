package Configurator

import (
	"CarbonCreditMarketPlaceAuthAPI/Helper/DevMode"
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type ConfiguratorInterface interface {
	InitiateConfig() error
	InitiateConnections() error
}

type ConfiguratorStruct struct {
	DBDRIVER      string
	DBCONNSTRING  string
	RDBCONNSTRING string
	ADDRESS       string
	db            *sql.DB
	rdb           *redis.Client
	rdbOption     redis.Options
	Mode          int
}

type configParser struct {
	DbDriver      string `mapstructure="DBDRIVER"`
	DbConnString  string `mapstructure="DBCONNSTRING"`
	RDbConnString string `mapstructure="RDBCONNSTRING"`
	Address       string `mapstructure="ADDRESS"`
}

func NewConfiguration(Mode int) (ConfiguratorStruct, error) {
	conf := ConfiguratorStruct{}

	conf.Mode = Mode

	err := conf.InitiateConfig()

	if err != nil {
		panic(err)
	}

	err = conf.InitiateConnections()

	if err != nil {
		panic(err)
	}

	return conf, nil
}

func (Conf *ConfiguratorStruct) InitiateConfig() error {

	var configParser configParser

	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	//viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = viper.Unmarshal(&configParser)

	if err != nil {
		log.Fatal(err)
		return err
	}

	Conf.DBDRIVER = configParser.DbDriver
	Conf.DBCONNSTRING = configParser.DbConnString
	Conf.ADDRESS = configParser.Address
	Conf.RDBCONNSTRING = configParser.RDbConnString

	return nil

}

func (Conf *ConfiguratorStruct) InitiateConnections() error {

	err := Conf.InitiateDBConnection()

	if err != nil {
		return err
	}

	err = Conf.InitiateRDBConnection()

	if err != nil {
		return err
	}

	return nil

}

func (Conf *ConfiguratorStruct) InitiateDBConnection() error {

	db, err := sql.Open(Conf.DBDRIVER, Conf.DBCONNSTRING)
	if err != nil {
		return nil
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	Conf.db = db
	return nil
}

func (Conf *ConfiguratorStruct) InitiateRDBConnection() error {

	if Conf.Mode == DevMode.QA || Conf.Mode == DevMode.PROD {

		opt, err := redis.ParseURL(Conf.RDBCONNSTRING)

		if err != nil {
			return err
		}

		Conf.rdb = redis.NewClient(opt)

		return nil

	} else {

		Conf.rdbOption = redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password
			DB:       0,  // use default DB
			Protocol: 2,
		}

		Conf.rdb = redis.NewClient(&Conf.rdbOption)
		return nil
	}
}
