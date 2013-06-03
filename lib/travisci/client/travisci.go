package main

import (
	"fmt"
	"os"

	"github.com/StalkR/goircbot/lib/travisci"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %v <user> <repo>)\n", os.Args[0])
		os.Exit(1)
	}
	user, repo := os.Args[1], os.Args[2]

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
		status := "passed"
		if !b.Success {
			status = "errored"
		}
		fmt.Printf(" - Build #%v: %v (%v) %v (%v/%v)\n", b.Number, status, b.Finished,
			b.Message, b.Branch, b.Commit)
	}
}
