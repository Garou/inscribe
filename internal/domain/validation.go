package domain

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ValidatorFunc validates a string value and returns an error if invalid.
type ValidatorFunc func(value string) error

var validators = map[string]ValidatorFunc{
	"dns-name": ValidateDNSName,
	"integer":  ValidateInteger,
	"string":   ValidateString,
	"port":     ValidatePort,
	"memory":   ValidateMemory,
	"cpu":      ValidateCPU,
}

// GetValidator returns a validator for the given type, or an error if unknown.
func GetValidator(validationType string) (ValidatorFunc, error) {
	v, ok := validators[validationType]
	if !ok {
		return nil, fmt.Errorf("unknown validation type: %q", validationType)
	}
	return v, nil
}

var dnsNameRegexp = regexp.MustCompile(`^[a-z0-9]([a-z0-9\-]{0,61}[a-z0-9])?$`)

// ValidateDNSName checks that value is a valid RFC 1123 DNS label.
func ValidateDNSName(value string) error {
	if value == "" {
		return fmt.Errorf("dns name cannot be empty")
	}
	if len(value) > 63 {
		return fmt.Errorf("dns name must be 63 characters or fewer")
	}
	if !dnsNameRegexp.MatchString(value) {
		return fmt.Errorf("must be a valid dns name: lowercase alphanumeric and hyphens, must start and end with alphanumeric")
	}
	return nil
}

// ValidateInteger checks that value is a valid integer.
func ValidateInteger(value string) error {
	if value == "" {
		return fmt.Errorf("integer value cannot be empty")
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("must be a valid integer")
	}
	if n < 0 {
		return fmt.Errorf("must be a non-negative integer")
	}
	return nil
}

// ValidateString checks that value is a non-empty string.
func ValidateString(value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("value cannot be empty")
	}
	return nil
}

// ValidatePort checks that value is a valid port number (1-65535).
func ValidatePort(value string) error {
	if value == "" {
		return fmt.Errorf("port cannot be empty")
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("port must be a valid integer")
	}
	if n < 1 || n > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	return nil
}

var memoryRegexp = regexp.MustCompile(`^[0-9]+(\.[0-9]+)?(Ki|Mi|Gi|Ti|Pi|Ei|k|M|G|T|P|E)?$`)

// ValidateMemory checks that value is a valid Kubernetes memory resource quantity.
func ValidateMemory(value string) error {
	if value == "" {
		return fmt.Errorf("memory value cannot be empty")
	}
	if !memoryRegexp.MatchString(value) {
		return fmt.Errorf("must be a valid memory quantity (e.g., 256Mi, 1Gi)")
	}
	return nil
}

var cpuRegexp = regexp.MustCompile(`^[0-9]+(\.[0-9]+)?m?$`)

// ValidateCPU checks that value is a valid Kubernetes CPU resource quantity.
func ValidateCPU(value string) error {
	if value == "" {
		return fmt.Errorf("cpu value cannot be empty")
	}
	if !cpuRegexp.MatchString(value) {
		return fmt.Errorf("must be a valid cpu quantity (e.g., 100m, 0.5, 2)")
	}
	return nil
}
