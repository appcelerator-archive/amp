package service

// import (
//         "fmt"
// 	"regexp"
// )
//
// const (
// 	portSpecRegex = `(?P<published_port>[\d]*):?(?P<target_port>[\d]+)(/(?P<protocol>tcp|udp))?`
// )
//
// var (
// 	portSpecParser *regexp.Regexp
// )
//
// // ParsePortSpec parses a string and returns a PortSpec
// func ParsePortSpec(s string) (portSpec PortSpec, err error) {
// 	if portSpecParser == nil {
// 		portSpecParser = regexp.MustCompile(portSpecRegex)
// 	}
//
//         m := portSpecParser.FindStringSubmatch(s)
//         fmt.Println(portSpecParser.SubexpNames())
//         fmt.Println(m)
//
//         if m == nil {
//                 err = fmt.Errorf("\"%s\" is not a valid PortSpec", s)
//                 return
//         }
//
//
//
//         return
// }
