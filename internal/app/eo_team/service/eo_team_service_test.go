package service

import (
	"context"
	"testing"

	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/google/uuid"
)

type stubEOTeamRepository struct {
	approveCalled  bool
	eventID        uuid.UUID
	registrationID uuid.UUID
	fee            float64
	status         enums.RegistrationStatus
}

func (r *stubEOTeamRepository) CheckEventOwnership(ctx context.Context, eventID, userID uuid.UUID) (bool, error) {
	return true, nil
}
func (r *stubEOTeamRepository) FindAllRegistrationsByEvent(ctx context.Context, eventID uuid.UUID, eventLevelID *uuid.UUID, institutionType *string) ([]entity.Registration, error) {
	return nil, nil
}
func (r *stubEOTeamRepository) FindRegistrationByIdAndEvent(ctx context.Context, registrationID, eventID uuid.UUID) (*entity.Registration, error) {
	return &entity.Registration{Id: registrationID}, nil
}
func (r *stubEOTeamRepository) UpdateRegistration(ctx context.Context, registration *entity.Registration) error {
	return nil
}
func (r *stubEOTeamRepository) ApproveRegistration(
	ctx context.Context,
	eventID uuid.UUID,
	registrationID uuid.UUID,
	totalFee float64,
	status enums.RegistrationStatus,
) error {
	r.approveCalled = true
	r.eventID = eventID
	r.registrationID = registrationID
	r.fee = totalFee
	r.status = status
	return nil
}
func (r *stubEOTeamRepository) GetStats(ctx context.Context, eventID uuid.UUID) (*dto.EOTeamStatsRes, error) {
	return nil, nil
}
func (r *stubEOTeamRepository) GetAssessmentStatus(ctx context.Context, registrationID uuid.UUID) (string, error) {
	return "", nil
}

type stubSettingRepository struct {
	setting entity.SystemSetting
}

func (r stubSettingRepository) Get(ctx context.Context) (*entity.SystemSetting, error) {
	return &r.setting, nil
}
func (r stubSettingRepository) Update(ctx context.Context, setting *entity.SystemSetting) error {
	return nil
}

func TestApproveTeamDelegatesSingleAtomicRepositoryOperation(t *testing.T) {
	eventID := uuid.New()
	registrationID := uuid.New()
	repo := &stubEOTeamRepository{}
	svc := &eoTeamService{
		repo: repo,
		settingRepo: stubSettingRepository{setting: entity.SystemSetting{
			ApprovalFee: 2,
			CoinRate:    1000,
		}},
	}

	err := svc.ApproveTeam(
		context.Background(),
		eventID,
		uuid.New(),
		registrationID,
		dto.EOTeamApproveReq{PaymentStatus: enums.FullPaid},
	)
	if err != nil {
		t.Fatalf("ApproveTeam returned error: %v", err)
	}
	if !repo.approveCalled {
		t.Fatal("atomic approval repository method was not called")
	}
	if repo.eventID != eventID || repo.registrationID != registrationID {
		t.Fatal("approval repository received incorrect IDs")
	}
	if repo.fee != 2000 || repo.status != enums.FullPaid {
		t.Fatalf("unexpected approval values: fee=%v status=%s", repo.fee, repo.status)
	}
}
