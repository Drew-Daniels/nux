// nux is a modern tmux session manager that builds sessions declaratively
// from YAML configs. It supports session groups, glob patterns, and
// multi-session operations. Install with: go install github.com/Drew-Daniels/nux@latest
package main

import "github.com/Drew-Daniels/nux/cmd"

func main() {
	cmd.Execute()
}
