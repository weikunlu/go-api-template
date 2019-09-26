package database

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jinzhu/gorm"
	"github.com/weikunlu/go-api-template/config"
	"log"
	"strconv"
)

var db *gorm.DB

func NewDatabase() (err error) {
	cfg := config.GetAppConfig()

	dbUri := fmt.Sprintf("host=%s port=%s dbname=%s sslmode=disable user=%s password=%s", cfg.DbHost, cfg.DbPort, cfg.DbName, cfg.DbUserName, cfg.DbPassword)

	fmt.Printf("init database connection\n")
	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Printf("fail to connect database %v\n", err.Error())
		return
	}

	err = conn.DB().Ping()
	if err != nil {
		fmt.Printf("fail to ping database %v\n", err.Error())
		return
	}
	fmt.Printf("database PING PONG\n")

	// Database migration
	if cfg.AppEnv == "prod" {
		conn.DB().SetMaxIdleConns(4)
	} else {
		conn.LogMode(true)

		// Migrate models using ORM
		//db.Debug().AutoMigrate(&models.Member{})
	}

	db = conn

	return
}

func PerformMigrations(direction string, step string) {
	cfg := config.GetAppConfig()
	driver, err := postgres.WithInstance(db.DB(), &postgres.Config{})
	if err != nil {
		log.Fatal(err)
		fmt.Printf("get database connection error %v\n", err.Error())
	}

	fsrc, err := (&file.File{}).Open("file://migrations")
	if err != nil {
		fmt.Printf("load migration files error %v\n", err.Error())
	}

	m, err := migrate.NewWithInstance("file", fsrc, cfg.DbName, driver)
	if err != nil {
		fmt.Printf("init migration error %v\n", err.Error())
	}
	defer m.Close()

	// overwrite down/up to step if step value exist
	var stepDirection string
	if len(step) > 0 && (direction == "up" || direction == "down") {
		stepDirection = direction
		direction = "step"
	}

	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			fmt.Printf("execute migrate up error %v\n", err.Error())
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			fmt.Printf("execute migrate up error %v\n", err.Error())
		}
	case "step":
		// Steps looks at the currently active migration version. It will migrate up if n > 0, and down if n < 0.
		iStep, _ := strconv.Atoi(step)
		if stepDirection == "down" {
			iStep = iStep * -1
		}
		if err := m.Steps(iStep); err != nil && err != migrate.ErrNoChange {
			fmt.Printf("execute migrate step %s %s version error %v\n", stepDirection, step, err.Error())
		}
	case "m_version":
		version, _ := strconv.ParseUint(step, 10, 64)
		if err := m.Migrate(uint(version)); err != nil && err != migrate.ErrNoChange {
			fmt.Printf("execute migrate specified version %s error %v\n", step, err.Error())
		}
	}

	ver, dirty, _ := m.Version()
	fmt.Printf("migration process completed ver:%v dirty:%v\n", ver, dirty)
}

func GetDb() *gorm.DB {
	return db
}
