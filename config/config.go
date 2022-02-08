package config

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"gorm.io/gorm/logger"

	"github.com/anik4good/go_email2sms/models"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	DB_USERNAME string `yaml:"DB_USERNAME"`
	DB_PASSWORD string `yaml:"DB_PASSWORD"`
	DB_IP       string `yaml:"DB_IP"`
	DB_PORT     string `yaml:"DB_PORT"`
	DB_NAME     string `yaml:"DB_NAME"`
}

var confg models.Config
var (
	GormDBConn *gorm.DB
)

func Yamlconfig() {

	configFile, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalln("error reading yaml file", err)
	}
	err = yaml.Unmarshal(configFile, &confg)
	if err != nil {
		log.Fatalln("error writting yaml file to struct: ", err)
	}

	//	log.Println("Config is ready")
}

// InitLogger initialize a new logger with specific file to log and return it
func InitLogger() *log.Logger {
	logFileName := "logs/" + time.Now().Format("2006-01-02") + ".log"

	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return log.New(file, "", log.LstdFlags)
}

// func InitDatabaseMysql() *sql.DB {

// 	cfg := mysql.Config{
// 		User:                 confg.DB_USERNAME,
// 		Pasword:              confg.DB_PASSWORD,
// 		Net:                  "tcp",
// 		Addr:                 confg.DB_IP + ":" + confg.DB_PORT,
// 		DBName:               confg.DB_NAME,
// 		AllowNativePasswords: true,
// 		ParseTime:            true,
// 	}

// 	// open the database connssection with the config. if encounter any error print the error and exit from program
// 	database, error := sql.Open("mysql", cfg.FormatDSN())
// 	if error != nil {
// 		log.Fatalln("Error connecting to database", error)
// 	}

// 	// ping the database to make sure connection is successfull
// 	error = database.Ping()
// 	if error != nil {
// 		log.Fatalln("Error on ping the database", error)
// 	}

// 	return database
// }

// connectDb
func InitDatabaseMysqlGorm() {

	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:@tcp(127.0.0.1:3306)/go_email?charset=utf8mb4&parseTime=True&loc=Local"
	/*
		NOTE:
		To handle time.Time correctly, you need to include parseTime as a parameter. (more parameters)
		To fully support UTF-8 encoding, you need to change charset=utf8 to charset=utf8mb4. See this article for a detailed explanation
	*/

	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}

	log.Println("connected")
	db.AutoMigrate(&models.User{}, &models.Queue{})
	GormDBConn = db

}
