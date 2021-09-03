//操作系统
package collector

import (
	"encoding/json"
	"time"

	"github.com/shirou/gopsutil/v3/load"
	"github.com/xunlbz/go-restful-template/pkg/lib"
)

var (
	label_load    = "load"
	name_load_avg = "avg"
	desc_load_avg = "系统负载"
)
var list []map[string]interface{}
var size int = 200

type loadCollector struct {
	metrics Metrics
}

func init() {
	registerCollector(label_load, NewLoadCollector)
}

func NewLoadCollector() (Collector, error) {
	lc := new(loadCollector)
	list = make([]map[string]interface{}, 0)
	go doSystemLoad()
	return lc, nil
}

func (c *loadCollector) Update() Metrics {
	c.metrics = NewMetrics(label_load, make([]Metric, 0))
	c.metrics.Entries = append(c.metrics.Entries, c.getSystemLoad())
	return c.metrics
}

func (c *loadCollector) Read() Metrics {
	c.Update()
	return c.metrics
}

func (c *loadCollector) getSystemLoad() (m Metric) {
	val, _ := json.Marshal(list)
	return NewMetric(label_load, name_load_avg, string(val), desc_load_avg)
}

func doSystemLoad() {
	for {
		if time.Now().Second()%15 == 0 {
			stat, err := load.Avg()
			if err != nil {
				return
			}
			var item map[string]interface{} = make(map[string]interface{})
			item[lib.GetNowTimeString()] = stat
			list = append(list, item)
			if len(list) > size {
				list = list[1:]
			}
		}
		time.Sleep(time.Second * 1)
	}
}
