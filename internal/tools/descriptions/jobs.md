## list_pipeline_jobs

List all jobs in a pipeline.

### project_id
Project ID (numeric) or full path (e.g. "group/project").

### pipeline_id
The pipeline ID.

### scope
Filter by job status: "created", "pending", "running", "failed", "success", "canceled", "skipped", "manual".

## get_job

Get detailed information about a single job.

### project_id
Project ID (numeric) or full path.

### job_id
The job ID.

## get_job_log

Get the trace/log output of a job. Supports grep filtering and head/tail truncation.

### project_id
Project ID (numeric) or full path.

### job_id
The job ID.

### grep
Regex pattern to filter log lines.

### head
Return only the first N lines (or matched lines if grep is used).

### tail
Return only the last N lines (or matched lines if grep is used). Ignored if head is also set.

## retry_job

Retry a failed or canceled job.

### project_id
Project ID (numeric) or full path.

### job_id
The job ID.

## cancel_job

Cancel a running job.

### project_id
Project ID (numeric) or full path.

### job_id
The job ID.
