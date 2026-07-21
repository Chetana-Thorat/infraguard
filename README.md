# InfraGuard

InfraGuard is a lightweight policy enforcement tool written in Go that scans Terraform and Kubernetes configuration files for security and operational misconfigurations before they reach production.

The tool is designed to be integrated into local development workflows and CI/CD pipelines, allowing teams to identify infrastructure policy violations early during code review.

---

## Why InfraGuard?

Infrastructure changes are often reviewed manually, making it easy for insecure or non-compliant configurations to be merged accidentally.

Examples include:

- Security groups exposing SSH to the public
- Kubernetes workloads deployed without CPU or memory limits
- Public cloud resources
- Misconfigured networking policies

InfraGuard automates these checks by scanning Infrastructure-as-Code (IaC) files and reporting violations before deployment.

---

## Current Features

- Scan Terraform (`.tf`) files
- Scan Kubernetes YAML manifests (`.yaml`, `.yml`)
- Detect Terraform security groups exposing SSH (`0.0.0.0/0` on port `22`)
- Detect Kubernetes containers missing CPU and memory resource limits
- Recursive directory scanning
- Extensible rule-based architecture
- CI-friendly exit codes
  - Exit `0` when no violations are found
  - Exit `1` when violations exist

---

## Project Structure

```text
infraguard/
├── cmd/
│   ├── root.go
│   └── scan.go
│
├── internal/
│   ├── checks/
│   │   ├── check.go
│   │   ├── ssh_check.go
│   │   └── k8s_limits_check.go
│   │
│   ├── config/
│   ├── models/
│   │   └── finding.go
│   ├── reporter/
│   └── scanner/
│       └── scanner.go
│
├── examples/
│   ├── sample.tf
│   └── sample.yaml
│
├── configs/
├── README.md
└── go.mod
```

---

## Architecture

```text
             +----------------+
             |  CLI (Cobra)   |
             +-------+--------+
                     |
                     v
             +----------------+
             | Directory Scan |
             +-------+--------+
                     |
          discovers Terraform/YAML
                     |
                     v
             +----------------+
             | Policy Engine  |
             +-------+--------+
                     |
         executes registered checks
                     |
                     v
             +----------------+
             | Findings       |
             +-------+--------+
                     |
                     v
             Terminal / CI Pipeline
```

---

## Supported Rules

### Terraform

| Rule | Severity |
|------|----------|
| SSH exposed to `0.0.0.0/0` | HIGH |

### Kubernetes

| Rule | Severity |
|------|----------|
| Missing CPU/Memory limits | MEDIUM |

---

## Example

### Terraform

```hcl
resource "aws_security_group" "web" {
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
```

Output

```text
❌ HIGH
Rule: SSH_OPEN_TO_WORLD
File: examples/sample.tf

Security group allows SSH from 0.0.0.0/0

Recommendation:
Restrict SSH access to trusted CIDR ranges.
```

---

### Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment

spec:
  template:
    spec:
      containers:
      - name: app
        image: nginx
```

Output

```text
❌ MEDIUM
Rule: K8S_RESOURCE_LIMITS
File: examples/sample.yaml

Container "app" is missing CPU and/or memory limits

Recommendation:
Add resources.limits.cpu and resources.limits.memory.
```

---

## Installation

Clone the repository

```bash
git clone https://github.com/Chetana-Thorat/infraguard.git
cd infraguard
```

Install dependencies

```bash
go mod tidy
```

Run

```bash
go run . --help
```

---

## Usage

Scan the current directory

```bash
go run . scan
```

Scan a specific directory

```bash
go run . scan examples
```

---

## Exit Codes

| Exit Code | Meaning |
|-----------|---------|
| 0 | No policy violations found |
| 1 | Policy violations detected |

These exit codes allow InfraGuard to fail CI/CD pipelines when infrastructure policies are violated.

---

## Design Principles

InfraGuard is built around a modular architecture.

Each policy check implements a common interface, making it easy to introduce new rules without modifying the scanning engine.

This separation allows:

- Independent rule development
- Reusable policy checks
- Improved maintainability
- Easier testing

---

## Roadmap

### Completed

- Recursive file scanning
- Terraform policy checks
- Kubernetes policy checks
- Modular rule engine
- CI-compatible exit codes

### Planned

- Markdown and JSON reports
- GitHub Actions integration
- Terraform HCL parser integration
- Kubernetes schema validation
- Custom rule configuration
- Severity filtering
- Ignore files
- SARIF report generation
- Additional Terraform and Kubernetes security rules

---

## Technologies

- Go
- Cobra
- YAML v3
- Standard Library

---

## Contributing

Contributions are welcome.

If you discover a bug or would like to add a new policy rule, feel free to open an issue or submit a pull request.

---

## License

MIT
