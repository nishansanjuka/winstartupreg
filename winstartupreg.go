package winstartupreg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// StartupRegistryType represents different startup registry locations
type StartupRegistryType int

const (
	CurrentUserRun StartupRegistryType = iota
	CurrentUserRunOnce
	AllUsersRun
	AllUsersRunOnce
)

// StartupEntry represents a Windows startup registry entry
type StartupEntry struct {
	Name    string
	Command string
}

// getRegistryPath returns the full registry path and root key for a given startup type
func getRegistryPath(registryType StartupRegistryType) (string, registry.Key) {
	switch registryType {
	case CurrentUserRun:
		return `Software\Microsoft\Windows\CurrentVersion\Run`, registry.CURRENT_USER
	case CurrentUserRunOnce:
		return `Software\Microsoft\Windows\CurrentVersion\RunOnce`, registry.CURRENT_USER
	case AllUsersRun:
		return `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.LOCAL_MACHINE
	case AllUsersRunOnce:
		return `SOFTWARE\Microsoft\Windows\CurrentVersion\RunOnce`, registry.LOCAL_MACHINE
	default:
		return `Software\Microsoft\Windows\CurrentVersion\Run`, registry.CURRENT_USER
	}
}

// AddStartupEntry adds an application to Windows startup registry
func AddStartupEntry(entry StartupEntry, registryType StartupRegistryType) error {
	// Validate input
	if entry.Name == "" {
		return fmt.Errorf("entry name cannot be empty")
	}

	// Normalize and validate command path
	fullPath, err := filepath.Abs(entry.Command)
	if err != nil {
		return fmt.Errorf("invalid command path: %w", err)
	}

	// Check if the executable exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("executable does not exist: %s", fullPath)
	}

	// Get registry path and root key
	keyPath, rootKey := getRegistryPath(registryType)

	// Open the registry key with write access
	k, err := registry.OpenKey(rootKey, keyPath, registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer k.Close()

	// Set the registry value
	err = k.SetStringValue(entry.Name, fullPath)
	if err != nil {
		return fmt.Errorf("failed to set registry value: %w", err)
	}

	return nil
}

// RemoveStartupEntry removes an application from Windows startup registry
func RemoveStartupEntry(entryName string, registryType StartupRegistryType) error {
	// Get registry path and root key
	keyPath, rootKey := getRegistryPath(registryType)

	// Attempt to open the registry key with write access
	k, err := registry.OpenKey(rootKey, keyPath, registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer k.Close()

	// Attempt to delete the value
	err = k.DeleteValue(entryName)
	if err != nil {
		// Check if the error indicates the value doesn't exist
		if strings.Contains(err.Error(), "The system cannot find the file specified") {
			return fmt.Errorf("startup entry '%s' not found in %s", entryName, keyPath)
		}
		return fmt.Errorf("failed to delete registry value: %w", err)
	}

	return nil
}

// SafeRemoveStartupEntry provides a comprehensive removal method
func SafeRemoveStartupEntry(entryName string) error {
	// List of registry types to check
	registryTypes := []StartupRegistryType{
		CurrentUserRun,
		CurrentUserRunOnce,
		AllUsersRun,
		AllUsersRunOnce,
	}

	var lastErr error
	var removedFromAny bool

	// Try to remove from all possible locations
	for _, registryType := range registryTypes {
		err := RemoveStartupEntry(entryName, registryType)
		if err == nil {
			removedFromAny = true
		} else {
			lastErr = err
		}
	}

	if !removedFromAny {
		return fmt.Errorf("failed to remove startup entry '%s' from any location: %w", entryName, lastErr)
	}

	return nil
}

// ListStartupEntries retrieves startup entries from a specific registry location
func ListStartupEntries(registryType StartupRegistryType) (map[string]string, error) {
	// Get registry path and root key
	keyPath, rootKey := getRegistryPath(registryType)

	// Open the registry key with read access
	k, err := registry.OpenKey(rootKey, keyPath, registry.QUERY_VALUE)
	if err != nil {
		return nil, fmt.Errorf("failed to open registry key: %w", err)
	}
	defer k.Close()

	// Get all value names
	valueNames, err := k.ReadValueNames(0)
	if err != nil {
		return nil, fmt.Errorf("failed to read value names: %w", err)
	}

	// Create a map to store startup entries
	entries := make(map[string]string)

	// Read each value
	for _, name := range valueNames {
		value, _, err := k.GetStringValue(name)
		if err == nil {
			entries[name] = value
		}
	}

	return entries, nil
}

// ListAllStartupEntries retrieves startup entries from all known locations
func ListAllStartupEntries() (map[StartupRegistryType]map[string]string, error) {
	// List of registry types to check
	registryTypes := []StartupRegistryType{
		CurrentUserRun,
		CurrentUserRunOnce,
		AllUsersRun,
		AllUsersRunOnce,
	}

	// Map to store all startup entries
	allEntries := make(map[StartupRegistryType]map[string]string)

	// Retrieve entries from each location
	for _, registryType := range registryTypes {
		entries, err := ListStartupEntries(registryType)
		if err == nil && len(entries) > 0 {
			allEntries[registryType] = entries
		}
	}

	return allEntries, nil
}
