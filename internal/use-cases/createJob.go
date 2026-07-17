package usecases

import (
	"encoding/json"
	"os"
	e "task-queue/internal/entities"

	"github.com/google/uuid"
)

func CreateJob(job e.Job) error {
	var db e.DataBase

	fileBytes, err := os.ReadFile("db.json")
	if err != nil {
		return err
	}

	if err = json.Unmarshal(fileBytes, &db); err != nil {
		return err
	}

	job.Id = uuid.New().String()

	db.Jobs = append(db.Jobs, job)

	updatedBytes, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}

	if err = os.WriteFile("db.json", updatedBytes, 0644); err != nil{
		return err
	}
	
	return nil
}