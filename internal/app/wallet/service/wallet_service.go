package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/domain/enums"
	uuidPkg "github.com/BangNopall/paskihub-be/pkg/uuid"
	"github.com/google/uuid"
)

type walletService struct {
	walletRepo contracts.WalletRepository
	eventRepo  contracts.EventRepository
	uuid       uuidPkg.UUIDInterface
	timeout    time.Duration
}

func NewWalletService(
	walletRepo contracts.WalletRepository,
	eventRepo contracts.EventRepository,
	uuid uuidPkg.UUIDInterface,
	timeout time.Duration,
) contracts.WalletService {
	return &walletService{
		walletRepo: walletRepo,
		eventRepo:  eventRepo,
		uuid:       uuid,
		timeout:    timeout,
	}
}

func (s *walletService) checkEventOwnership(ctx context.Context, eventId uuid.UUID, userId uuid.UUID) error {
	event, err := s.eventRepo.FetchOneById(ctx, eventId)
	if err != nil {
		return err
	}
	if event.UserId != userId {
		return domain.ErrForbidden
	}
	return nil
}

func (s *walletService) GetWalletInfo(ctx context.Context, eventId string, userId string) (*dto.WalletResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	eId, err := uuid.Parse(eventId)
	if err != nil {
		return nil, domain.ErrInternalServer
	}

	uId, err := uuid.Parse(userId)
	if err != nil {
		return nil, domain.ErrInternalServer
	}

	if err := s.checkEventOwnership(ctx, eId, uId); err != nil {
		return nil, err
	}

	wallet, err := s.walletRepo.GetWalletByEventId(ctx, eId)
	if err != nil {
		return nil, err
	}

	return dto.WalletEntityToResponse(wallet), nil
}

func (s *walletService) GetTransactionLogs(ctx context.Context, eventId string, userId string) ([]dto.WalletTransactionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	eId, err := uuid.Parse(eventId)
	if err != nil {
		return nil, domain.ErrInternalServer
	}

	uId, err := uuid.Parse(userId)
	if err != nil {
		return nil, domain.ErrInternalServer
	}

	if err := s.checkEventOwnership(ctx, eId, uId); err != nil {
		return nil, err
	}

	wallet, err := s.walletRepo.GetWalletByEventId(ctx, eId)
	if err != nil {
		return nil, err
	}

	transactions, err := s.walletRepo.GetTransactionLogs(ctx, wallet.Id)
	if err != nil {
		return nil, err
	}

	var responses []dto.WalletTransactionResponse
	for _, tx := range transactions {
		responses = append(responses, *dto.WalletTransactionEntityToResponse(&tx))
	}

	return responses, nil
}

func (s *walletService) RequestTopUp(ctx context.Context, eventId string, userId string, req *dto.TopUpRequest, proofFile *multipart.FileHeader) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	eId, err := uuid.Parse(eventId)
	if err != nil {
		return domain.ErrInternalServer
	}

	uId, err := uuid.Parse(userId)
	if err != nil {
		return domain.ErrInternalServer
	}

	if err := s.checkEventOwnership(ctx, eId, uId); err != nil {
		return err
	}

	wallet, err := s.walletRepo.GetWalletByEventId(ctx, eId)
	if err != nil {
		return err
	}

	if req.Amount < 50000 {
		return domain.ErrBadRequest
	}

	// Validate file extension, assuming common image formats
	ext := filepath.Ext(proofFile.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return domain.ErrInvalidProofType
	}

	txId, err := s.uuid.New()
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s-topup-%d%s", txId.String(), time.Now().Unix(), ext)
	path := filepath.Join("public", "uploads", "wallets", filename)

	if err := s.saveFile(proofFile, path); err != nil {
		return domain.ErrInternalServer
	}

	transaction := &entity.WalletTransaction{
		Id:        txId,
		WalletId:  wallet.Id,
		Type:      enums.TopUp,
		Amount:    req.Amount,
		ProofPath: path,
		Status:    enums.Pending,
	}

	return s.walletRepo.CreateTransaction(ctx, transaction)
}

func (s *walletService) ApproveTopUp(ctx context.Context, transactionId string, adminUserId string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	txId, err := uuid.Parse(transactionId)
	if err != nil {
		return domain.ErrInternalServer
	}

	// Assume router handles admin authorization via JWT token middlewares

	return s.walletRepo.ApproveTransaction(ctx, txId)
}

func (s *walletService) RejectTopUp(ctx context.Context, transactionId string, adminUserId string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	txId, err := uuid.Parse(transactionId)
	if err != nil {
		return domain.ErrInternalServer
	}

	// Assume router handles admin authorization via JWT token middlewares

	return s.walletRepo.RejectTransaction(ctx, txId)
}

func (s *walletService) saveFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
