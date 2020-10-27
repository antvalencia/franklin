package models

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

func SetupModels() *gorm.DB {
	pwd, _ := os.Getwd()
	feed_csv := filepath.Join(
		pwd,
		"models",
		"data",
		"feed.csv",
	)

	conn_str := getConnStr()
	db, err := gorm.Open("postgres", conn_str)
	if err != nil {
		panic("Failed to connect to database!")
	}

	db.AutoMigrate(&Feed{})
	db.AutoMigrate(&NewsItem{})

	lines, err := readCsv(feed_csv)
	if err != nil {
		panic(err)
	}

	for _, line := range lines {
		f := Feed{
			Category: line[0],
			Source:   line[1],
			Title:    line[2],
			URL:      line[3],
		}
		db.Create(&f)
	}

	return db
}

func getConnStr() string {
	viper.AutomaticEnv()

	v_user := viper.Get("POSTGRES_USER")
	v_password := viper.Get("POSTGRES_PASSWORD")
	v_db := viper.Get("POSTGRES_DB")
	v_host := viper.Get("POSTGRES_HOST")
	v_port := viper.Get("POSTGRES_PORT")

	return fmt.Sprintf(
		"host=%v port=%v user=%v dbname=%v password=%v sslmode=disable",
		v_host,
		v_port,
		v_user,
		v_db,
		v_password,
	)
}

func readCsv(filename string) ([][]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}
