# Conventional Commits Reference Guide

We use the [Conventional Commits Specification](https://www.conventionalcommits.org/en/v1.0.0/) to maintain a clean, machine-readable Git history. This format allows us to automatically update versions (SemVer) and generate changelogs during releases.

## Commit Message Structure

Every commit message must follow this structure:

```text
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

---

## Commit Types

### Core Types

- **`feat`**: A new feature for the user, not a new feature for build script. (Triggers a **MINOR** version bump)
- **`fix`**: A bug fix for the user. (Triggers a **PATCH** version bump)

### Other Common Types

- **`docs`**: Documentation-only changes (e.g., editing `README.md`).
- **`style`**: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc.).
- **`refactor`**: A code change that neither fixes a bug nor adds a feature.
- **`perf`**: A code change that improves performance.
- **`test`**: Adding missing tests or correcting existing tests.
- **`build`**: Changes that affect the build system or external dependencies (e.g., npm, glp).
- **`ci`**: Changes to CI configuration files and scripts (e.g., GitHub Actions, Travis).
- **`chore`**: Other changes that don't modify src or test files (e.g., updating `.gitignore`).
- **`revert`**: Reverts a previous commit.

---

## Formatting Rules

1. **Imperative Mood**: Write the description in the imperative mood ("add" instead of "added" or "adds").
2. **No Capitalization**: Do not capitalize the first letter of the description.
3. **No Period**: Do not place a period (`.`) at the end of the description line.
4. **Scope**: Wrap the scope in parentheses to provide contextual location (e.g., `feat(auth):`).

---

## Breaking Changes

A breaking change triggers a **MAJOR** version bump. It must be indicated by an exclamation mark (`!`) before the colon, or by a `BREAKING CHANGE:` footer.

### Example with Exclamation Mark

```text
refactor(api)!: drop support for Node v18
```

### Example with Footer

```text
feat(auth): switch to OAuth2 protocol

BREAKING CHANGE: The old token endpoint is no longer active.
```

---

## Full Examples

### Simple Commit

```text
fix: resolve memory leak in background workers
```

### Commit with Scope

```text
feat(ui): add dark mode toggle switch
```

### Commit with Body and Footer

```text
fix(database): enforce strict unique constraint on emails

The previous validation checked rows via application logic. This change moves
the safety constraint directly to the PostgreSQL database layer.

Closes #142
```
