#!/bin/bash

# Copyright (c) 2025 A Bit of Help, Inc.

# Script to replace "System" with "Service" in filenames and file content

echo "Implementing changes to replace 'System' with 'Service'..."
echo

# 1. Files to be renamed
echo "Files being renamed:"
echo "----------------------------"
echo "docs/SDD_FamilySystem.md -> docs/SDD_FamilyService.md"
echo "docs/SRS_FamilySystem.md -> docs/SRS_FamilyService.md"
echo "docs/STP_FamilySystem.md -> docs/STP_FamilyService.md"
echo

# Rename the files
mv docs/SDD_FamilySystem.md docs/SDD_FamilyService.md
mv docs/SRS_FamilySystem.md docs/SRS_FamilyService.md
mv docs/STP_FamilySystem.md docs/STP_FamilyService.md
echo "Files renamed successfully."
echo

# 2. References to these files in other files
echo "References being updated in Makefile:"
echo "------------------------------------"
grep -n "FamilySystem" Makefile
echo
echo "Making the following changes:"
echo "  Line 444: @echo \"- Software Requirements Specification (SRS): docs/SRS_FamilySystem.md\" -> @echo \"- Software Requirements Specification (SRS): docs/SRS_FamilyService.md\""
echo "  Line 445: @echo \"- Software Design Document (SDD): docs/SDD_FamilySystem.md\" -> @echo \"- Software Design Document (SDD): docs/SDD_FamilyService.md\""
echo "  Line 446: @echo \"- Software Test Plan (STP): docs/STP_FamilySystem.md\" -> @echo \"- Software Test Plan (STP): docs/STP_FamilyService.md\""
echo "  Line 447: @echo \"- Deployment Document: docs/Deployment_FamilySystem.md\" -> @echo \"- Deployment Document: docs/Deployment_FamilyService.md\""
echo

echo "References being updated in README.md:"
echo "-------------------------------------"
grep -n "FamilySystem" README.md
echo
echo "Making the following changes:"
echo "  Line 191: - [Software Requirements Specification (SRS)](./docs/SRS_FamilySystem.md) -> - [Software Requirements Specification (SRS)](./docs/SRS_FamilyService.md)"
echo "  Line 192: - [Software Design Document (SDD)](./docs/SDD_FamilySystem.md) -> - [Software Design Document (SDD)](./docs/SDD_FamilyService.md)"
echo "  Line 193: - [Software Test Plan (STP)](./docs/STP_FamilySystem.md) -> - [Software Test Plan (STP)](./docs/STP_FamilyService.md)"
echo "  Line 194: - [Deployment Document](./docs/Deployment_FamilySystem.md) -> - [Deployment Document](./docs/Deployment_FamilyService.md)"
echo

echo "References being updated in docs/Secrets_Setup_Guide.md:"
echo "-----------------------------------------------------"
grep -n "FamilySystem" docs/Secrets_Setup_Guide.md
echo
echo "Making the following changes:"
echo "  Line 73: For more information on how these services are configured, please refer to the [Deployment Document](Deployment_FamilySystem.md). -> For more information on how these services are configured, please refer to the [Deployment Document](Deployment_FamilyService.md)."
echo

# Update references in files
# Using -i '' for macOS compatibility
sed -i '' 's/SRS_FamilySystem.md/SRS_FamilyService.md/g' Makefile README.md
sed -i '' 's/SDD_FamilySystem.md/SDD_FamilyService.md/g' Makefile README.md
sed -i '' 's/STP_FamilySystem.md/STP_FamilyService.md/g' Makefile README.md
sed -i '' 's/Deployment_FamilySystem.md/Deployment_FamilyService.md/g' Makefile README.md docs/Secrets_Setup_Guide.md
echo "References updated successfully."
echo

echo "All changes have been successfully implemented:"
echo "1. Files have been renamed"
echo "2. References in Makefile, README.md, and docs/Secrets_Setup_Guide.md have been updated"
echo
echo "The migration from 'FamilySystem' to 'FamilyService' is now complete."
