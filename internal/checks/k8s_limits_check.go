package checks

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/Chetana-Thorat/infraguard/internal/models"
	"gopkg.in/yaml.v3"
)

type K8sResourceLimitsCheck struct{}

func (c K8sResourceLimitsCheck) ID() string {
	return "K8S_RESOURCE_LIMITS"
}

func (c K8sResourceLimitsCheck) Description() string {
	return "Detects containers without CPU and memory limits"
}

func (c K8sResourceLimitsCheck) Run(filePath string, content []byte) ([]models.Finding, error) {
	dec := yaml.NewDecoder(bytes.NewReader(content))

	var findings []models.Finding

	for {
		var doc any
		if err := dec.Decode(&doc); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("decode yaml: %w", err)
		}

		if doc == nil {
			continue
		}

		walkYAML(doc, func(container map[string]any) {
			if hasResourceLimits(container) {
				return
			}

			name, _ := container["name"].(string)
			message := "Container is missing CPU and memory limits"
			if name != "" {
				message = fmt.Sprintf("Container %q is missing CPU and/or memory limits", name)
			}

			findings = append(findings, models.Finding{
				RuleID:         c.ID(),
				Severity:       models.SeverityMedium,
				FilePath:       filePath,
				Message:        message,
				Recommendation: "Add resources.limits.cpu and resources.limits.memory to the container.",
			})
		})
	}

	return findings, nil
}

func walkYAML(node any, visit func(map[string]any)) {
	switch v := node.(type) {
	case map[string]any:
		if containers, ok := v["containers"].([]any); ok {
			for _, item := range containers {
				if container, ok := item.(map[string]any); ok {
					visit(container)
				}
			}
		}

		for _, child := range v {
			walkYAML(child, visit)
		}

	case []any:
		for _, child := range v {
			walkYAML(child, visit)
		}
	}
}

func hasResourceLimits(container map[string]any) bool {
	resources, ok := container["resources"].(map[string]any)
	if !ok {
		return false
	}

	limits, ok := resources["limits"].(map[string]any)
	if !ok {
		return false
	}

	cpu := strings.TrimSpace(fmt.Sprint(limits["cpu"]))
	memory := strings.TrimSpace(fmt.Sprint(limits["memory"]))

	return cpu != "" && cpu != "<nil>" && memory != "" && memory != "<nil>"
}
