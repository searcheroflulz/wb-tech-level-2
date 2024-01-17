package main

import (
	"bufio"
	"bytes"
	"fmt"
	gops "github.com/mitchellh/go-ps"
	"log"
	"os"
	osExec "os/exec"
	"path/filepath"
	"strconv"
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

// меняем текущую директорию на выбранную пользователем
func cd(dir string) error {
	err := os.Chdir(dir)
	if err != nil {
		return err
	}
	return nil
}

// выводим пользователю текущую директорию
func pwd() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return fmt.Sprint(dir), nil
}

// возвращаем тект пользователя
func echo(text string) string {
	return fmt.Sprintln(text)
}

// убиваем процесс по pid
func kill(pid string) error {
	piDig, err := strconv.Atoi(pid)
	if err != nil {
		return err
	}
	process, err := os.FindProcess(piDig)
	if err != nil {
		return err
	}
	return process.Kill()
}

// выводим все процессы на компьютере списком
func ps() (string, error) {
	var result string
	processes, err := gops.Processes()
	if err != nil {
		return "", err
	}
	result = fmt.Sprintln("   PID\t| Executable")
	for i := range processes {
		proc := processes[i]
		result = result + fmt.Sprintf("%d\t| %s\n", proc.Pid(), proc.Executable())
	}
	return result, nil
}

// проверка на наличие файла
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// запускаем файл
func exec(arg string) (*bytes.Buffer, error) {
	var output bytes.Buffer
	var filePath string

	args := strings.Split(arg, " ")
	//если путь до файла не абсолютный, подставляем текущую директорию
	if !filepath.IsAbs(args[0]) {
		absPath, err := filepath.Abs(filepath.Join(".", args[0]))
		if err != nil {
			return nil, err
		}
		filePath = absPath
	}
	//если файла нет по абсолютному пути, пытаемся найти его в $PATH
	if !fileExists(filePath) {
		cmdPath, err := osExec.LookPath(args[0])
		if err != nil {
			return nil, err
		}
		cmd := osExec.Command(cmdPath, args[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = &output
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return nil, err
		}

		return &output, nil
	}

	//если файл все-таки находится в текущей директории, запускаем его
	cmd := osExec.Command(filePath, args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = &output
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func runCommands(command string, output *bytes.Buffer, pipes bool) (*bytes.Buffer, error) {
	command = strings.TrimLeft(command, " ")
	command = strings.TrimSuffix(command, "\r")
	args := strings.Split(command, " ")

	switch args[0] {
	case "cd":
		if pipes {
			return output, cd(output.String())
		}
		return output, cd(args[1])
	case "pwd":
		s, err := pwd()
		output.Reset()
		output.WriteString(s)
		return output, err
	case "echo":
		if pipes {
			s := echo(output.String())
			output.Reset()
			output.WriteString(s)
			return output, nil
		}
		output.WriteString(echo(strings.Join(args[1:], " ")))
		return output, nil
	case "kill":
		return output, kill(args[1])
	case "ps":
		s, err := ps()
		output.Reset()
		output.WriteString(s)
		return output, err
	//case "fork":
	case "exec":
		buffer, err := exec(args[1])
		if err != nil {
			return nil, err
		}
		output = buffer
		return output, nil
	case "\\quit":
		os.Exit(0)
	default:
		fmt.Println("Unknown command")
	}
	return nil, nil
}

func checkCommands(commands string) error {
	coms := strings.Split(commands, "|")
	var pipes bool
	countPipes := len(coms) - 1
	if countPipes != 0 {
		pipes = true
	}
	var output *bytes.Buffer
	var err error

	output = new(bytes.Buffer)

	for _, command := range coms {
		output, err = runCommands(command, output, pipes)
		if err != nil {
			return err
		}
		if countPipes == 0 {
			if output.String() == "" {
				return nil
			}
			fmt.Println(output.String())
			return nil
		}
		if countPipes > 0 {
			countPipes--
		}
	}
	if output == nil || len(output.String()) == 0 {
		fmt.Fprint(os.Stdout, "")
	} else {
		fmt.Fprint(os.Stdout, output.String())
	}

	return nil
}

func main() {
	//считываем информарцию с консоли
	reader := bufio.NewReader(os.Stdin)

	for {
		//выдаем пользователю информацию о текущей директории
		dir, err := pwd()
		fmt.Printf("%v$ ", dir)
		//считываем пользовательские команды
		commands, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		commands = strings.TrimSuffix(commands, "\n")
		//выполняем команды пользователя
		err = checkCommands(commands)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
