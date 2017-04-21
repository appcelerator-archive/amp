package labels

//Amp labels name
const (
	// LabelsNameRole describe the role of the container/service in amp context
	LabelsNameRole = "io.amp.role"

	// LabelsMapping set a maaping between a logical host name and a service port
	// format "name=myName account=myAccount port=myPort" -> will create a reverse proxy on http://myName.stackName.myAccount.local.appcelerator.io -> service:port
	// by default name is the service name
	LabelsNameMapping = "io.amp.mapping"
)

//Amp label values
const (
	// LabelsValuesRoleInfrastructure amp role for infrastructure container
	LabelsValuesRoleInfrastructure = "infrastructure"

	// LabelsValuesRoleTools amp role for tools containers
	LabelsValuesRoleTools = "tools"
)
