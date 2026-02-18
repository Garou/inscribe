# Inscribe

An interactive CLI/TUI tool for generating Kubernetes manifest files from templates. Inscribe connects to live clusters to auto-list values like namespaces and CNPG clusters, validates input, and produces correctly-structured YAML.

## Install

Requires Go 1.25+.

```sh
git clone <repo-url> && cd inscribe
make build
```

The binary is output to `bin/inscribe`.

## Quick Start

```sh
# Generate a CNPG cluster manifest interactively (launches TUI wizard)
./bin/inscribe cluster cnpg

# Generate non-interactively with all flags
./bin/inscribe cluster cnpg \
  --name=mydb \
  --namespace=default \
  --instances=3 \
  --cnpg-resource-templates="Production - 4Gi/2CPU" \
  --context=minikube \
  --filename=mydb-cluster.yaml

# Run a parent command to auto-select or pick a template
./bin/inscribe cluster
```

If all required flags are provided, Inscribe renders the manifest directly. If any are missing, it launches an interactive TUI wizard with the provided values pre-filled.

## Commands

```
inscribe
├── cluster                  # Generate cluster manifests
│   └── cnpg                 # CNPG PostgreSQL Cluster
├── backup                   # Generate backup manifests
│   └── cnpg                 # CNPG One-Off Backup
├── scheduled-backup         # Generate scheduled backup manifests
│   └── cnpg                 # CNPG Scheduled Backup
└── env [path]               # Output shell config for template directory
```

### `inscribe cluster cnpg`

Generates a CloudNativePG PostgreSQL cluster manifest.

| Flag | Description |
|---|---|
| `--name` | Cluster name (must be a valid DNS name) |
| `--namespace` | Kubernetes namespace (auto-listed from cluster if omitted) |
| `--instances` | Number of PostgreSQL instances |
| `--cnpg-resource-templates` | Resource profile: `"Production - 4Gi/2CPU"`, `"QA - 2Gi/1CPU"`, or `"Test - 512Mi/500m"` |
| `--context` | Kubernetes context to use |
| `--filename` | Output filename |

### `inscribe backup cnpg`

Generates a CloudNativePG one-off backup manifest.

| Flag | Description |
|---|---|
| `--name` | Backup name (must be a valid DNS name) |
| `--namespace` | Kubernetes namespace (auto-listed from cluster if omitted) |
| `--cnpg-clusters` | CNPG cluster to back up (auto-listed from cluster if omitted) |
| `--backup-methods` | Backup method: `barmanObjectStore` or `volumeSnapshot` |
| `--context` | Kubernetes context to use |
| `--filename` | Output filename |

### `inscribe scheduled-backup cnpg`

Generates a CloudNativePG scheduled backup manifest.

| Flag | Description |
|---|---|
| `--name` | Backup name (must be a valid DNS name) |
| `--namespace` | Kubernetes namespace (auto-listed from cluster if omitted) |
| `--schedule` | Cron schedule expression (e.g. `"0 0 * * *"`) |
| `--cnpg-clusters` | CNPG cluster to back up (auto-listed from cluster if omitted) |
| `--backup-methods` | Backup method: `barmanObjectStore` or `volumeSnapshot` |
| `--context` | Kubernetes context to use |
| `--filename` | Output filename |

### Global Flags

| Flag | Env Variable | Default | Description |
|---|---|---|---|
| `--template-dir` | `INSCRIBE_TEMPLATE_DIR` | `templates` | Path to template directory |
| `-o`, `--output-dir` | | `.` | Output directory for generated manifests |

### `inscribe env`

Outputs a shell export statement for `INSCRIBE_TEMPLATE_DIR`. Add to your shell profile for persistent configuration:

```sh
# Set for current session
eval "$(inscribe env /path/to/templates)"

# Add to ~/.zshrc or ~/.bashrc
echo 'eval "$(inscribe env /path/to/your/templates)"' >> ~/.zshrc
```

## Templates

Templates live in the directory specified by `--template-dir` or `INSCRIBE_TEMPLATE_DIR`. Inscribe scans the directory recursively for `.yaml`/`.yml` files with an `inscribe:` header comment.

### Template Types

**Main template** — defines a manifest with placeholder fields:

```yaml
{{/* inscribe: type="template" name="cnpg-cluster" command="cluster cnpg" description="CNPG PostgreSQL Cluster" */}}
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata:
  name: {{ input "name" "dns-name" }}
  namespace: {{ autoList "namespace" }}
spec:
  instances: {{ input "instances" "integer" }}
  resources:
{{ templateGroup "cnpg-resource-templates" | indent 4 }}
```

**Sub-template** — a reusable fragment selected by the user from a group:

```yaml
{{/* inscribe: type="sub-template" group="cnpg-resource-templates" description="Production - 4Gi/2CPU" */}}
requests:
  memory: "4Gi"
  cpu: "2"
limits:
  memory: "4Gi"
  cpu: "2"
```

**Static list** — a predefined set of values to pick from:

```yaml
{{/* inscribe: type="list" name="backup-methods" */}}
- barmanObjectStore
- volumeSnapshot
```

### Template Functions

| Function | Description | Example |
|---|---|---|
| `input "name" "validation"` | User-provided field with validation | `{{ input "name" "dns-name" }}` |
| `autoList "source"` | Auto-populated from Kubernetes | `{{ autoList "namespace" }}` |
| `templateGroup "group"` | Pick from sub-template group | `{{ templateGroup "cnpg-resource-templates" \| indent 4 }}` |
| `staticList "name"` | Pick from static list | `{{ staticList "backup-methods" }}` |
| `indent N` | Indent piped content by N spaces | `{{ templateGroup "grp" \| indent 4 }}` |

### Validation Types

Used with `input` fields:

| Type | Rules |
|---|---|
| `dns-name` | RFC 1123 DNS label: lowercase alphanumeric and hyphens, 1-63 chars |
| `integer` | Non-negative integer |
| `string` | Non-empty string |
| `port` | Integer between 1 and 65535 |
| `memory` | Kubernetes memory quantity (e.g. `256Mi`, `4Gi`) |
| `cpu` | Kubernetes CPU quantity (e.g. `100m`, `0.5`, `2`) |
| `cron-schedule` | Standard 5-field cron expression (e.g. `0 0 * * *`, `*/5 * * * 1-5`) |
| `filename` | Valid output filename (no path separators or directory traversal) |
| `path` | Non-empty directory path |

### Auto-Detect Sources

Used with `autoList` fields:

| Source | Description |
|---|---|
| `namespace` | Lists namespaces from the selected Kubernetes context |
| `cnpg-clusters` | Lists CNPG clusters from the selected context and namespace |

## Writing Custom Templates

1. Create a `.yaml` file in your template directory with an `inscribe:` header
2. Use template functions for dynamic fields
3. Create sub-templates and static lists as needed
4. The template will automatically appear in the CLI based on its `command` field

The `command` header field maps to the CLI structure. For example, `command="cluster cnpg"` maps to `inscribe cluster cnpg`.

## Development

```sh
make build          # Build binary to bin/inscribe
make test           # Run unit tests
make test-coverage  # Run tests with coverage report
make lint           # Run golangci-lint
make clean          # Remove build artifacts
```

### Running Integration Tests

Integration tests run against a live Kubernetes cluster:

```sh
go test ./internal/kubernetes/ -tags=integration -v
```

## Project Structure

```
inscribe/
├── cmd/inscribe/          # Entry point
├── internal/
│   ├── domain/            # Types, validation, interfaces
│   ├── engine/            # Template parsing, registry, rendering
│   ├── kubernetes/        # K8s client (contexts, namespaces, CRDs)
│   ├── tui/               # Interactive wizard (huh-based)
│   │   └── components/    # Atomic design: atoms, molecules, organisms
│   ├── cli/               # Cobra commands and bridge logic
│   └── output/            # Manifest file writer
└── templates/cnpg/        # Bundled CNPG templates
```
