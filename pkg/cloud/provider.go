package cloud

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/appcelerator/amp/pkg/cloud/aws"
)

// amplifier has to detect which cloud provider it's deployed on
type Provider string

const (
	ProviderUnknown = Provider("unknown")
	ProviderLocal   = Provider("local")
	ProviderAWS     = Provider("AWS")
	ProviderAzure   = Provider("Azure")
	ProviderDO      = Provider("DO")
	ProviderGCP     = Provider("GCP")
)

// CloudProvider returns the cloud provider
func CloudProvider() (Provider, error) {
	dataLen := 3
	uuidFile, err := os.Open("/sys/hypervisor/uuid")
	if err != nil {
		uuidFile, err = os.Open("/sys/devices/virtual/dmi/id/product_uuid")
	}
	if err != nil {
		// file does not exist, so we'll consider it's not a cloud deployment
		return ProviderLocal, nil
	}
	data := make([]byte, dataLen)
	count, err := uuidFile.Read(data)
	if err != nil {
		return ProviderUnknown, fmt.Errorf("Unable to establish provider from uuid file: %s", err)
	}
	if count != dataLen {
		log.Infoln("Unable to establish provider, empty uuid file")
		return ProviderLocal, nil
	}
	switch string(data) {
	case "ec2", "EC2":
		return ProviderAWS, nil
	default:
		return ProviderLocal, nil
	}
}

func Region() (string, error) {
	p, err := CloudProvider()
	if err != nil {
		return "", err
	}
	switch p {
	case ProviderLocal:
		return "", nil
	case ProviderAWS:
		return aws.Region()
	default:
		return "", fmt.Errorf("provider not implemented: [%s]", p)
	}

	return "", nil
}
