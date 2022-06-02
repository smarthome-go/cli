#!/usr/bin/env python3
# This file is used to update the version number in all relevant places
# The SemVer (https://semver.org) versioning system is used.
import re

makefile_path = "./Makefile"
root_go_path = "./cmd/root.go"
readme_path = "./README.md"

# Extract old version from root.go
with open(root_go_path, "r") as root_go:
    content = root_go.read()
    old_version = content.split('const Version = "')[1].split('"\n')[0]
    print(f"Found old version in {root_go_path}: {old_version}")

# Attempt to read new version from user input
try:
    VERSION = input(
        f"Current version: {old_version}\nNew version (without 'v' prefix): ")
except KeyboardInterrupt:
    print("\nCanceled by user")
    quit()

if VERSION == "":
    VERSION = old_version

if not re.match(r"^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$", VERSION):
    print(
        f"\x1b[31mThe version: '{VERSION}' is not a valid SemVer version.\x1b[0m")
    quit()

# Update version in root.go
with open(root_go_path, 'w') as root_go:
    root_go.write(content.replace(old_version, VERSION))

# Update version in Makefile
with open(makefile_path, 'r') as makefile:
    content = makefile.read()
    old_version = content.split("version := ")[1].split('\n')[0]
    print(f"Found old version in {makefile_path}: {old_version}")

with open(makefile_path, 'w') as makefile:
    makefile.write(content.replace(old_version, VERSION))

with open(readme_path, "r") as readme:
    content = readme.read()
    old_version = content.split('## Installation of v')[
        1].split(' (for Linux/AMD64)\n')[0]
    print(f"Found old version in {readme_path}: {old_version}")

with open(readme_path, "r") as readme:
    content = readme.read()
    old_version = content.split('cd /tmp && wget https://github.com/smarthome-go/cli/releases/download/v')[
        1].split('/homescript_linux_amd64.tar.gz')[0]
    print(f"Found old version 2 in {readme_path}: {old_version}")

with open(readme_path, "w") as readme:
    readme.write(content.replace(old_version, VERSION))

print(f"Version has been changed from '{old_version}' -> '{VERSION}'")
