// Binary x86 explains an x86/x86-64 assembly instruction.
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/StalkR/goircbot/lib/x86"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %v <instruction>\n", os.Args[0])
		os.Exit(1)
	}

	ref, err := x86.New()
	if err != nil {
		log.Fatal(err)
	}
	mnemonics := ref.Mnemonics()
	instr := strings.ToUpper(os.Args[1])
	desc, ok := mnemonics[instr]
	if !ok {
		fmt.Println("not found")
		os.Exit(1)
	}
	fmt.Printf("%s: %s\n", instr, desc)
}
