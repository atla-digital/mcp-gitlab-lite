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

## list_issue_discussions

List all discussion threads on an issue. Each thread has an `id` and ordered `notes` (first is the root, the rest are replies). Use the thread `id` with reply_to_issue_discussion.

### project_id
Project ID (numeric) or full path.

### issue_iid
The internal ID of the issue.

## reply_to_issue_discussion

Reply inside an existing discussion thread on an issue.

### project_id
Project ID (numeric) or full path.

### issue_iid
The internal ID of the issue.

### discussion_id
ID of the thread to reply to (from list_issue_discussions).

### body
Markdown body of the reply.

## create_issue_note

Add a new top-level comment to an issue (starts a new thread; to answer an existing thread use reply_to_issue_discussion).

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
