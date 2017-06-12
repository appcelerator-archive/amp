package pstats

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//PrometheusSource definition
type PrometheusSource struct {
	name      string
	url       string
	port      string
	metricMap map[string]*PrometheusMetric
}

func newSource(name string, url string, port string) *PrometheusSource {
	return &PrometheusSource{
		name:      name,
		url:       url,
		port:      port,
		metricMap: make(map[string]*PrometheusMetric),
	}
}

func (s *PrometheusSource) load() {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("http://%s:%s/metrics", s.url, s.port), nil)
	req.Header.Add("Accept", "text/plain")
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Error loading source %s: %v", s.name, err)
		return
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	lines := strings.Split(string(respBody), "\n")
	mtype := ""
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			headers := strings.Split(line, " ")
			if len(headers) >= 4 {
				mtype = headers[3]
			}
		} else {
			param := ""
			items := strings.Split(line, " ")[0]
			params := strings.Split(items, "{")
			_, ok := s.metricMap[params[0]]
			if !ok {
				if len(params) > 1 {
					param = strings.Split(params[1], "}")[0]
				}
				s.addMetric(params[0], mtype, param)
			}
		}
	}
}

func (s *PrometheusSource) addMetric(name string, mtype string, param string) {
	params := make([]string, 0)
	list := strings.Split(param, ",")
	for _, item := range list {
		params = append(params, strings.Split(item, "=")[0])
	}
	s.metricMap[name] = newMetric(name, mtype, params)
}
