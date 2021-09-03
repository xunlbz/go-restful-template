package collector

import (
	"fmt"
	"sync"
	"time"

	"github.com/xunlbz/go-restful-template/pkg/log"
)

var factories = make(map[string]func() (Collector, error))
var collectors map[string]Collector = make(map[string]Collector)

type Metric struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Desc  string `json:"desc"`
}

type Metrics struct {
	Label     string   `json:"label"`
	Entries   []Metric `json:"entries"`
	Timestamp int64    `json:"timestamp"`
}

type Collector interface {
	Update() Metrics
	Read() Metrics
}

type EdgeCollector struct {
}

func registerCollector(label string, f func() (Collector, error)) {
	if _, ok := factories[label]; !ok {
		factories[label] = f
		return
	}
	log.Errorf("register collertor %s is already exits", label)
}

func Register() {
	for label, f := range factories {
		log.Debugf("register collertor %s", label)
		c, err := f()
		if err != nil {
			log.Errorf("create collector %s error %v", label, err)
		}
		collectors[label] = c
	}
}

func NewMetric(label, name, value, desc string) Metric {
	name = fmt.Sprintf("%s_%s", label, name)
	return Metric{Name: name, Value: value, Desc: desc}
}

func NewMetrics(label string, entries []Metric) Metrics {
	metrics := Metrics{}
	metrics.Label = label
	metrics.Entries = entries
	metrics.Timestamp = time.Now().Unix()
	return metrics
}

func NewEdgeCollector() EdgeCollector {
	ec := EdgeCollector{}
	return ec
}

func (c EdgeCollector) Collect() {
	t1 := time.Now()
	var sw sync.WaitGroup
	for _, c := range collectors {
		sw.Add(1)
		go func(c Collector) {
			defer sw.Done()
			c.Update()
		}(c)
	}
	sw.Wait()
	log.Infof("edge collector collect end takes %v", time.Since(t1))
}

func (c EdgeCollector) GetModuleMitrics(label string) Metrics {
	if c1, ok := collectors[label]; ok {
		return c1.Read()
	}
	return Metrics{}
}
