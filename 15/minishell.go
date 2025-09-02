package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	signal.Ignore(syscall.SIGINT)

	scanner := bufio.NewScanner(os.Stdin)
	stat, _ := os.Stdin.Stat()
	interactive := (stat.Mode() & os.ModeCharDevice) != 0

	for {
		if interactive {
			fmt.Print("minish> ")
		}
		if !scanner.Scan() {
			if interactive {
				fmt.Println()
			}
			return
		}

		line := scanner.Text()
		if line == "" {
			continue
		}

		signal.Reset(syscall.SIGINT)
		err := execute(line)
		signal.Ignore(syscall.SIGINT)

		if err != nil {
			fmt.Fprintln(os.Stderr, "minish:", err)
		}

		if !interactive {
			return
		}
	}
}

func runPipeline(line string) error {
	pipes := strings.Split(line, "|")
	commands := make([]*exec.Cmd, len(pipes))

	for i, part := range pipes {
		args := strings.Fields(strings.TrimSpace(part))
		if len(args) == 0 {
			continue
		}
		cmd := exec.Command(args[0], args[1:]...)
		commands[i] = cmd
	}
	for i := 0; i < len(commands)-1; i++ {
		stdout, err := commands[i].StdoutPipe()
		if err != nil {
			return err
		}
		commands[i+1].Stdin = stdout
	}

	commands[len(commands)-1].Stdout = os.Stdout
	commands[len(commands)-1].Stderr = os.Stderr

	for _, cmd := range commands {
		if err := cmd.Start(); err != nil {
			return err
		}
	}
	for _, cmd := range commands {
		if err := cmd.Wait(); err != nil {
			return err
		}
	}

	return nil
}

func runCommand(line string) error {
	args := strings.Fields(line)
	if len(args) == 0 {
		return nil
	}

	cmd := args[0]

	switch cmd {
	case "cd":
		path := os.Getenv("HOME")
		if len(args) > 1 {
			path = args[1]
		}
		return os.Chdir(path)

	case "pwd":
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		fmt.Println(pwd)
		return nil
	case "echo":
		fmt.Println(strings.TrimPrefix(line, "echo "))
		return nil

	case "kill":
		if len(args) < 2 {
			return fmt.Errorf("kill: no pid")
		}
		pid, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}
		proc, err := os.FindProcess(pid)
		if err != nil {
			return err
		}
		return proc.Signal(os.Kill)

	case "ps":
		c := exec.Command("ps", "aux")
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	default:
		c := exec.Command(args[0], args[1:]...)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	}

}
func execute(line string) error {
	line = os.ExpandEnv(line)
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}
	if strings.Contains(line, "&&") || strings.Contains(line, "||") {
		return runConditional(line)
	}

	if strings.Contains(line, ">") || strings.Contains(line, "<") {
		return runRedirect(line)
	}
	if strings.Contains(line, "|") {
		return runPipeline(line)
	}
	return runCommand(line)
}

func runRedirect(line string) error {
	if strings.Contains(line, ">>") {
		parts := strings.SplitN(line, ">>", 2)
		return handlerAppend(parts[0], parts[1])
	} else if strings.Contains(line, ">") {
		parts := strings.SplitN(line, ">", 2)
		return handlerOutput(parts[0], parts[1])
	} else if strings.Contains(line, "<") {
		parts := strings.SplitN(line, "<", 2)
		return handlerInput(parts[0], parts[1])
	}
	return fmt.Errorf("unknow redirect type")
}

func handlerOutput(cmdLine, filename string) error {
	cmd, err := parseCommand(cmdLine)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	cmd.Stdout = file
	return cmd.Run()
}
func handlerInput(cmdLine, filename string) error {
	cmd, err := parseCommand(cmdLine)
	if err != nil {
		return err
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	cmd.Stdin = file
	return cmd.Run()
}

func handlerAppend(cmdLine, filename string) error {
	cmd, err := parseCommand(cmdLine)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	cmd.Stdout = file
	return cmd.Run()
}

func parseCommand(line string) (*exec.Cmd, error) {
	args := strings.Fields(line)
	if len(args) == 0 {
		return nil, nil
	}
	return exec.Command(args[0], args[1:]...), nil
}

func runConditional(line string) error {
	andParts := strings.Split(line, "&&")
	for i, andPart := range andParts {
		orParts := strings.Split(andPart, "||")
		var lastErr error

		for _, orPart := range orParts {
			orPart = strings.TrimSpace(orPart)
			if orPart == "" {
				continue
			}
			err := execute(orPart)
			if err == nil {
				break
			}
			lastErr = err
		}
		if lastErr != nil {
			if i == len(andParts)-1 {
				return lastErr
			}
			return nil
		}
	}
	return nil
}
