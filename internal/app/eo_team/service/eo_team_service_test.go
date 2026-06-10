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
	approvalFee    float64
	coinRate       float64
	status         enums.RegistrationStatus
}

func (r *stubEOTeamRepository) CheckEventOwnership(ctx context.Context, eventId, eoId uuid.UUID) (bool, error) {
	return true, nil
}
func (r *stubEOTeamRepository) FindAllRegistrationsByEvent(ctx context.Context, eventId uuid.UUID, eventLevelId *uuid.UUID, institutionType *string) ([]entity.Registration, error) {
	return nil, nil
}
func (r *stubEOTeamRepository) FindRegistrationByIdAndEvent(ctx context.Context, registrationId, eventId uuid.UUID) (*entity.Registration, error) {
	return nil, nil
}
func (r *stubEOTeamRepository) UpdateRegistration(ctx context.Context, registration *entity.Registration) error {
	return nil
}
func (r *stubEOTeamRepository) ApproveRegistration(
	ctx context.Context,
	eventId uuid.UUID,
	registrationId uuid.UUID,
	totalFee float64,
	approvalFee float64,
	coinRate float64,
	status enums.RegistrationStatus,
) error {
	r.approveCalled = true
	r.eventID = eventId
	r.registrationID = registrationId
	r.fee = totalFee
	r.approvalFee = approvalFee
	r.coinRate = coinRate
	r.status = status
	return nil
}
func (r *stubEOTeamRepository) GetStats(ctx context.Context, eventId uuid.UUID) (*dto.EOTeamStatsRes, error) {
	return nil, nil
}
func (r *stubEOTeamRepository) GetAssessmentStatus(ctx context.Context, registrationId uuid.UUID) (string, error) {
	return "", nil
}

type stubSettingRepository struct{}

func (r *stubSettingRepository) Get(ctx context.Context) (*entity.SystemSetting, error) {
	return &entity.SystemSetting{
		ApprovalFee: 10,
		CoinRate:    1000,
	}, nil
}
func (r *stubSettingRepository) Update(ctx context.Context, setting *entity.SystemSetting) error {
	return nil
}

func TestApproveTeamDelegatesSingleAtomicRepositoryOperation(t *testing.T) {
	repo := &stubEOTeamRepository{}
	settingRepo := &stubSettingRepository{}
	svc := NewEOTeamService(repo, settingRepo)

	eventID := uuid.New()
	userID := uuid.New()
	registrationID := uuid.New()

	err := svc.ApproveTeam(
		context.Background(),
		eventID,
		userID,
		registrationID,
		dto.EOTeamApproveReq{PaymentStatus: enums.FullPaid},
	)

	if err != nil {
		t.Fatalf("ApproveTeam returned error: %v", err)
	}
	if !repo.approveCalled {
		t.Fatalf("Expected ApproveRegistration to be called")
	}
}
