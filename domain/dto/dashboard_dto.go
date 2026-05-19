package dto

import "github.com/google/uuid"

type OrganizerDashboardRes struct {
	Stats            OrganizerStats        `json:"stats"`
	RecentActivities []EOActivityRes      `json:"recent_activities"`
	UpcomingEvents   []UpcomingEventRes    `json:"upcoming_events"`
}

type OrganizerStats struct {
	TotalEvent   StatValue `json:"total_event"`
	TotalTeam    StatValue `json:"total_team"`
	CoinBalance  CoinValue `json:"coin_balance"`
	Revenue      StatValue `json:"revenue"`
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
	Stats            ParticipantStats       `json:"stats"`
	RecentActivities []ParticipantActivity `json:"recent_activities"`
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
