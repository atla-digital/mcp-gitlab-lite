## list_issues

List issues in a project. Returns compact references (no description body).

### project_id
Project ID (numeric) or full path (e.g. "group/project").

### state
Filter by state: "opened", "closed", or "all". Default: "opened".

### labels
Comma-separated label names to filter by.

### assignee_username
Filter by assignee username.

### search
Search issues by title or description.

### page
Page number for paginated results (default: 1).

## get_issue

Get full details of a single issue, including description.

### project_id
Project ID (numeric) or full path.

### issue_iid
The internal ID of the issue within the project.

## create_issue

Create a new issue in a project.

### project_id
Project ID (numeric) or full path.

### title
Title for the new issue.

### description
Markdown description body.

### labels
Comma-separated label names to apply.

## update_issue

Update an existing issue's title, description, state, or labels.

### project_id
Project ID (numeric) or full path.

### issue_iid
The internal ID of the issue within the project.

### title
New title (leave empty to keep current).

### description
New description (leave empty to keep current).

### state_event
State transition: "close" or "reopen".

### labels
Comma-separated label names (replaces existing labels).

## list_issue_notes

List all comments on an issue.

### project_id
Project ID (numeric) or full path.

### issue_iid
The internal ID of the issue.

## create_issue_note

Add a comment to an issue.

### project_id
Project ID (numeric) or full path.

### issue_iid
The internal ID of the issue.

### body
Markdown body of the comment.

## search_issues

Search issues across all projects the user has access to.

### search
Search query string.

### state
Filter by state: "opened", "closed", or "all". Default: "opened".

### page
Page number for paginated results (default: 1).
