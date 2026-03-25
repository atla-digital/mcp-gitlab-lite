## list_branches

List branches in a project repository.

### project_id
Project ID (numeric) or full path (e.g. "group/project").

### search
Filter branches by name (fuzzy match).

## get_branch

Get details of a single branch.

### project_id
Project ID (numeric) or full path.

### branch
Branch name.

## create_branch

Create a new branch from a ref.

### project_id
Project ID (numeric) or full path.

### branch
Name for the new branch.

### ref
Branch name, tag, or commit SHA to create the branch from.

## delete_branch

Delete a branch from the repository.

### project_id
Project ID (numeric) or full path.

### branch
Branch name to delete.

## list_repository_tree

List files and directories in a repository tree.

### project_id
Project ID (numeric) or full path.

### path
Path inside the repository to list (default: root).

### ref
Branch, tag, or commit to list the tree at.

### recursive
List tree recursively including subdirectories.

## get_file_content

Get the content of a file from the repository.

### project_id
Project ID (numeric) or full path.

### file_path
Path to the file within the repository.

### ref
Branch, tag, or commit to read the file at (default: default branch).

## list_commits

List commits in a project repository.

### project_id
Project ID (numeric) or full path.

### ref_name
Branch or tag name to list commits for.

### path
Filter commits affecting this file path.

### page
Page number for paginated results (default: 1).

## get_commit

Get details of a single commit.

### project_id
Project ID (numeric) or full path.

### sha
Commit SHA.

## compare_refs

Compare two branches, tags, or commits.

### project_id
Project ID (numeric) or full path.

### from
Source branch, tag, or commit SHA.

### to
Target branch, tag, or commit SHA.
