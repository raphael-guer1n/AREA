package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/raphael-guer1n/AREA/CronService/internal/domain"
	"github.com/raphael-guer1n/AREA/CronService/internal/repository"
	"github.com/robfig/cron/v3"
)

type CronService struct {
	repo           *repository.ActionRepository
	cron           *cron.Cron
	jobs           map[int]cron.EntryID
	jobsMutex      sync.RWMutex
	areaServiceURL string
	internalSecret string
}

func NewCronService(repo *repository.ActionRepository, areaServiceURL, internalSecret string) *CronService {
	return &CronService{
		repo:           repo,
		cron:           cron.New(),
		jobs:           make(map[int]cron.EntryID),
		areaServiceURL: areaServiceURL,
		internalSecret: internalSecret,
	}
}

func (s *CronService) Start() {
	s.cron.Start()
	log.Println("Cron service started")

	actions, err := s.repo.GetAll()
	if err != nil {
		log.Printf("Failed to load existing actions: %v", err)
		return
	}

	for _, action := range actions {
		if action.Active {
			if err := s.scheduleAction(action); err != nil {
				log.Printf("Failed to schedule action %d: %v", action.ActionID, err)
			}
		}
	}
}

func (s *CronService) Stop() {
	s.cron.Stop()
	log.Println("Cron service stopped")
}

func (s *CronService) CreateAction(action *domain.Action) error {
	if err := s.repo.Create(action); err != nil {
		return fmt.Errorf("failed to create action: %w", err)
	}

	if action.Active {
		if err := s.scheduleAction(action); err != nil {
			return fmt.Errorf("failed to schedule action: %w", err)
		}
	}

	return nil
}

func (s *CronService) ActivateAction(actionID int) error {
	action, err := s.repo.GetByActionID(actionID)
	if err != nil {
		return fmt.Errorf("failed to get action: %w", err)
	}

	if action.Active {
		return nil
	}

	action.Active = true
	if err := s.repo.Update(action); err != nil {
		return fmt.Errorf("failed to update action: %w", err)
	}

	if err := s.scheduleAction(action); err != nil {
		return fmt.Errorf("failed to schedule action: %w", err)
	}

	return nil
}

func (s *CronService) DeactivateAction(actionID int) error {
	action, err := s.repo.GetByActionID(actionID)
	if err != nil {
		return fmt.Errorf("failed to get action: %w", err)
	}

	if !action.Active {
		return nil
	}

	action.Active = false
	if err := s.repo.Update(action); err != nil {
		return fmt.Errorf("failed to update action: %w", err)
	}

	s.unscheduleAction(actionID)

	return nil
}

func (s *CronService) DeleteAction(actionID int) error {
	s.unscheduleAction(actionID)

	if err := s.repo.Delete(actionID); err != nil {
		return fmt.Errorf("failed to delete action: %w", err)
	}

	return nil
}

func (s *CronService) scheduleAction(action *domain.Action) error {
	var delaySeconds int
	for _, input := range action.Input {
		if input.Name == "delay" {
			val, err := strconv.Atoi(input.Value)
			if err != nil {
				return fmt.Errorf("invalid delay value: %w", err)
			}
			delaySeconds = val
			break
		}
	}

	if delaySeconds <= 0 {
		return fmt.Errorf("delay must be greater than 0")
	}

	cronExpr := fmt.Sprintf("@every %ds", delaySeconds)

	entryID, err := s.cron.AddFunc(cronExpr, func() {
		s.triggerAction(action)
	})

	if err != nil {
		return fmt.Errorf("failed to add cron job: %w", err)
	}

	s.jobsMutex.Lock()
	s.jobs[action.ActionID] = entryID
	s.jobsMutex.Unlock()

	log.Printf("Scheduled action %d with delay %d seconds", action.ActionID, delaySeconds)

	return nil
}

func (s *CronService) unscheduleAction(actionID int) {
	s.jobsMutex.Lock()
	defer s.jobsMutex.Unlock()

	if entryID, exists := s.jobs[actionID]; exists {
		s.cron.Remove(entryID)
		delete(s.jobs, actionID)
		log.Printf("Unscheduled action %d", actionID)
	}
}

func (s *CronService) triggerAction(action *domain.Action) {
	log.Printf("Triggering action %d", action.ActionID)

	outputFields := []domain.OutputField{}
	for _, input := range action.Input {
		if input.Name == "delay" {
			outputFields = append(outputFields, domain.OutputField{
				Name:  "delay",
				Value: input.Value,
			})
		}
	}

	triggerReq := domain.TriggerAreaRequest{
		ActionID:     action.ActionID,
		OutputFields: outputFields,
	}

	if err := s.callAreaService(triggerReq); err != nil {
		log.Printf("Failed to trigger area service for action %d: %v", action.ActionID, err)
	}
}

func (s *CronService) callAreaService(req domain.TriggerAreaRequest) error {
	url := fmt.Sprintf("%s/triggerArea", s.areaServiceURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Internal-Secret", s.internalSecret)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("area service returned status %d", resp.StatusCode)
	}

	log.Printf("Successfully triggered area service for action %d", req.ActionID)
	return nil
}
