# Configurable Timing Intervals

## Overview
Both the workflow loop interval and job check interval have been made configurable to allow flexibility in how often the main processes run and check job statuses.

## Changes Made

### 1. Configuration Structure
Added two new fields to the `Update_Groups` struct in `pkg_updates/structs.go`:

```go
WorkflowLoopInterval int `json:"workflow_loop_interval"`  // For Start_Workflow function
JobCheckInterval     int `json:"job_check_interval"`      // For Check_Jobs function
```

### 2. Implementation

#### Workflow Loop Interval
Modified the `Start_Workflow` function in `pkg_updates/start.go`:
- If `workflow_loop_interval` is set to a positive value in the YAML config, it uses that value
- If not set or set to 0 or negative, it defaults to 20 seconds (maintaining backward compatibility)

#### Job Check Interval  
Modified the `Check_Jobs` function in `pkg_updates/check_jobs.go`:
- If `job_check_interval` is set to a positive value in the YAML config, it uses that value
- If not set or set to 0 or negative, it defaults to 60 seconds (maintaining backward compatibility)

### 3. Usage in YAML Configuration

Add the following parameters to your YAML configuration file:

```yaml
# Controls how long the main workflow loop waits between iterations (in seconds)
# Default: 20 seconds if not specified or set to 0/negative value
workflow_loop_interval: 30

# Controls how long between job status checks (in seconds) 
# Default: 60 seconds if not specified or set to 0/negative value
job_check_interval: 45
```

### 4. Examples

**Fast processing:**
```yaml
workflow_loop_interval: 10  # Check workflow every 10 seconds
job_check_interval: 30      # Check jobs every 30 seconds
```

**Standard processing (defaults):**
```yaml
workflow_loop_interval: 20  # OR omit for default
job_check_interval: 60      # OR omit for default
```

**Slower processing:**
```yaml
workflow_loop_interval: 60  # Check workflow every minute
job_check_interval: 120     # Check jobs every 2 minutes
```

**Mixed scenarios:**
```yaml
# Fast workflow checks but slower job checks
workflow_loop_interval: 15
job_check_interval: 90

# Slow workflow checks but fast job checks  
workflow_loop_interval: 45
job_check_interval: 30
```

### 5. Benefits

- **Flexibility**: Adjust both workflow and job check frequencies independently based on system needs
- **Performance Tuning**: Reduce system load by increasing intervals or speed up processing by decreasing intervals
- **Backward Compatibility**: Existing configurations without these parameters continue to work with defaults (20s workflow, 60s job checks)
- **Environment-Specific**: Different environments can have different processing speeds (e.g., test vs production)
- **Independent Control**: Workflow execution and job monitoring can be tuned separately

### 6. Considerations

#### Workflow Loop Interval
- **Minimum Recommended**: 5 seconds (to avoid excessive system load)
- **Maximum Practical**: 300 seconds (5 minutes - to ensure timely processing)
- **Production Recommendation**: 20-60 seconds depending on system capacity

#### Job Check Interval
- **Minimum Recommended**: 10 seconds (jobs need time to execute)
- **Maximum Practical**: 600 seconds (10 minutes - to ensure timely status updates)
- **Production Recommendation**: 30-120 seconds depending on job complexity and urgency

#### General Guidelines
- Set job check interval higher than workflow interval if jobs are long-running
- Use shorter intervals in test environments for faster feedback
- Consider system resources when setting very short intervals
- Monitor system performance and adjust accordingly
