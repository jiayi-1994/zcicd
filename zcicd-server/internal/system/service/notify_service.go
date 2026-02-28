package service

import (
	"github.com/zcicd/zcicd-server/internal/system/model"
	"github.com/zcicd/zcicd-server/internal/system/repository"
)

type NotifyService struct {
	channelRepo *repository.NotifyRepository
	ruleRepo    *repository.RuleRepository
}

func NewNotifyService(cr *repository.NotifyRepository, rr *repository.RuleRepository) *NotifyService {
	return &NotifyService{channelRepo: cr, ruleRepo: rr}
}

func (s *NotifyService) CreateChannel(req CreateChannelReq) (*model.NotifyChannel, error) {
	c := &model.NotifyChannel{
		Name:        req.Name,
		ChannelType: req.ChannelType,
		Enabled:     true,
	}
	return c, s.channelRepo.CreateChannel(c)
}

func (s *NotifyService) GetChannel(id string) (*model.NotifyChannel, error) {
	return s.channelRepo.GetChannel(id)
}

func (s *NotifyService) UpdateChannel(id string, req UpdateChannelReq) (*model.NotifyChannel, error) {
	c, err := s.channelRepo.GetChannel(id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		c.Name = req.Name
	}
	if req.Enabled != nil {
		c.Enabled = *req.Enabled
	}
	return c, s.channelRepo.UpdateChannel(c)
}

func (s *NotifyService) DeleteChannel(id string) error {
	return s.channelRepo.DeleteChannel(id)
}

func (s *NotifyService) ListChannels() ([]model.NotifyChannel, error) {
	return s.channelRepo.ListChannels()
}

func (s *NotifyService) CreateRule(req CreateRuleReq) (*model.NotifyRule, error) {
	r := &model.NotifyRule{
		Name:      req.Name,
		EventType: req.EventType,
		ChannelID: req.ChannelID,
		ProjectID: req.ProjectID,
		Severity:  req.Severity,
		Enabled:   true,
	}
	if r.Severity == "" {
		r.Severity = "all"
	}
	return r, s.ruleRepo.Create(r)
}

func (s *NotifyService) UpdateRule(id string, req UpdateRuleReq) (*model.NotifyRule, error) {
	r, err := s.ruleRepo.Get(id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		r.Name = req.Name
	}
	if req.EventType != "" {
		r.EventType = req.EventType
	}
	if req.Severity != "" {
		r.Severity = req.Severity
	}
	if req.Enabled != nil {
		r.Enabled = *req.Enabled
	}
	return r, s.ruleRepo.Update(r)
}

func (s *NotifyService) ListRules() ([]model.NotifyRule, error) {
	return s.ruleRepo.List()
}
