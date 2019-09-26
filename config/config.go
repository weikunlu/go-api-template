package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

type AppConfig struct {
	AppEnv        string
	AppLogPath    string
	ServerHost    string
	ServerPort    string
	DbHost        string
	DbPort        string
	DbName        string
	DbUserName    string
	DbPassword    string
	RedisHost     string
	RedisPassword string
	RedisPort     string
	RedisDb       int
}

type BuildConfig struct {
	Version  string
	CommitId string
	Date     string
}

const projectDirName = "kkday-appphone-api"

var appConfig *AppConfig

var buildConfig BuildConfig

func init() {
	// Find root path
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	// Load .env file
	isTesting, _ := os.LookupEnv("TESTING")
	if isTesting == "true" {
		fmt.Println("Loading .env for testing")

		if err := godotenv.Load(string(rootPath) + `/.env`); err != nil {
			panic("No .env found")
		}
	} else {
		// look up execute path
		executeFolder := ""
		if err := godotenv.Load(); err != nil {
			// fail to get env at execution path
			executePath, _ := os.Executable()
			executeFolder = filepath.Dir(executePath)
			maxTries := 2
			for i := 0; i <= maxTries; i++ {
				if i == maxTries {
					// fail to get env at execute path
					panic("No .env found")
				}
				if err := godotenv.Load(executeFolder + `/.env`); err == nil {
					break
				}
				executeFolder = filepath.Dir(executeFolder)
			}
		}

		fmt.Printf("%s/.env loaded\n", executeFolder)
	}

	appLogPath := os.Getenv("APP_LOG_PATH")
	if isTesting == "true" {
		appLogPath = string(rootPath) + `/logs/service_test.log`
	}

	redisDb, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	appConfig = &AppConfig{
		AppEnv:        os.Getenv("APP_ENV"),
		AppLogPath:    appLogPath,
		ServerHost:    os.Getenv("SERVER_HOST"),
		ServerPort:    os.Getenv("SERVER_PORT"),
		DbHost:        os.Getenv("DB_HOST"),
		DbPort:        os.Getenv("DB_PORT"),
		DbName:        os.Getenv("DB_DATABASE"),
		DbUserName:    os.Getenv("DB_USERNAME"),
		DbPassword:    os.Getenv("DB_PASSWORD"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisDb:       redisDb,
	}

	fmt.Println("App Env: ", appConfig.AppEnv)
}

func GetAppConfig() *AppConfig {
	return appConfig
}

func SetBuildConfig(version string, build string, buildDate string) {
	buildConfig = BuildConfig{
		Version:  version,
		CommitId: build,
		Date:     buildDate,
	}
}

func GetBuildConfig() BuildConfig {
	return buildConfig
}
