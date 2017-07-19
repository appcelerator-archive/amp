package constants

// cluster constant
const (
	LabelKeyRole = "io.amp.role"

	SecretAmplifier   = "amplifier_yml"
	SecretCertificate = "certificate_amp"
)

// cluster vars
var (
	Labels  = []string{LabelKeyRole}
	Secrets = []string{SecretAmplifier, SecretCertificate}
)
