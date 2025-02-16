package bpf

import (
	"bytes"
	"embed"
	"fmt"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

//go:embed cpu.bpf.o mem.bpf.o
var bpfPrograms embed.FS

func LoadCpuProgram() (*ebpf.Collection, error) {
	file, err := bpfPrograms.ReadFile("cpu.bpf.o")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded eBPF program: %v", err)
	}

	spec, err := ebpf.LoadCollectionSpecFromReader(bytes.NewReader(file))
	if err != nil {
		return nil, fmt.Errorf("failed to load collection spec: %v", err)
	}

	coll, err := ebpf.NewCollection(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to create collection: %v", err)
	}

	return coll, nil
}

func LoadMemProgram() (*ebpf.Collection, error) {
	file, err := bpfPrograms.ReadFile("mem.bpf.o")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded eBPF program: %v", err)
	}

	spec, err := ebpf.LoadCollectionSpecFromReader(bytes.NewReader(file))
	if err != nil {
		return nil, fmt.Errorf("failed to load collection spec: %v", err)
	}

	coll, err := ebpf.NewCollection(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to create collection: %v", err)
	}

	return coll, nil
}

func MonitorCPU(resourceType string) {
	if err := rlimit.RemoveMemlock(); err != nil {
		fmt.Printf("Failed to remove memlock limit: %v\n", err)
		return
	}

	coll, err := LoadCpuProgram()
	if err != nil {
		fmt.Printf("Failed to load eBPF program: %v\n", err)
		return
	}
	defer coll.Close()

	prog := coll.Programs["cpu_monitor"]
	if prog == nil {
		fmt.Printf("Failed to find 'cpu_monitor' program\n")
		return
	}

	hook, err := link.Tracepoint("sched", "sched_switch", prog, nil)
	if err != nil {
		fmt.Printf("Failed to attach eBPF program: %v\n", err)
		return
	}
	defer hook.Close()

	fmt.Printf("Monitoring CPU for %s...\n", resourceType)
}

func MonitorMemory(resourceType string) {
	if err := rlimit.RemoveMemlock(); err != nil {
		fmt.Printf("Failed to remove memlock limit: %v\n", err)
		return
	}

	coll, err := LoadMemProgram()
	if err != nil {
		fmt.Printf("Failed to load eBPF program: %v\n", err)
		return
	}
	defer coll.Close()

	prog := coll.Programs["mem_monitor"]
	if prog == nil {
		fmt.Printf("Failed to find 'mem_monitor' program\n")
		return
	}

	hook, err := link.Tracepoint("kmem", "mm_page_alloc", prog, nil)
	if err != nil {
		fmt.Printf("Failed to attach eBPF program: %v\n", err)
		return
	}
	defer hook.Close()

	fmt.Printf("Monitoring memory for %s...\n", resourceType)
}
