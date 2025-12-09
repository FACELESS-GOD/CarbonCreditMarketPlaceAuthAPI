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
	dBDRIVER      string
	dBCONNSTRING  string
	rDBCONNSTRING string
	aDDRESS       string
	DB            *sql.DB
	RDB           *redis.Client
	rdbOption     redis.Options
	Mode          int
	TxOption      sql.TxOptions
	JwtSecretKey  string
}

type configParser struct {
	DbDriver      string `mapstructure:"DBDRIVER"`
	DbConnString  string `mapstructure:"DBCONNSTRING"`
	RDbConnString string `mapstructure:"RDBCONNSTRING"`
	Address       string `mapstructure:"ADDRESS"`
	JWTKEY        string `mapstructure:"JWTKEY"`
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

	txOption := sql.TxOptions{
		Isolation: sql.LevelSerializable,
	}

	conf.TxOption = txOption

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

	Conf.dBDRIVER = configParser.DbDriver
	Conf.dBCONNSTRING = configParser.DbConnString
	Conf.aDDRESS = configParser.Address
	Conf.rDBCONNSTRING = configParser.RDbConnString
	Conf.JwtSecretKey = configParser.JWTKEY

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

	db, err := sql.Open(Conf.dBDRIVER, Conf.dBCONNSTRING)
	if err != nil {
		return nil
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	Conf.DB = db
	return nil
}

func (Conf *ConfiguratorStruct) InitiateRDBConnection() error {

	if Conf.Mode == DevMode.QA || Conf.Mode == DevMode.PROD {

		opt, err := redis.ParseURL(Conf.rDBCONNSTRING)

		if err != nil {
			return err
		}

		Conf.RDB = redis.NewClient(opt)

		return nil

	} else {

		Conf.rdbOption = redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password
			DB:       0,  // use default DB
			Protocol: 2,
		}

		Conf.RDB = redis.NewClient(&Conf.rdbOption)
		return nil
	}
}
