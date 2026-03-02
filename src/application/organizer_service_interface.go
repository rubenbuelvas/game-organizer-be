package application

import "github.com/rubenbuelvas/game-organizer-be/src/domain"

type IOrganizerService interface {
	Organize(dto OrganizerDTO) ([]domain.Team, error)
}
