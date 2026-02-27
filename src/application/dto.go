package application

import "github.com/rubenbuelvas/game-organizer-be/src/domain"

type OrganizerDTO struct {
	NumberOfTeams    int             `json:"number_of_teams"`
	BreakFamilies    bool            `json:"break_families"`
	StatBalanceType  string          `json:"stat_balance_type"`  // Inside PlayerStats
	SkillBalanceType string          `json:"skill_balance_type"` // Inside Skills
	Priority         string          `json:"priority"`           // "skills" or "stats"
	Players          []domain.Player `json:"players"`
}
