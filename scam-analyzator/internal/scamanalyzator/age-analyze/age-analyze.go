package ageanalyze

import (
	"time"
)

type AnalysisResult struct {
	IsNewToken bool      `json:"is_new_token"` // вердикт
	AgeHours   float64   `json:"age_hours"`    // возраст в часах для наглядности
	CreatedAt  time.Time `json:"created_at"`   // дата создания в читаемом виде
}

func AnalyzeAge(timestamp int64) AnalysisResult {
	// конвертируем Unix (секунды) в объект времени Go
	createdTime := time.Unix(timestamp, 0)

	// считаем разницу между сейчас и временем создания
	duration := time.Since(createdTime)

	// переводим разницу в часы
	hours := duration.Hours()

	// наш критерий скама - токену меньше 24 часов
	isScam := hours < 24

	return AnalysisResult{
		IsNewToken: isScam,
		AgeHours:   hours,
		CreatedAt:  createdTime,
	}
}
