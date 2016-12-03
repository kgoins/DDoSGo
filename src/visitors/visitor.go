package visitors

import "subsystems"

var (
	alertVisitor *AlertVisitor
)

type AlertVisitor struct {
	alertingSys *subsystems.AlertSystem
}

func NewAlertSystemVisitor(alertingSys *subsystems.AlertSystem) *AlertVisitor {
	return &AlertVisitor{alertingSys: alertingSys}
}

func (alertVisitor *AlertVisitor) HeldInstance() *subsystems.AlertSystem {
	return alertVisitor.alertingSys
}
