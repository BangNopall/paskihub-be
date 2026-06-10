package repository

import (
	"context"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type privateFileRepository struct {
	db *gorm.DB
}

func NewPrivateFileRepository(db *gorm.DB) contracts.PrivateFileRepository {
	return &privateFileRepository{db: db}
}

func (r *privateFileRepository) FindPrivateFile(
	ctx context.Context,
	resourceType string,
	resourceID uuid.UUID,
) (contracts.PrivateFileRecord, error) {
	var record contracts.PrivateFileRecord
	query := r.db.WithContext(ctx)

	switch resourceType {
	case "team-recommendation":
		query = query.
			Table("teams AS t").
			Select("t.rec_letter_path AS path, i.user_id AS participant_user_id").
			Joins("JOIN institutions AS i ON i.id = t.insti_id").
			Where("t.id = ?", resourceID)
	case "member-id-card":
		query = query.
			Table("team_members AS tm").
			Select("tm.id_card_path AS path, i.user_id AS participant_user_id").
			Joins("JOIN teams AS t ON t.id = tm.team_id").
			Joins("JOIN institutions AS i ON i.id = t.insti_id").
			Where("tm.id = ?", resourceID)
	case "member-photo":
		query = query.
			Table("team_members AS tm").
			Select("tm.photo_path AS path, i.user_id AS participant_user_id").
			Joins("JOIN teams AS t ON t.id = tm.team_id").
			Joins("JOIN institutions AS i ON i.id = t.insti_id").
			Where("tm.id = ?", resourceID)
	case "registration-payment":
		query = query.
			Table("registrations AS r").
			Select("r.payment_proof_path AS path, i.user_id AS participant_user_id, e.user_id AS organizer_user_id").
			Joins("JOIN teams AS t ON t.id = r.team_id").
			Joins("JOIN institutions AS i ON i.id = t.insti_id").
			Joins("JOIN event_levels AS el ON el.id = r.event_level_id").
			Joins("JOIN events AS e ON e.id = el.event_id").
			Where("r.id = ?", resourceID)
	case "wallet-topup-proof":
		query = query.
			Table("wallet_transactions AS wt").
			Select("wt.proof_path AS path, e.user_id AS organizer_user_id").
			Joins("JOIN wallets AS w ON w.id = wt.wallet_id").
			Joins("JOIN events AS e ON e.id = w.event_id").
			Where("wt.id = ?", resourceID)
	default:
		return record, domain.ErrBadRequest
	}

	if err := query.Take(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return record, domain.ErrNotFound
		}
		return record, domain.ErrInternalServer
	}
	return record, nil
}
