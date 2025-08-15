package agent

import (
	"fmt"
	"os/exec"
)

// i do not know how to structure this function
// my javascript brain cannot comprehend types anymore,
// this json shit aint for me man

// Some things to note
// Should definitely specify which shell is being used
// and with bash -c configuration, i think i wont need to parse the args and create an array
func runTerminalCommand() {

	cmd := exec.Command("bash", "-c", "echo")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))
}
