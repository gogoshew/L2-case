package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах

Реализовать утилиту netcat (nc) клиент
принимать данные из stdin и отправлять в соединение (tcp/udp)
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func commandExec(command string) error {
	command = strings.TrimSuffix(command, "\n")
	comArgs := strings.Split(command, " ")

	switch comArgs[0] {
	case "cd":
		if len(comArgs) < 2 {
			return errors.New("path required")
		}
		return os.Chdir(comArgs[1])
	case "exit":
		os.Exit(0)
	}

	cmd := exec.Command(comArgs[0], comArgs[1:]...)

	cmd.Stderr, cmd.Stdout = os.Stderr, os.Stdout

	return cmd.Run()
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		cmd, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		if err = commandExec(cmd); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
