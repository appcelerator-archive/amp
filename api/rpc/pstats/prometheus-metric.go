package pstats

import "log"

//PrometheusMetric definition
type PrometheusMetric struct {
	name   string
	mtype  string
	labels []string
}

func newMetric(name string, mtype string, labels []string) *PrometheusMetric {
	log.Printf("add metric name=%s type=%s labels=%v\n", name, mtype, labels)
	return &PrometheusMetric{
		name:   name,
		mtype:  mtype,
		labels: labels,
	}
}
