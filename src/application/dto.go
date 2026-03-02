package application

import "github.com/rubenbuelvas/game-organizer-be/src/domain"

type OrganizerDTO struct {
	NumberOfTeams int             `json:"number_of_teams"`
	BreakFamilies bool            `json:"break_families"`
	Stats         domain.Stats    `json:"stats"` // value between 0 and 10 indicating how much to prioritize balancing the specified stats
	Players       []domain.Player `json:"players"`
}
