//go:build ignore
#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <linux/types.h>

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__type(key, u32);
	__type(value, u64);
	__uint(max_entries, 1024);
} mem_usage SEC(".maps");

SEC("tracepoint/kmem/mm_page_alloc")
int mem_monitor(struct trace_event_raw_mm_page_alloc *ctx) {
	u32 pid = bpf_get_current_pid_tgid();
	u64 *usage = bpf_map_lookup_elem(&mem_usage, &pid);
	if (usage) {
		(*usage) += ctx->bytes_alloc;
	} else {
		u64 init_usage = ctx->bytes_alloc;
		bpf_map_update_elem(&mem_usage, &pid, &init_usage, BPF_NOEXIST);
	}
	return 0;
}

char _license[] SEC("license") = "GPL";