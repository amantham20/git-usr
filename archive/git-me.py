#!/usr/bin/env python3
"""
Git Profile Switcher - Quickly switch between different git user profiles
"""
import json
import os
import subprocess
import sys
from pathlib import Path


def get_config_path():
    """Get the path to the config file"""
    if sys.platform == "win32":
        config_dir = Path(os.environ.get("APPDATA", Path.home() / "AppData/Roaming")) / "git-me"
    else:
        config_dir = Path.home() / ".config" / "git-me"
    
    config_dir.mkdir(parents=True, exist_ok=True)
    return config_dir / "profiles.json"


def load_profiles():
    """Load profiles from config file"""
    config_path = get_config_path()
    
    if not config_path.exists():
        # Create default profiles
        default_profiles = {
            "work": {
                "name": "Your Work Name",
                "email": "you@work.com"
            },
            "personal": {
                "name": "Your Personal Name",
                "email": "you@personal.com"
            }
        }
        save_profiles(default_profiles)
        return default_profiles
    
    with open(config_path, 'r') as f:
        return json.load(f)


def save_profiles(profiles):
    """Save profiles to config file"""
    config_path = get_config_path()
    with open(config_path, 'w') as f:
        json.dump(profiles, f, indent=2)


def set_git_config(name, email, scope='local'):
    """Set git user name and email"""
    try:
        subprocess.run(['git', 'config', f'--{scope}', 'user.name', name], check=True)
        subprocess.run(['git', 'config', f'--{scope}', 'user.email', email], check=True)
        return True
    except subprocess.CalledProcessError as e:
        print(f"Error setting git config: {e}")
        return False


def get_current_git_config():
    """Get current git user name and email"""
    try:
        name = subprocess.check_output(['git', 'config', 'user.name'], text=True).strip()
        email = subprocess.check_output(['git', 'config', 'user.email'], text=True).strip()
        return name, email
    except subprocess.CalledProcessError:
        return None, None


def list_profiles():
    """List all available profiles"""
    profiles = load_profiles()
    current_name, current_email = get_current_git_config()
    
    print("\n📋 Available profiles:")
    print("-" * 50)
    
    for profile_name, profile_data in profiles.items():
        is_current = (profile_data['name'] == current_name and 
                     profile_data['email'] == current_email)
        marker = "👉 " if is_current else "   "
        print(f"{marker}{profile_name}")
        print(f"   Name:  {profile_data['name']}")
        print(f"   Email: {profile_data['email']}")
        print()


def switch_profile(profile_name, scope='local'):
    """Switch to a specific profile"""
    profiles = load_profiles()
    
    if profile_name not in profiles:
        print(f"❌ Profile '{profile_name}' not found!")
        print(f"\nAvailable profiles: {', '.join(profiles.keys())}")
        print(f"\nUse 'git-me add' to create a new profile")
        return False
    
    profile = profiles[profile_name]
    
    if set_git_config(profile['name'], profile['email'], scope):
        scope_text = "globally" if scope == "global" else "for this repository"
        print(f"✅ Switched to '{profile_name}' profile {scope_text}")
        print(f"   Name:  {profile['name']}")
        print(f"   Email: {profile['email']}")
        return True
    
    return False


def add_profile(profile_name, name=None, email=None):
    """Add or update a profile"""
    profiles = load_profiles()
    
    if profile_name in profiles and not (name and email):
        print(f"Profile '{profile_name}' already exists:")
        print(f"  Name:  {profiles[profile_name]['name']}")
        print(f"  Email: {profiles[profile_name]['email']}")
        print("\nTo update, provide both name and email.")
        return
    
    # Interactive mode if name/email not provided
    if not name:
        name = input("Enter name: ").strip()
    if not email:
        email = input("Enter email: ").strip()
    
    if not name or not email:
        print("❌ Name and email are required!")
        return
    
    profiles[profile_name] = {
        "name": name,
        "email": email
    }
    
    save_profiles(profiles)
    print(f"✅ Profile '{profile_name}' saved!")
    print(f"   Name:  {name}")
    print(f"   Email: {email}")
    print(f"\nUse: git-me {profile_name}")


def remove_profile(profile_name):
    """Remove a profile"""
    profiles = load_profiles()
    
    if profile_name not in profiles:
        print(f"❌ Profile '{profile_name}' not found!")
        return
    
    del profiles[profile_name]
    save_profiles(profiles)
    print(f"✅ Profile '{profile_name}' removed!")


def show_help():
    """Show help message"""
    print("""
🔧 Git Profile Switcher

Usage:
  git-me <profile>              Switch to profile (local scope)
  git-me <profile> --global     Switch to profile (global scope)
  git-me list                   List all profiles
  git-me add <profile>          Add/update a profile (interactive)
  git-me add <profile> "Name" "email@example.com"
  git-me remove <profile>       Remove a profile
  git-me current                Show current git config
  git-me help                   Show this help

Examples:
  git-me work                   Switch to work profile (local)
  git-me personal --global      Switch to personal profile (global)
  git-me add work "John Doe" "john@company.com"
  git-me list                   List all available profiles

Config location: {}
""".format(get_config_path()))


def show_current():
    """Show current git configuration"""
    name, email = get_current_git_config()
    
    if name and email:
        print("\n📝 Current git configuration:")
        print(f"   Name:  {name}")
        print(f"   Email: {email}")
    else:
        print("❌ No git configuration found in this repository")


def main():
    """Main entry point"""
    if len(sys.argv) < 2:
        show_help()
        return
    
    command = sys.argv[1]
    
    # Check for global flag
    scope = 'global' if '--global' in sys.argv else 'local'
    
    if command in ['help', '--help', '-h']:
        show_help()
    
    elif command == 'list':
        list_profiles()
    
    elif command == 'current':
        show_current()
    
    elif command == 'add':
        if len(sys.argv) < 3:
            print("❌ Profile name required!")
            print("Usage: git-me add <profile> [name] [email]")
            return
        
        profile_name = sys.argv[2]
        name = sys.argv[3] if len(sys.argv) > 3 else None
        email = sys.argv[4] if len(sys.argv) > 4 else None
        add_profile(profile_name, name, email)
    
    elif command == 'remove':
        if len(sys.argv) < 3:
            print("❌ Profile name required!")
            print("Usage: git-me remove <profile>")
            return
        
        profile_name = sys.argv[2]
        remove_profile(profile_name)
    
    else:
        # Assume it's a profile name
        switch_profile(command, scope)


if __name__ == "__main__":
    main()
