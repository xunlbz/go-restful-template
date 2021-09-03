package collector

import (
	"encoding/json"

	"github.com/xunlbz/go-restful-template/pkg/lib"
	"github.com/xunlbz/go-restful-template/pkg/log"
)

var (
	label_service     = "service"
	name_service_list = "list"
	desc_service_list = "系统服务列表"
)

type serviceCollector struct {
	metrics Metrics
}

func init() {
	registerCollector(label_service, NewServiceCollector)
}

func NewServiceCollector() (Collector, error) {
	return new(serviceCollector), nil
}

func (c *serviceCollector) Update() Metrics {
	c.metrics = NewMetrics(label_service, make([]Metric, 0))
	c.metrics.Entries = append(c.metrics.Entries, c.getServices())
	return c.metrics
}

func (c *serviceCollector) Read() Metrics {
	return c.metrics
}

func (c *serviceCollector) getServices() Metric {
	var metric Metric
	var ss []Service = make([]Service, 0)
	units, err := lib.ListUnit([]string{"wyzl*", "docker.service"})
	if err != nil {
		log.Errorf("get units error:%v", err)
		return metric
	}
	log.Infof("services %v", units)
	for _, unit := range units {
		s := Service{
			Name:        unit.Name,
			Description: unit.Description,
			Sub:         unit.SubState,
			Load:        unit.LoadState,
			Active:      unit.ActiveState,
		}
		ss = append(ss, s)
	}
	data, err := json.Marshal(ss)
	if err != nil {
		return metric
	}
	return NewMetric(label_service, name_service_list, string(data), desc_service_list)
}

type Service struct {
	Name        string `json:"name"`
	Load        string `json:"load"`
	Active      string `json:"active"`
	Sub         string `json:"sub"`
	Description string `json:"desc"`
}

func (s Service) String() string {
	b, _ := json.Marshal(s)
	return string(b)
}
