package db

import "github.com/ipreferwater/pikmin/go/model"

var (
	PikminRepo    PikminRepository
)

type PikminRepository interface {
	CreatePikmin(model model.Pikmin) (string, error)
	UpdatePikmin(id string, newModel string) error
	GetPikminsByColor(color string) ([]model.Pikmin, error)
	GiveBombs(pikmins []model.Pikmin) (int64, error)
}
