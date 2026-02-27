package application

import "github.com/rubenbuelvas/ringover/src/domain"

type IOrganizerService interface {
	GetProjects() ([]domain.Project, error)
	GetProjectByID(id int64) (domain.Project, error)
	CreateProject(data CreateProjectDTO) error
	UpdateProject(id int64, data UpdateProjectDTO) error
	DeleteProject(id int64) error
}
