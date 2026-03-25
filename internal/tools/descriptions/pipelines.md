## list_pipelines

List pipelines in a project.

### project_id
Project ID (numeric) or full path (e.g. "group/project").

### status
Filter by status: "running", "pending", "success", "failed", "canceled", "skipped", "manual".

### ref
Filter by git ref (branch or tag name).

### page
Page number for paginated results (default: 1).

## get_pipeline

Get detailed information about a single pipeline.

### project_id
Project ID (numeric) or full path.

### pipeline_id
The pipeline ID.

## create_pipeline

Trigger a new pipeline on a ref.

### project_id
Project ID (numeric) or full path.

### ref
Branch or tag name to run the pipeline on.

### variables
JSON object of key-value pairs for pipeline variables (e.g. {"KEY":"value"}).

## cancel_pipeline

Cancel a running pipeline.

### project_id
Project ID (numeric) or full path.

### pipeline_id
The pipeline ID.

## retry_pipeline

Retry all failed jobs in a pipeline.

### project_id
Project ID (numeric) or full path.

### pipeline_id
The pipeline ID.

## wait_pipeline_change

Poll a pipeline until its status changes from the current value, or timeout.

### project_id
Project ID (numeric) or full path.

### pipeline_id
The pipeline ID.

### timeout_seconds
Maximum seconds to wait (default: 300).

### poll_interval_seconds
Seconds between polls (default: 5, minimum: 3).

## wait_pipeline_finish

Poll a pipeline until it reaches a terminal state (success, failed, canceled, skipped, manual), or timeout.

### project_id
Project ID (numeric) or full path.

### pipeline_id
The pipeline ID.

### timeout_seconds
Maximum seconds to wait (default: 1800).
