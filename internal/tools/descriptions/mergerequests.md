## list_merge_requests

List merge requests in a project. Returns compact references.

### project_id
Project ID (numeric) or full path (e.g. "group/project").

### state
Filter by state: "opened", "closed", "merged", or "all". Default: "opened".

### target_branch
Filter by target branch name.

### source_branch
Filter by source branch name.

### search
Search MRs by title or description.

### page
Page number for paginated results (default: 1).

## get_merge_request

Get full details of a single merge request.

### project_id
Project ID (numeric) or full path.

### mr_iid
The internal ID of the merge request within the project.

## create_merge_request

Create a new merge request.

### project_id
Project ID (numeric) or full path.

### title
Title for the merge request.

### source_branch
Branch containing the changes.

### target_branch
Branch to merge into.

### description
Markdown description body.

### remove_source_branch
Delete source branch after merge.

### squash
Squash commits when merging.

## update_merge_request

Update an existing merge request.

### project_id
Project ID (numeric) or full path.

### mr_iid
The internal ID of the merge request.

### title
New title (leave empty to keep current).

### description
New description (leave empty to keep current).

### target_branch
New target branch.

### state_event
State transition: "close" or "reopen".

### labels
Comma-separated label names (replaces existing labels).

## merge_merge_request

Accept and merge a merge request.

### project_id
Project ID (numeric) or full path.

### mr_iid
The internal ID of the merge request.

### merge_commit_message
Custom merge commit message.

### squash
Squash commits when merging.

### should_remove_source_branch
Delete source branch after merge.

## approve_merge_request

Approve a merge request.

### project_id
Project ID (numeric) or full path.

### mr_iid
The internal ID of the merge request.

## list_mr_notes

List all comments on a merge request.

### project_id
Project ID (numeric) or full path.

### mr_iid
The internal ID of the merge request.

## create_mr_note

Add a comment to a merge request.

### project_id
Project ID (numeric) or full path.

### mr_iid
The internal ID of the merge request.

### body
Markdown body of the comment.

## create_mr_discussion

Create a new discussion thread on a merge request diff.

### project_id
Project ID (numeric) or full path.

### mr_iid
The internal ID of the merge request.

### body
Markdown body of the discussion comment.

### base_sha
Base commit SHA of the merge request diff.

### head_sha
Head commit SHA of the merge request diff.

### start_sha
Start commit SHA of the merge request diff.

### new_path
File path in the new version (for new/modified files).

### old_path
File path in the old version (for renamed/deleted files).

### new_line
Line number in the new version to comment on.

### old_line
Line number in the old version to comment on.

## list_mr_diffs

List all file diffs in a merge request.

### project_id
Project ID (numeric) or full path.

### mr_iid
The internal ID of the merge request.

## list_mr_commits

List all commits in a merge request.

### project_id
Project ID (numeric) or full path.

### mr_iid
The internal ID of the merge request.

## list_mr_pipelines

List all pipelines for a merge request.

### project_id
Project ID (numeric) or full path.

### mr_iid
The internal ID of the merge request.

## search_merge_requests

Search merge requests across all projects the user has access to.

### search
Search query string.

### state
Filter by state: "opened", "closed", "merged", or "all". Default: "opened".

### page
Page number for paginated results (default: 1).
