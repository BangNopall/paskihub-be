package dto

import (
	"time"

	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/google/uuid"
)

type OrganizerDashboardRes struct {
	Stats            OrganizerStats     `json:"stats"`
	RecentActivities []EOActivityRes    `json:"recent_activities"`
	UpcomingEvents   []UpcomingEventRes `json:"upcoming_events"`
}

type OrganizerStats struct {
	TotalEvent  StatValue `json:"total_event"`
	TotalTeam   StatValue `json:"total_team"`
	CoinBalance CoinValue `json:"coin_balance"`
	Revenue     StatValue `json:"revenue"`
}

type StatValue struct {
	Value float64 `json:"value"`
	Trend string  `json:"trend"`
}

type CoinValue struct {
	Value float64 `json:"value"`
	Coins float64 `json:"coins"`
}

type EOActivityRes struct {
	Id        uuid.UUID `json:"id"`
	TeamName  string    `json:"team_name"`
	EventName string    `json:"event_name"`
	TimeAgo   string    `json:"time_ago"`
	Status    string    `json:"status"`
}

type UpcomingEventRes struct {
	Id              uuid.UUID `json:"id"`
	Title           string    `json:"title"`
	Date            string    `json:"date"`
	RegisteredTeams int       `json:"registered_teams"`
	Status          string    `json:"status"`
}

type ParticipantDashboardRes struct {
	Stats            ParticipantStats              `json:"stats"`
	RecentActivities []ParticipantActivity         `json:"recent_activities"`
	UpcomingEvents   []ParticipantUpcomingEventRes `json:"upcoming_events"`
}

type ParticipantStats struct {
	TotalTeam      int `json:"total_team"`
	ActiveEvent    int `json:"active_event"`
	FinishedEvent  int `json:"finished_event"`
	PendingPayment int `json:"pending_payment"`
}

type ParticipantActivity struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Time        string `json:"time"`
}

type ParticipantUpcomingEventRes struct {
	Id              uuid.UUID `json:"id"`
	Title           string    `json:"title"`
	Date            string    `json:"date"`
	RegisteredTeams int       `json:"registered_teams"`
	Status          string    `json:"status"`
	DetailURLId     uuid.UUID `json:"detail_url_id"`
}

type AdminDashboardRes struct {
	Stats              AdminDashboardStats               `json:"stats"`
	RecentTransactions []AdminDashboardTransactionRes    `json:"recent_transactions"`
	EORegistrations    []AdminDashboardEORegistrationRes `json:"eo_registrations"`
}

type HomeStatsResponse struct {
	TotalEvents       int64 `json:"total_events"`
	TotalOrganizers   int64 `json:"total_organizers"`
	TotalParticipants int64 `json:"total_participants"`
	TotalTeams        int64 `json:"total_teams"`
}

type AdminDashboardStats struct {
	TotalRevenue      AdminDashboardStatValue `json:"total_revenue"`
	TotalEO           AdminDashboardStatValue `json:"total_eo"`
	TotalParticipants AdminDashboardStatValue `json:"total_participants"`
	PendingTopups     AdminDashboardStatValue `json:"pending_topups"`
}

type AdminDashboardStatValue struct {
	Value float64 `json:"value"`
	Trend string  `json:"trend"`
}

type AdminDashboardTransactionRes struct {
	Id         uuid.UUID               `json:"id"`
	EOName     string                  `json:"eo_name"`
	Amount     float64                 `json:"amount"`
	AmountKoin float64                 `json:"amount_koin"`
	TimeAgo    string                  `json:"time_ago"`
	Status     enums.TransactionStatus `json:"status"`
}

type AdminDashboardEORegistrationRes struct {
	Id           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	RegisteredAt string    `json:"registered_at"`
}

type AdminEORegistrationRaw struct {
	Id        uuid.UUID
	Name      string
	Email     string
	CreatedAt time.Time
}
