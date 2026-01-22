package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const version = "1.0.0"

// Profile represents a git user profile
type Profile struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Config holds all user profiles
type Config struct {
	Profiles map[string]Profile `json:"profiles"`
}

// getConfigPath returns the path to the configuration file
func getConfigPath() (string, error) {
	var configDir string

	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		configDir = filepath.Join(appData, "git-usr")
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(home, ".config", "git-usr")
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(configDir, "profiles.json"), nil
}

// loadProfiles loads profiles from the config file
func loadProfiles() (map[string]Profile, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	// If file doesn't exist, create default profiles
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultProfiles := map[string]Profile{
			"work": {
				Name:  "Your Work Name",
				Email: "you@work.com",
			},
			"personal": {
				Name:  "Your Personal Name",
				Email: "you@personal.com",
			},
		}
		if err := saveProfiles(defaultProfiles); err != nil {
			return nil, err
		}
		return defaultProfiles, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var profiles map[string]Profile
	if err := json.Unmarshal(data, &profiles); err != nil {
		return nil, err
	}

	return profiles, nil
}

// saveProfiles saves profiles to the config file
func saveProfiles(profiles map[string]Profile) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(profiles, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// setGitConfig sets git user name and email
func setGitConfig(name, email, scope string) error {
	cmd := exec.Command("git", "config", "--"+scope, "user.name", name)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set user.name: %w", err)
	}

	cmd = exec.Command("git", "config", "--"+scope, "user.email", email)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set user.email: %w", err)
	}

	return nil
}

// getCurrentGitConfig gets the current git user name and email
func getCurrentGitConfig() (string, string, error) {
	nameCmd := exec.Command("git", "config", "user.name")
	nameOut, err := nameCmd.Output()
	if err != nil {
		return "", "", nil // Not an error, just no config
	}

	emailCmd := exec.Command("git", "config", "user.email")
	emailOut, err := emailCmd.Output()
	if err != nil {
		return "", "", nil
	}

	return strings.TrimSpace(string(nameOut)), strings.TrimSpace(string(emailOut)), nil
}

// listProfiles lists all available profiles
func listProfiles() error {
	profiles, err := loadProfiles()
	if err != nil {
		return err
	}

	currentName, currentEmail, _ := getCurrentGitConfig()

	fmt.Println("\n📋 Available profiles:")
	fmt.Println(strings.Repeat("-", 50))

	for name, profile := range profiles {
		isCurrent := profile.Name == currentName && profile.Email == currentEmail
		marker := "   "
		if isCurrent {
			marker = "👉 "
		}
		fmt.Printf("%s%s\n", marker, name)
		fmt.Printf("   Name:  %s\n", profile.Name)
		fmt.Printf("   Email: %s\n", profile.Email)
		fmt.Println()
	}

	return nil
}

// switchProfile switches to a specific profile
func switchProfile(profileName, scope string) error {
	profiles, err := loadProfiles()
	if err != nil {
		return err
	}

	profile, exists := profiles[profileName]
	if !exists {
		fmt.Printf("❌ Profile '%s' not found!\n", profileName)
		fmt.Println("\nAvailable profiles:", getProfileNames(profiles))
		fmt.Println("\nUse 'git usr add' to create a new profile")
		return fmt.Errorf("profile not found")
	}

	if err := setGitConfig(profile.Name, profile.Email, scope); err != nil {
		return err
	}

	scopeText := "for this repository"
	if scope == "global" {
		scopeText = "globally"
	}

	fmt.Printf("✅ Switched to '%s' profile %s\n", profileName, scopeText)
	fmt.Printf("   Name:  %s\n", profile.Name)
	fmt.Printf("   Email: %s\n", profile.Email)

	return nil
}

// addProfile adds or updates a profile
func addProfile(profileName, name, email string) error {
	profiles, err := loadProfiles()
	if err != nil {
		return err
	}

	// If profile exists and no new data provided
	if _, exists := profiles[profileName]; exists && (name == "" || email == "") {
		fmt.Printf("Profile '%s' already exists:\n", profileName)
		fmt.Printf("  Name:  %s\n", profiles[profileName].Name)
		fmt.Printf("  Email: %s\n", profiles[profileName].Email)
		fmt.Println("\nTo update, provide both name and email.")
		return nil
	}

	// Interactive mode if name/email not provided
	if name == "" {
		fmt.Print("Enter name: ")
		fmt.Scanln(&name)
	}
	if email == "" {
		fmt.Print("Enter email: ")
		fmt.Scanln(&email)
	}

	if name == "" || email == "" {
		return fmt.Errorf("❌ Name and email are required!")
	}

	profiles[profileName] = Profile{
		Name:  name,
		Email: email,
	}

	if err := saveProfiles(profiles); err != nil {
		return err
	}

	fmt.Printf("✅ Profile '%s' saved!\n", profileName)
	fmt.Printf("   Name:  %s\n", name)
	fmt.Printf("   Email: %s\n", email)
	fmt.Printf("\nUse: git usr %s\n", profileName)

	return nil
}

// removeProfile removes a profile
func removeProfile(profileName string) error {
	profiles, err := loadProfiles()
	if err != nil {
		return err
	}

	if _, exists := profiles[profileName]; !exists {
		return fmt.Errorf("❌ Profile '%s' not found!", profileName)
	}

	delete(profiles, profileName)

	if err := saveProfiles(profiles); err != nil {
		return err
	}

	fmt.Printf("✅ Profile '%s' removed!\n", profileName)
	return nil
}

// showCurrent shows the current git configuration
func showCurrent() error {
	name, email, err := getCurrentGitConfig()
	if err != nil {
		return err
	}

	if name != "" && email != "" {
		fmt.Println("\n📝 Current git configuration:")
		fmt.Printf("   Name:  %s\n", name)
		fmt.Printf("   Email: %s\n", email)
	} else {
		fmt.Println("❌ No git configuration found in this repository")
	}

	return nil
}

// showHelp displays help information
func showHelp() {
	configPath, _ := getConfigPath()
	
	fmt.Println(`
🔧 Git User Profile Switcher

Usage:
  git usr <profile>              Switch to profile (local scope)
  git usr <profile> --global     Switch to profile (global scope)
  git usr list                   List all profiles
  git usr add <profile>          Add/update a profile (interactive)
  git usr add <profile> "Name" "email@example.com"
  git usr remove <profile>       Remove a profile
  git usr current                Show current git config
  git usr completion [bash|zsh|fish|powershell]  Generate completion script
  git usr version                Show version information
  git usr help                   Show this help

Examples:
  git usr work                   Switch to work profile (local)
  git usr personal --global      Switch to personal profile (global)
  git usr add work "John Doe" "john@company.com"
  git usr list                   List all available profiles

Config location: ` + configPath)
}

// showVersion displays version information
func showVersion() {
	fmt.Print(`
            __
           / _)
    .-^^^-/ /
 __/       /
<__.|_|-|_|

git-usr version ` + version + `
Made by Aman (thammina@msu.edu)
`)
}

// getProfileNames returns a comma-separated list of profile names
func getProfileNames(profiles map[string]Profile) string {
	names := make([]string, 0, len(profiles))
	for name := range profiles {
		names = append(names, name)
	}
	return strings.Join(names, ", ")
}

// generateCompletion generates shell completion scripts
func generateCompletion(shell string) error {
	profiles, err := loadProfiles()
	if err != nil {
		return err
	}

	profileNames := make([]string, 0, len(profiles))
	for name := range profiles {
		profileNames = append(profileNames, name)
	}

	switch shell {
	case "bash":
		fmt.Println(getBashCompletion(profileNames))
	case "zsh":
		fmt.Println(getZshCompletion(profileNames))
	case "fish":
		fmt.Println(getFishCompletion(profileNames))
	case "powershell":
		fmt.Println(getPowershellCompletion(profileNames))
	default:
		return fmt.Errorf("❌ Unsupported shell: %s. Supported: bash, zsh, fish, powershell", shell)
	}

	return nil
}

func getBashCompletion(profiles []string) string {
	return `# bash completion for git-usr
_git_usr() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    # Main commands
    local commands="list current add remove help version completion ` + strings.Join(profiles, " ") + `"
    
    # Completion for subcommands
    case "${prev}" in
        completion)
            COMPREPLY=( $(compgen -W "bash zsh fish powershell" -- ${cur}) )
            return 0
            ;;
        remove)
            COMPREPLY=( $(compgen -W "` + strings.Join(profiles, " ") + `" -- ${cur}) )
            return 0
            ;;
        *)
            ;;
    esac

    # Complete with available commands and profiles
    COMPREPLY=( $(compgen -W "${commands} --global" -- ${cur}) )
    return 0
}

complete -F _git_usr git-usr

# Installation: Add this to ~/.bashrc or ~/.bash_completion
# Or save to /etc/bash_completion.d/git-usr`
}

func getZshCompletion(profiles []string) string {
	profileList := strings.Join(profiles, " ")
	return `#compdef git-usr

_git_usr() {
    local -a commands profiles
    commands=(
        'list:List all profiles'
        'current:Show current git config'
        'add:Add or update a profile'
        'remove:Remove a profile'
        'version:Show version information'
        'help:Show help'
        'completion:Generate completion script'
    )
    
    profiles=(` + profileList + `)

    _arguments -C \
        '1: :->command' \
        '2: :->args' \
        '*::arg:->args' \
        '--global[Apply globally]'

    case $state in
        command)
            _describe -t commands 'git-usr commands' commands
            _describe -t profiles 'profiles' profiles
            ;;
        args)
            case $words[1] in
                completion)
                    _values 'shell' bash zsh fish powershell
                    ;;
                remove)
                    _describe -t profiles 'profiles' profiles
                    ;;
            esac
            ;;
    esac
}

_git_usr "$@"

# Installation: Save to a file in $fpath, e.g., ~/.zsh/completions/_git-usr
# Then add to ~/.zshrc: fpath=(~/.zsh/completions $fpath) && autoload -U compinit && compinit`
}

func getFishCompletion(profiles []string) string {
	completions := `# fish completion for git-usr

# Main commands
complete -c git-usr -f -n "__fish_use_subcommand" -a "list" -d "List all profiles"
complete -c git-usr -f -n "__fish_use_subcommand" -a "current" -d "Show current git config"
complete -c git-usr -f -n "__fish_use_subcommand" -a "add" -d "Add or update a profile"
complete -c git-usr -f -n "__fish_use_subcommand" -a "remove" -d "Remove a profile"
complete -c git-usr -f -n "__fish_use_subcommand" -a "version" -d "Show version information"
complete -c git-usr -f -n "__fish_use_subcommand" -a "help" -d "Show help"
complete -c git-usr -f -n "__fish_use_subcommand" -a "completion" -d "Generate completion script"

# Profiles
`
	for _, profile := range profiles {
		completions += fmt.Sprintf("complete -c git-usr -f -n \"__fish_use_subcommand\" -a \"%s\" -d \"Switch to %s profile\"\n", profile, profile)
	}

	completions += `
# Completion for completion subcommand
complete -c git-usr -f -n "__fish_seen_subcommand_from completion" -a "bash zsh fish powershell"

# Completion for remove subcommand
`
	for _, profile := range profiles {
		completions += fmt.Sprintf("complete -c git-usr -f -n \"__fish_seen_subcommand_from remove\" -a \"%s\"\n", profile)
	}

	completions += `
# Global flag
complete -c git-usr -l global -d "Apply globally"

# Installation: Save to ~/.config/fish/completions/git-usr.fish`

	return completions
}

func getPowershellCompletion(profiles []string) string {
	profileList := "'" + strings.Join(profiles, "', '") + "'"
	return `# PowerShell completion for git-usr

Register-ArgumentCompleter -Native -CommandName git-usr -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)

    $commands = @('list', 'current', 'add', 'remove', 'version', 'help', 'completion')
    $profiles = @(` + profileList + `)
    $shells = @('bash', 'zsh', 'fish', 'powershell')

    $tokens = $commandAst.ToString() -split '\s+'
    
    if ($tokens.Count -eq 2) {
        # Complete main commands and profiles
        $allOptions = $commands + $profiles + @('--global')
        $allOptions | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
            [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
        }
    }
    elseif ($tokens.Count -eq 3) {
        switch ($tokens[1]) {
            'completion' {
                $shells | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
                    [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
                }
            }
            'remove' {
                $profiles | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
                    [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
                }
            }
        }
    }
}

# Installation: Add this to your PowerShell profile ($PROFILE)
# Or dot-source this file: . path\to\git-usr-completion.ps1`
}

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	command := os.Args[1]
	scope := "local"

	// Check for --global flag
	for _, arg := range os.Args {
		if arg == "--global" {
			scope = "global"
			break
		}
	}

	var err error

	switch command {
	case "help", "--help", "-h":
		showHelp()

	case "version", "--version", "-v":
		showVersion()

	case "list":
		err = listProfiles()

	case "current":
		err = showCurrent()

	case "add":
		if len(os.Args) < 3 {
			fmt.Println("❌ Profile name required!")
			fmt.Println("Usage: git usr add <profile> [name] [email]")
			return
		}
		profileName := os.Args[2]
		name := ""
		email := ""
		if len(os.Args) > 3 {
			name = os.Args[3]
		}
		if len(os.Args) > 4 {
			email = os.Args[4]
		}
		err = addProfile(profileName, name, email)

	case "remove":
		if len(os.Args) < 3 {
			fmt.Println("❌ Profile name required!")
			fmt.Println("Usage: git usr remove <profile>")
			return
		}
		err = removeProfile(os.Args[2])

	case "completion":
		if len(os.Args) < 3 {
			fmt.Println("❌ Shell type required!")
			fmt.Println("Usage: git usr completion [bash|zsh|fish|powershell]")
			return
		}
		err = generateCompletion(os.Args[2])

	default:
		// Assume it's a profile name
		err = switchProfile(command, scope)
	}

	if err != nil {
		os.Exit(1)
	}
}
