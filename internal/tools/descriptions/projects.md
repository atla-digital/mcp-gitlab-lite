## list_projects

List GitLab projects. By default only returns projects the authenticated user is a member of.

### search
Filter projects by name or path (fuzzy match).

### membership
If true, only return projects the authenticated user is a member of. If false, return all visible projects. Default: true.

### visibility
Filter by visibility: "public", "internal", or "private".

### page
Page number for paginated results (default: 1).

## get_project

Get detailed information about a single project.

### project_id
Project ID (numeric) or full path (e.g. "group/project").

## list_group_projects

List projects within a GitLab group.

### group_id
Group ID (numeric) or full path.

### search
Filter projects by name or path.

### include_subgroups
Include projects from subgroups.

### page
Page number for paginated results (default: 1).
