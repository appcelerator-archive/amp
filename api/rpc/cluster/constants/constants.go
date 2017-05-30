package constants

// cluster constant
const (
	LabelKeyRole = "io.amp.role"

	SecretAmplifier   = "amplifier_yml"
	SecretCertificate = "certificate_atomiq"
)

// cluster vars
var (
	Labels  = []string{LabelKeyRole}
	Secrets = []string{SecretAmplifier, SecretCertificate}
)
