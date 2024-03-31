package scheduler

import (
	"context"
	"k8s.io/kubernetes/pkg/scheduler"
	"os"
)

func Init() {
	os.Args = append(os.Args,
		"--port=10222",
		"--authentication-skip-lookup=true",
		"--v=4",
		"--feature-gates=DisablePodDisruptionBudgetInformer=true,ResizePod=true",
		"--config=/tmp/koord-scheduler.config",
		"--leader-elect-renew-deadline=1000s",
		"--leader-elect-lease-duration=2000s",
		"--leader-elect",
	)
}

// cp schedule.yaml /tmp/koord-scheduler.config

func main() {
	var sched *scheduler.Scheduler
	ctx := context.Background()
	sched.Run(ctx) // ✅ 队列之间的入队
}
