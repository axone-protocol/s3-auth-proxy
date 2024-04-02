package dataverse

import "strings"

type GovernanceExecAnswer struct {
	Result   string
	Evidence string
}

type ExecutionOrderContext struct {
	Zone string

	/// Executions contains for each execution id related to the order, the linked statuses.
	Executions map[string][]string
}

func (e *ExecutionOrderContext) ExecutionsInProgress() []string {
	var results []string
	for exec, statuses := range e.Executions {
		inProgress := false
		for _, status := range statuses {
			if strings.HasSuffix(status, "Cancelled") || strings.HasSuffix(status, "Delivered") || strings.HasSuffix(status, "Failed") {
				inProgress = false
				break
			}
			inProgress = strings.HasSuffix(status, "InExecution")
		}

		if inProgress {
			results = append(results, exec)
		}
	}

	return results
}
