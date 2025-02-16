//go:build ignore

#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <linux/types.h>

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __type(key, u32);
    __type(value, u64);
    __uint(max_entries, 1024);
} cpu_usage SEC(".maps");

SEC("tracepoint/sched/sched_switch")
int cpu_monitor(struct trace_event_raw_sched_switch *ctx) {
    u32 pid = bpf_get_current_pid_tgid();
    u64 *usage = bpf_map_lookup_elem(&cpu_usage, &pid);
    if (usage) {
        (*usage)++;
    } else {
        u64 init_usage = 1;
        bpf_map_update_elem(&cpu_usage, &pid, &init_usage, BPF_NOEXIST);
    }
    return 0;
}

char _license[] SEC("license") = "GPL";