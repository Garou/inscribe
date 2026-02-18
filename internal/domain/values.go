package domain

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ValueObject represents a validated domain value.
type ValueObject interface {
	String() string
}

// ParseValue creates a validated value object from a type name and raw string.
// Returns an error if the value is invalid for the given type.
func ParseValue(typeName string, value string) (ValueObject, error) {
	switch typeName {
	case "dns-name":
		return NewDNSName(value)
	case "integer":
		return NewInteger(value)
	case "string":
		return NewNonEmptyString(value)
	case "port":
		return NewPort(value)
	case "memory":
		return NewMemory(value)
	case "cpu":
		return NewCPU(value)
	case "cron-schedule":
		return NewCronSchedule(value)
	case "filename":
		return NewFilename(value)
	case "path":
		return NewPath(value)
	default:
		return nil, fmt.Errorf("unknown validation type: %q", typeName)
	}
}

// DNSName is a valid RFC 1123 DNS label.
type DNSName struct{ value string }

var dnsNameRegexp = regexp.MustCompile(`^[a-z0-9]([a-z0-9\-]{0,61}[a-z0-9])?$`)

func NewDNSName(s string) (DNSName, error) {
	if s == "" {
		return DNSName{}, fmt.Errorf("dns name cannot be empty")
	}
	if len(s) > 63 {
		return DNSName{}, fmt.Errorf("dns name must be 63 characters or fewer")
	}
	if !dnsNameRegexp.MatchString(s) {
		return DNSName{}, fmt.Errorf("must be a valid dns name: lowercase alphanumeric and hyphens, must start and end with alphanumeric")
	}
	return DNSName{value: s}, nil
}

func (d DNSName) String() string { return d.value }

// Integer is a non-negative integer.
type Integer struct{ value string }

func NewInteger(s string) (Integer, error) {
	if s == "" {
		return Integer{}, fmt.Errorf("integer value cannot be empty")
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return Integer{}, fmt.Errorf("must be a valid integer")
	}
	if n < 0 {
		return Integer{}, fmt.Errorf("must be a non-negative integer")
	}
	return Integer{value: s}, nil
}

func (i Integer) String() string { return i.value }

// NonEmptyString is a non-empty string value.
type NonEmptyString struct{ value string }

func NewNonEmptyString(s string) (NonEmptyString, error) {
	if strings.TrimSpace(s) == "" {
		return NonEmptyString{}, fmt.Errorf("value cannot be empty")
	}
	return NonEmptyString{value: s}, nil
}

func (n NonEmptyString) String() string { return n.value }

// Port is a valid port number (1-65535).
type Port struct{ value string }

func NewPort(s string) (Port, error) {
	if s == "" {
		return Port{}, fmt.Errorf("port cannot be empty")
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return Port{}, fmt.Errorf("port must be a valid integer")
	}
	if n < 1 || n > 65535 {
		return Port{}, fmt.Errorf("port must be between 1 and 65535")
	}
	return Port{value: s}, nil
}

func (p Port) String() string { return p.value }

// Memory is a valid Kubernetes memory resource quantity.
type Memory struct{ value string }

var memoryRegexp = regexp.MustCompile(`^[0-9]+(\.[0-9]+)?(Ki|Mi|Gi|Ti|Pi|Ei|k|M|G|T|P|E)?$`)

func NewMemory(s string) (Memory, error) {
	if s == "" {
		return Memory{}, fmt.Errorf("memory value cannot be empty")
	}
	if !memoryRegexp.MatchString(s) {
		return Memory{}, fmt.Errorf("must be a valid memory quantity (e.g., 256Mi, 1Gi)")
	}
	return Memory{value: s}, nil
}

func (m Memory) String() string { return m.value }

// CPU is a valid Kubernetes CPU resource quantity.
type CPU struct{ value string }

var cpuRegexp = regexp.MustCompile(`^[0-9]+(\.[0-9]+)?m?$`)

func NewCPU(s string) (CPU, error) {
	if s == "" {
		return CPU{}, fmt.Errorf("cpu value cannot be empty")
	}
	if !cpuRegexp.MatchString(s) {
		return CPU{}, fmt.Errorf("must be a valid cpu quantity (e.g., 100m, 0.5, 2)")
	}
	return CPU{value: s}, nil
}

func (c CPU) String() string { return c.value }

// CronSchedule is a valid 5-field standard cron expression.
type CronSchedule struct{ value string }

var cronFieldRegexp = regexp.MustCompile(`^(\*|[0-9]+(-[0-9]+)?)(\/[0-9]+)?(,(\*|[0-9]+(-[0-9]+)?)(\/[0-9]+)?)*$`)

func NewCronSchedule(s string) (CronSchedule, error) {
	if strings.TrimSpace(s) == "" {
		return CronSchedule{}, fmt.Errorf("cron schedule cannot be empty")
	}
	fields := strings.Fields(s)
	if len(fields) != 5 {
		return CronSchedule{}, fmt.Errorf("cron schedule must have exactly 5 fields (minute hour day month weekday)")
	}
	for _, field := range fields {
		if !cronFieldRegexp.MatchString(field) {
			return CronSchedule{}, fmt.Errorf("invalid cron field %q: must contain numbers, *, -, /, or ,", field)
		}
	}
	return CronSchedule{value: s}, nil
}

func (c CronSchedule) String() string { return c.value }

// Filename is a validated output filename (no path separators or traversal).
type Filename struct{ value string }

func NewFilename(s string) (Filename, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Filename{}, fmt.Errorf("filename cannot be empty")
	}
	if strings.Contains(s, "/") || strings.Contains(s, "\\") {
		return Filename{}, fmt.Errorf("filename must not contain path separators")
	}
	if strings.HasPrefix(s, "..") {
		return Filename{}, fmt.Errorf("filename must not contain path traversal")
	}
	return Filename{value: s}, nil
}

func (f Filename) String() string { return f.value }

// Path is a validated non-empty directory path.
type Path struct{ value string }

func NewPath(s string) (Path, error) {
	if strings.TrimSpace(s) == "" {
		return Path{}, fmt.Errorf("path cannot be empty")
	}
	return Path{value: s}, nil
}

func (p Path) String() string { return p.value }
