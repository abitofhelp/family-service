# Documentation Tools Scripts

## Overview

This directory contains scripts for managing and maintaining the documentation in the Family Service project.

## Scripts

### enforce_naming_convention.sh

This script enforces the naming convention for `.puml` and `.svg` files in the `DOCS/diagrams` directory. The naming convention is to use underscores and lowercase for all file names.

#### Usage

```bash
./enforce_naming_convention.sh
```

#### Features

- Checks all `.puml` and `.svg` files in the `DOCS/diagrams` directory
- Renames files that don't follow the naming convention (lowercase with underscores)
- Updates corresponding `.svg` files when a `.puml` file is renamed
- Provides feedback on renamed files

#### Integration with Makefile

The script is integrated with the project's Makefile and can be run using:

```bash
make enforce-naming-convention
```

It's also included in the pre-commit checks:

```bash
make pre-commit
```

## Best Practices

1. **Run Before Committing**: Always run `make enforce-naming-convention` before committing changes to ensure all diagram files follow the naming convention.
2. **Update Both Files**: When creating new diagrams, ensure both the `.puml` and `.svg` files follow the naming convention.
3. **Use Lowercase and Underscores**: All diagram filenames should be lowercase with underscores separating words.

## Related Components

- [DOCS/diagrams](../../diagrams/README.md) - The diagrams directory where the files are stored
- [DOCS/tools](../README.md) - Other documentation tools