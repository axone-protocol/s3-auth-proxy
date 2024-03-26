package dataverse

import "strings"

type GovernanceExecAnswer struct {
	Result   string
	Evidence string
}

type ExecutionOrderContext struct {
	Zone      string
	Statuses  []string
	Resources []string
}

func (e *ExecutionOrderContext) IsInProgress() bool {
	inProgress := false
	for _, status := range e.Statuses {
		if strings.HasSuffix(status, "Cancelled") || strings.HasSuffix(status, "Delivered") || strings.HasSuffix(status, "Failed") {
			return false
		}
		inProgress = strings.HasSuffix(status, "InExecution")
	}

	return inProgress
}

func (e *ExecutionOrderContext) HasResource(resource string) bool {
	for _, r := range e.Resources {
		if r == resource {
			return true
		}
	}
	return false
}
