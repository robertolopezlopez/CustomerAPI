package cron

import (
	"api/dao"
	"api/logging"
	"time"

	"github.com/go-co-op/gocron"
)

// Scheduler configures and starts the scheduler asynchronously
func Scheduler() (*gocron.Scheduler, error) {
	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(1).Second().Tag("clean up old entries").Do(deleteOldEntries)
	if err != nil {
		return nil, err
	}
	s.StartAsync()
	return s, nil
}

func deleteOldEntries() {
	rows, err := dao.DAO.DeleteOld(300)
	if err != nil {
		logging.ErrorLogger.Printf("CRON: %s", err.Error())
	}
	if rows != 0 {
		logging.InfoLogger.Printf("CRON: deleted %d old entries", rows)
	}
}
