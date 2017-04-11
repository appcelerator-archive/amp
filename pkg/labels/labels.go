package labels

//Amp labels name
const (
	// LabelsNameRole describe the role of the container/service in amp context
	LabelsNameRole = "io.amp.role"
)

//Amp label values
const (
	// LabelsValuesRoleInfrastructure amp role for infrastructure container
	LabelsValuesRoleInfrastructure = "infrastructure"

	// LabelsValuesRoleTools amp role for tools containers
	LabelsValuesRoleTools = "tools"
)
