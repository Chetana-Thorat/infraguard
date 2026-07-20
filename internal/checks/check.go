package checks

import "github.com/Chetana-Thorat/infraguard/internal/models"

type Check interface {
	ID() string
	Description() string
	Run(filePath string, content []byte) ([]models.Finding, error)
}
