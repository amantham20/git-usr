package main

import (
	"encoding/json"
	"testing"
)

// TestGetConfigPath tests the config path generation
func TestGetConfigPath(t *testing.T) {
	path, err := getConfigPath()
	if err != nil {
		t.Fatalf("getConfigPath() failed: %v", err)
	}

	if path == "" {
		t.Error("getConfigPath() returned empty path")
	}
}

// TestGetProfileNames tests profile name extraction
func TestGetProfileNames(t *testing.T) {
	profiles := map[string]Profile{
		"work":     {Name: "John", Email: "john@work.com"},
		"personal": {Name: "Jane", Email: "jane@personal.com"},
	}

	names := getProfileNames(profiles)

	if names == "" {
		t.Error("getProfileNames returned empty string")
	}
}

// TestProfileJSONSerialization tests JSON serialization
func TestProfileJSONSerialization(t *testing.T) {
	profile := Profile{
		Name:  "Test User",
		Email: "test@example.com",
	}

	data, err := json.Marshal(profile)
	if err != nil {
		t.Fatalf("Failed to marshal profile: %v", err)
	}

	var loaded Profile
	err = json.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("Failed to unmarshal profile: %v", err)
	}

	if loaded.Name != profile.Name || loaded.Email != profile.Email {
		t.Error("Profile data mismatch after serialization")
	}
}

// TestGenerateCompletionBash tests bash completion generation
func TestGenerateCompletionBash(t *testing.T) {
	profiles := []string{"work", "personal"}
	completion := getBashCompletion(profiles)

	if completion == "" {
		t.Error("Bash completion is empty")
	}

	if !contains(completion, "work") || !contains(completion, "personal") {
		t.Error("Bash completion missing profiles")
	}
}

// TestGenerateCompletionZsh tests zsh completion generation
func TestGenerateCompletionZsh(t *testing.T) {
	profiles := []string{"work"}
	completion := getZshCompletion(profiles)

	if completion == "" {
		t.Error("Zsh completion is empty")
	}

	if !contains(completion, "#compdef") {
		t.Error("Zsh completion missing #compdef directive")
	}
}

// TestGenerateCompletionFish tests fish completion generation
func TestGenerateCompletionFish(t *testing.T) {
	profiles := []string{"work"}
	completion := getFishCompletion(profiles)

	if completion == "" {
		t.Error("Fish completion is empty")
	}

	if !contains(completion, "complete -c git-usr") {
		t.Error("Fish completion missing complete command")
	}
}

// TestGenerateCompletionPowershell tests powershell completion generation
func TestGenerateCompletionPowershell(t *testing.T) {
	profiles := []string{"work"}
	completion := getPowershellCompletion(profiles)

	if completion == "" {
		t.Error("PowerShell completion is empty")
	}

	if !contains(completion, "Register-ArgumentCompleter") {
		t.Error("PowerShell completion missing Register-ArgumentCompleter")
	}
}

// TestEmptyProfileHandling tests handling of empty profile sets
func TestEmptyProfileHandling(t *testing.T) {
	emptyProfiles := map[string]Profile{}
	names := getProfileNames(emptyProfiles)

	if names != "" {
		t.Errorf("Expected empty string for empty profiles, got: %s", names)
	}
}

// Helper function
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
