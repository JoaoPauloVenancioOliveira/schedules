package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
)

type Schedule struct {
	ID           uint      `gorm:"primaryKey" json:"id"` // ID auto incrementado
	ScheduleTime time.Time `json:"schedule_time"`
}

var db *gorm.DB

func init() {
	var err error
	dsn := os.Getenv("DATABASE_URL")
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.Migrator().DropTable(&Schedule{})
	db.AutoMigrate(&Schedule{}) // Cria a tabela se não existir
}

// Gera os horários disponíveis para as massagens
func generateSchedules() []Schedule {
	var newSchedules []Schedule

	location, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		panic(err)
	}

	// Intervalo da manhã: 09:00 - 12:00
	morningStart := time.Date(2024, 9, 22, 9, 0, 0, 0, location)
	morningEnd := time.Date(2024, 9, 22, 12, 0, 0, 0, location)

	// Intervalo da tarde: 14:00 - 17:00
	afternoonStart := time.Date(2024, 9, 22, 14, 0, 0, 0, location)
	afternoonEnd := time.Date(2024, 9, 22, 17, 0, 0, 0, location)

	// Gera os horários da manhã
	for t := morningStart; t.Before(morningEnd); t = t.Add(15 * time.Minute) {
		newSchedules = append(newSchedules, Schedule{
			ScheduleTime: t,
		})
	}

	// Gera os horários da tarde
	for t := afternoonStart; t.Before(afternoonEnd); t = t.Add(15 * time.Minute) {
		newSchedules = append(newSchedules, Schedule{
			ScheduleTime: t,
		})
	}

	// Salva os horários no banco de dados
	db.Create(&newSchedules)
	return newSchedules
}

func getSchedules(w http.ResponseWriter, r *http.Request) {
	var schedules []Schedule
	db.Find(&schedules) // Busca todos os horários no banco de dados
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedules)
}

func main() {
	http.HandleFunc("/schedules", getSchedules)
	generateSchedules() // Gera e armazena os horários na inicialização
	fmt.Println("Servidor rodando na porta :8080")
	http.ListenAndServe(":8080", nil)
}
