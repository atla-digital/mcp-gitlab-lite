## v1.0.0 (2026-04-09)

### BREAKING CHANGE

- the MR tools now expect \`merge_request_iid\` instead of
\`mr_iid\`, and \`list_projects\` expects a boolean \`membership\` instead
of a string.

### Fix

- **tools**: align parameter names with GitLab API and fix dead params

## v0.2.1 (2026-04-09)

### Fix

- **args**: accept numeric strings for int parameters (#2)

## v0.2.0 (2026-03-25)

### Feat

- add CI/CD pipelines

### Fix

- move max-issues settings to correct config section

## v0.1.0 (2026-03-25)

### Feat

- initial implementation of gitlab-mcp-lite
