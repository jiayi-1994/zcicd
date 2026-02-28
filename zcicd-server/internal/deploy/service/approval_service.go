package service

import (
	"fmt"
	"time"

	"github.com/zcicd/zcicd-server/internal/deploy/model"
	"github.com/zcicd/zcicd-server/internal/deploy/repository"
)

type ApprovalService struct {
	approvalRepo *repository.ApprovalRepository
	deployRepo   *repository.DeployRepository
}

func NewApprovalService(approvalRepo *repository.ApprovalRepository, deployRepo *repository.DeployRepository) *ApprovalService {
	return &ApprovalService{approvalRepo: approvalRepo, deployRepo: deployRepo}
}

func (s *ApprovalService) CreateApproval(historyID, envID, requestedBy string) (*model.ApprovalRecord, error) {
	record := &model.ApprovalRecord{
		DeployHistoryID: historyID,
		EnvironmentID:   envID,
		RequestedBy:     requestedBy,
		Status:          "pending",
	}
	if err := s.approvalRepo.Create(record); err != nil {
		return nil, err
	}
	return record, nil
}

func (s *ApprovalService) Approve(id, approverID string, req ApproveReq) (*model.ApprovalRecord, error) {
	record, err := s.approvalRepo.Get(id)
	if err != nil {
		return nil, err
	}
	if record.Status != "pending" {
		return nil, fmt.Errorf("approval already decided")
	}

	now := time.Now()
	record.ApproverID = &approverID
	record.Status = "approved"
	record.Comment = req.Comment
	record.DecidedAt = &now

	if err := s.approvalRepo.Update(record); err != nil {
		return nil, err
	}

	// Update deploy history status to allow sync
	history, err := s.deployRepo.GetHistory(record.DeployHistoryID)
	if err == nil {
		history.Status = "pending"
		s.deployRepo.UpdateHistory(history)
	}

	return record, nil
}

func (s *ApprovalService) Reject(id, approverID string, req RejectReq) (*model.ApprovalRecord, error) {
	record, err := s.approvalRepo.Get(id)
	if err != nil {
		return nil, err
	}
	if record.Status != "pending" {
		return nil, fmt.Errorf("approval already decided")
	}

	now := time.Now()
	record.ApproverID = &approverID
	record.Status = "rejected"
	record.Comment = req.Comment
	record.DecidedAt = &now

	if err := s.approvalRepo.Update(record); err != nil {
		return nil, err
	}

	// Update deploy history status to cancelled
	history, err := s.deployRepo.GetHistory(record.DeployHistoryID)
	if err == nil {
		history.Status = "cancelled"
		finished := time.Now()
		history.FinishedAt = &finished
		s.deployRepo.UpdateHistory(history)
	}

	return record, nil
}

func (s *ApprovalService) Get(id string) (*model.ApprovalRecord, error) {
	return s.approvalRepo.Get(id)
}

func (s *ApprovalService) ListPending(approverID string, page, pageSize int) ([]model.ApprovalRecord, int64, error) {
	return s.approvalRepo.ListPending(approverID, page, pageSize)
}
