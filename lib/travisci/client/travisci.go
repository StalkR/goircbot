package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/StalkR/goircbot/lib/travisci"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %v <user>/<repo>)\n", os.Args[0])
		os.Exit(1)
	}
	userRepo := strings.SplitN(os.Args[1], "/", 2)
	if len(userRepo) != 2 {
		fmt.Printf("Invalid user/repo: %v\n", os.Args[1])
		os.Exit(1)
	}
	user, repo := userRepo[0], userRepo[1]

	builds, err := travisci.Builds(user, repo)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if len(builds) == 0 {
		fmt.Printf("No such user/repo or no build yet: https://www.travis-ci.org/%v/%v\n",
			user, repo)
		os.Exit(1)
	}

	fmt.Printf("Builds for %v/%v:\n", user, repo)
	for i, j := 0, len(builds)-1; i < j; i, j = i+1, j-1 {
		builds[i], builds[j] = builds[j], builds[i]
	}
	for _, b := range builds {
		var status string
		if b.State == "finished" {
			status = "passed"
			if !b.Success {
				status = "errored"
			}
		} else {
			status = "in progress"
		}
		if b.Finished.IsZero() {
			fmt.Printf(" - Build #%v: %v, %v (%v/%v)\n", b.Number, status,
				b.Message, b.Branch, b.Commit)
		} else {
			fmt.Printf(" - Build #%v: %v (%v) %v (%v/%v)\n", b.Number, status, b.Finished,
				b.Message, b.Branch, b.Commit)
		}
	}
}
