package mem

import (
	"fmt"

	"github.com/saiaunghlyanhtet/k-moon/pkg/bpf"
	"github.com/saiaunghlyanhtet/k-moon/pkg/k8s"
	"github.com/spf13/cobra"
)

var (
	nodeFlag      bool
	podFlag       bool
	containerFlag bool
)

var Cmd = &cobra.Command{
	Use:   "mem",
	Short: "Monitor memory usage for nodes, pods, or containers",
	Run: func(cmd *cobra.Command, args []string) {
		clientset, err := k8s.NewClient()
		if err != nil {
			fmt.Printf("Failed to create Kubernetes client: %v\n", err)
			return
		}

		if nodeFlag {
			nodes, err := k8s.GetNodes(clientset)
			if err != nil {
				fmt.Printf("Failed to fetch nodes: %v\n", err)
				return
			}
			fmt.Println("Nodes:", len(nodes))
			bpf.MonitorMemory("node")
		} else if podFlag {
			pods, err := k8s.GetPods(clientset)
			if err != nil {
				fmt.Printf("Failed to fetch pods: %v\n", err)
				return
			}
			fmt.Println("Pods:", len(pods))
			bpf.MonitorMemory("pod")
		} else if containerFlag {
			pods, err := k8s.GetPods(clientset)
			if err != nil {
				fmt.Printf("Failed to fetch pods: %v\n", err)
				return
			}
			for _, pod := range pods {
				containers := k8s.GetContainers(pod)
				fmt.Println("Containers in pod", pod.Name, ":", len(containers))
			}
			bpf.MonitorMemory("container")
		} else {
			fmt.Println("Please specify a resource type: --node, --pod, or --container")
		}
	},
}

func init() {
	Cmd.Flags().BoolVar(&nodeFlag, "node", false, "Monitor memory usage for nodes")
	Cmd.Flags().BoolVar(&podFlag, "pod", false, "Monitor memory usage for pods")
	Cmd.Flags().BoolVar(&containerFlag, "container", false, "Monitor memory usage for containers")
}
