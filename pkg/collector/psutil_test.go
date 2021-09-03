package collector

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"runtime"
	"testing"
	"time"
)

func TestGetMemery(t *testing.T) {
	v, _ := mem.VirtualMemory()
	fmt.Printf("Total: %v, Free:%.1f, UsedPercent:%f%%\n", v.Total/1e9, float64(v.Free)/1e9, v.UsedPercent)
	fmt.Println(v)
}

func TestGetDiskIO(t *testing.T) {
	name := "/mnt/sata"
	if runtime.GOOS == "windows" {
		name = "C:"
	}
	ret, err := disk.IOCounters(name)
	if err != nil {
		t.Errorf("error %v", err)
	}
	if len(ret) == 0 {
		t.Errorf("ret is empty")
	}
	empty := disk.IOCountersStat{}
	for part, io := range ret {
		fmt.Println(part, io)
		if io == empty {
			t.Errorf("io_counter error %v, %v", part, io)
		}
	}
}
func TestGetNetIO(t *testing.T) {
	var i int
	for i < 50 {
		ret, _ := net.IOCounters(true)
		for _, io := range ret {
			if io.Name == "以太网" {
				fmt.Println("接收流量", float64(io.BytesRecv)/1e6)
				fmt.Println("发送流量", float64(io.BytesSent)/1e6)
				fmt.Println("发送包", float64(io.PacketsSent)/1e6)
				fmt.Println("接收包", float64(io.PacketsRecv)/1e6)
			}
		}
		time.Sleep(time.Second)
		i++
	}

}

func TestGetLoad(t *testing.T) {
	var i int
	for i < 50 {
		info, _ := load.Avg()
		fmt.Printf("load avg %+v", info)
		time.Sleep(time.Second)
	}

}

func TestAppend(t *testing.T) {
	var ints []int = make([]int, 0)
	ints = append(ints, 1)
	ints = append(ints, 2)
	ints = append(ints, 3)
	ints = append(ints, 4)
	fmt.Println(ints)
	fmt.Println(cap(ints))
	fmt.Println(len(ints))
}
