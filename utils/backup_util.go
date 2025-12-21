package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func StartBackupScheduler() {
	c := cron.New()

	// Запуск в 01:00 ночи
	_, err := c.AddFunc("0 1 * * *", func() {
		ExecuteBackup()
	})

	if err != nil {
		log.Fatalf("Ошибка планировщика: %v", err)
	}

	c.Start()
	fmt.Println("Планировщик бекапов запущен (данные из .env)")
}

func ExecuteBackup() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Не удалось загрузить .env файл, используем системные переменные")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	backupDir := os.Getenv("BACKUP_PATH")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		os.Mkdir(backupDir, os.ModePerm)
	}

	fileName := fmt.Sprintf("%s/%s_%s.sql", backupDir, dbName, time.Now().Format("2006-01-02_15-04-05"))

	cmd := exec.Command(
		"pg_dump",
		"-h", dbHost,
		"-p", dbPort,
		"-U", dbUser,
		"-d", dbName,
		"-f", fileName,
	)

	cmd.Env = append(os.Environ(), "PGPASSWORD="+dbPass)
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("Ошибка бекапа: %v", err)
		return
	}

	log.Printf("Бекап сохранен: %s", fileName)
}
