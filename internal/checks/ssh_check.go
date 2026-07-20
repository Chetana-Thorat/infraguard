package checks

import (
	"bytes"

	"github.com/Chetana-Thorat/infraguard/internal/models"
)

type SSHOpenToWorldCheck struct{}

func (c SSHOpenToWorldCheck) ID() string {
	return "SSH_OPEN_TO_WORLD"
}

func (c SSHOpenToWorldCheck) Description() string {
	return "Detects security groups that allow SSH from 0.0.0.0/0"
}

func (c SSHOpenToWorldCheck) Run(filePath string, content []byte) ([]models.Finding, error) {

	hasSSH := bytes.Contains(content, []byte("from_port")) &&
		bytes.Contains(content, []byte("22")) &&
		bytes.Contains(content, []byte(`cidr_blocks = ["0.0.0.0/0"]`))

	if !hasSSH {
		return nil, nil
	}

	return []models.Finding{
		{
			RuleID:         c.ID(),
			Severity:       models.SeverityHigh,
			FilePath:       filePath,
			Message:        "Security group allows SSH from 0.0.0.0/0",
			Recommendation: "Restrict SSH access to trusted CIDR ranges.",
		},
	}, nil
}
