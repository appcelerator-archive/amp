package service

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	publishSpecRegex = `((((?P<name>[a-zA-Z0-9]+([a-zA-Z0-9-]*[a-zA-Z0-9])*):)?(?P<publish_port>[\d]{1,5}):)?)(?P<internal_port>[\d]{1,5})(/(?P<protocol>(tcp|udp)))?([^:]$|$)`
)

var (
	publishSpecParser *regexp.Regexp
)

// ParsePublishSpec parses a string and returns a PublishSpec
func ParsePublishSpec(s string) (publishSpec PublishSpec, err error) {
	if publishSpecParser == nil {
		publishSpecParser = regexp.MustCompile(publishSpecRegex)
	}

	m := publishSpecParser.FindStringSubmatch(s)
	if m == nil {
		err = fmt.Errorf("\"%s\" is not a valid PublishSpec", s)
		return
	}

	names := publishSpecParser.SubexpNames()
	nameMap := mapNames(names)
	for name, index := range nameMap {
		val := m[index]
		switch name {
		case "name":
			publishSpec.Name = val
		case "publish_port":
			err = portAtoi(val, &publishSpec.PublishPort)
			if err != nil {
				return
			}
		case "internal_port":
			err = portAtoi(val, &publishSpec.InternalPort)
			if err != nil {
				return
			}
		case "protocol":
			publishSpec.Protocol = val
		}
	}

	return
}

// ParseNetwork parses a string and returns a NetworkAttachement
func ParseNetwork(s string) *NetworkAttachment {
	list := strings.Split(s, ":")
	network := NetworkAttachment{
		Target: list[0],
	}
	if len(list) > 1 {
		network.Aliases = list[1:]
	}
	return &network
}

func portAtoi(s string, port *uint32) error {
	var u64 uint64
	if s == "" {
		s = "0"
	}
	u64, err = strconv.ParseUint(s, 10, 32)
	if err != nil {
		return err
	}
	*port = uint32(u64)
	return nil
}

func mapNames(names []string) map[string]int {
	nameMap := make(map[string]int)
	for i, name := range names {
		if name != "" {
			nameMap[name] = i
		}
	}
	return nameMap
}
