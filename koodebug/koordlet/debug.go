package koordlet

import "os"

func init() {
	os.Args = append(os.Args,
		"--prediction-checkpoint-filepath=/tmp/prediction-checkpoint",
		"--runtime-hooks-addr=/etc/runtime/hookserver.d/koordlet.sock",
	)

}
func Init() {
	name, _ := os.Hostname()
	os.Setenv("NODE_NAME", name)
	os.Setenv("agent_mode", "hostMode")
}
