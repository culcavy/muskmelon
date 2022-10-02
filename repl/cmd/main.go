package main

import (
	"fmt"
	"github.com/hollykbuck/muskmelon/repl"
	"log"
	"os"
	"os/user"
)

func _main() error {
	current, err := user.Current()
	if err != nil {
		return fmt.Errorf("获取当前用户失败: %w", err)
	}
	fmt.Printf("Hello %s! This is muskmelon programming language!\n", current.Username)
	fmt.Printf("feel free to type in commands\n")
	err = repl.Start(os.Stdin, os.Stdout)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := _main()
	if err != nil {
		log.Println(err)
		return
	}
}
