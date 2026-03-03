package enums

type EventStatus string
type Role string
type InstitutionType string
type TeamType string
type RegistrationStatus string
type WalletType string
type TransactionStatus string

const (
	Organizer Role = "ORGANIZER"
	Peserta   Role = "PESERTA"
	Admin     Role = "ADMIN"

	Draft    EventStatus = "DRAFT"
	Open     EventStatus = "OPEN"
	Closed   EventStatus = "CLOSED"
	Archived EventStatus = "ARCHIVED"

	SD    InstitutionType = "SD"
	SMP   InstitutionType = "SMP"
	SMA   InstitutionType = "SMA"
	PURNA InstitutionType = "PURNA"
	UMUM  InstitutionType = "UMUM"

	Pasukan  TeamType = "PASUKAN"
	Danpas   TeamType = "DANPAS"
	Official TeamType = "OFFICIAL"
	Pelatih  TeamType = "PELATIH"

	Waiting  RegistrationStatus = "WAITING"
	DPPaid   RegistrationStatus = "DP_PAID"
	FullPaid RegistrationStatus = "FULL_PAID"
	Rejected RegistrationStatus = "REJECTED"

	TopUp    WalletType = "TOPUP"
	Withdraw WalletType = "WITHDRAW"

	Pending    TransactionStatus = "PENDING"
	Approve    TransactionStatus = "APPROVE"
	TSRejected TransactionStatus = "REJECTED"
)

func IsValidEventStatus(s string) bool {
	switch EventStatus(s) {
	case Draft, Open, Closed, Archived:
		return true
	default:
		return false
	}
}

func IsValidRole(r string) bool {
	switch Role(r) {
	case Organizer, Peserta, Admin:
		return true
	default:
		return false
	}
}

func IsValidInstitutionType(s string) bool {
	switch InstitutionType(s) {
	case SD, SMP, SMA, PURNA, UMUM:
		return true
	default:
		return false
	}
}
