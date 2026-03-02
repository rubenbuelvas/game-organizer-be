package application

import "github.com/rubenbuelvas/game-organizer-be/src/domain"

type OrganizerDTO struct {
	NumberOfTeams int             `json:"number_of_teams"`
	BreakFamilies bool            `json:"break_families"`
	BalanceBy     string          `json:"balance_by"`
	Players       []domain.Player `json:"players"`
}
