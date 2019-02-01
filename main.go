package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	mp "github.com/mackerelio/go-mackerel-plugin"
	statsHTTP "github.com/mercari/go-httpstats"
)

var (
	Version string
)

type ChoconPlugin struct {
	Target string
	Prefix string
}

var (
	DefaultHost        = "127.0.0.1"
	DefaultPrefix      = "chocon"
	DefaultPort        = "80"
	DefaultTempFile    = "/tmp/mackerel-plugin-chocon.json"
	MetricsHTTPRequest = []mp.Metrics{
		{
			Label: "http count in a period",
			Name:  "count",
			Diff:  true,
		},
	}
	MetricsHTTPResponseTime = []mp.Metrics{
		{
			Label: "average http response time",
			Name:  "avg_time",
		},
	}
	MetricsHTTPRequestPerStatus       = buildMetricsHTTPRequestPerStatus()
	MetricsHTTPResponsePercetiledTime = buildMetricsHTTPResponsePercentiledTime()

	graphdef = map[string]mp.Graphs{
		"http.requests": mp.Graphs{
			Label:   "http requests",
			Unit:    mp.UnitInteger,
			Metrics: append(MetricsHTTPRequest, MetricsHTTPRequestPerStatus...),
		},
		"http.latency": mp.Graphs{
			Label:   "http latency",
			Unit:    mp.UnitFloat,
			Metrics: append(MetricsHTTPResponseTime, MetricsHTTPResponsePercetiledTime...),
		},
	}
	percents  = []int{90, 95, 99}
	status4xx = []int{
		http.StatusNotFound,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusMethodNotAllowed,
	}
	status5xx = []int{
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
	}
)

func buildMetricsHTTPRequestPerStatus() []mp.Metrics {
	var (
		statuses = []string{
			"200",
			"403",
			"4xx",
			"5xx",
		}
		metrics = make([]mp.Metrics, len(statuses))
	)
	for i, s := range statuses {
		metrics[i] = mp.Metrics{
			Label: fmt.Sprintf("http request of status %s", s),
			Name:  fmt.Sprintf("count_%s", s),
			Diff:  true,
		}
	}
	return append(metrics)
}

func buildMetricsHTTPResponsePercentiledTime() []mp.Metrics {
	var (
		metricsPercentile = make([]mp.Metrics, len(percents))
	)
	for i, p := range percents {
		metricsPercentile[i] = mp.Metrics{
			Label: fmt.Sprintf("%d percentile of http response time", p),
			Name:  fmt.Sprintf("time_percentile_%d", p),
		}
	}
	return append(metricsPercentile)
}

func (p ChoconPlugin) FetchMetrics() (map[string]float64, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/.api/http-stats", p.Target))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var chocon statsHTTP.Data
	if err := json.NewDecoder(resp.Body).Decode(&chocon); err != nil {
		return nil, err
	}
	stats := map[string]float64{
		"count":    float64(chocon.Request.Count),
		"avg_time": chocon.Response.AverageTime,
	}
	var count4xx int64
	for _, status := range status4xx {
		count4xx += chocon.Request.StatusCount[status]
	}
	var count5xx int64
	for _, status := range status5xx {
		count5xx += chocon.Request.StatusCount[status]
	}
	stats["count_200"] = float64(chocon.Request.StatusCount[200])
	stats["count_403"] = float64(chocon.Request.StatusCount[403])
	stats["count_4xx"] = float64(count4xx)
	stats["count_5xx"] = float64(count5xx)

	for _, p := range percents {
		stats[fmt.Sprintf("time_percentile_%d", p)] = chocon.Response.PercentiledTime[p]
	}

	return stats, nil
}

func (p ChoconPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

func (p ChoconPlugin) MetricKeyPrefix() string {
	return p.Prefix
}

func main() {
	var (
		optHost     = flag.String("host", DefaultHost, "hostname")
		optPort     = flag.String("port", DefaultPort, "port")
		optPrefix   = flag.String("prefix", DefaultPrefix, "metrics name prefix")
		optTempfile = flag.String("tempfile", DefaultTempFile, "temporary file")
		optVersion  = flag.Bool("version", false, "print version")
	)
	flag.Parse()
	log.SetOutput(os.Stderr)
	if *optVersion {
		fmt.Printf("version: %s\n", Version)
		return
	}
	if _, err := os.Stat(*optTempfile); err != nil {
		if err := os.MkdirAll(filepath.Dir(*optTempfile), 0666); err != nil {
			log.Println(err)
		}
		f, err := os.Create(*optTempfile)
		if err != nil {
			log.Println(err)
		}
		f.Close()
	}

	helper := mp.NewMackerelPlugin(ChoconPlugin{
		Target: fmt.Sprintf("%s:%s", *optHost, *optPort),
		Prefix: *optPrefix,
	})
	helper.Tempfile = *optTempfile
	helper.Run()
}
