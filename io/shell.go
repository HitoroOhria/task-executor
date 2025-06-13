package io

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-task/task/v3/errors"
)

const (
	incrementalSearchTool = "peco"
	maxVarPromptWidth     = 18
)

var (
	searchTaskfiles = []string{"Taskfile.yml", "Taskfile.yaml"}

	selectTaskNameCommand = func(taskfile string) string {
		return fmt.Sprintf(`
	  task -t %s -l --sort none | \
	    tail -n +2 | \
	    sed 's/^\*//' | \
	    %s | \
	    sed -E 's/^ ([^ ]+):.*/\1/' | \
	    sed -E 's/:$//'
	`, taskfile, incrementalSearchTool)
	}

	ErrFileNotFound              = errors.New("file not found")
	ErrTaskfileNotFound          = errors.New("taskfile not found")
	ErrCanceledIncrementalSearch = errors.New("canceled incremental search")
)

// findFileName はファイルパスからファイル名を取得する
// ファイルパスは相対パスでも絶対パスでも良い
// ファイルが見つからなかった場合、ErrFileNotFound を返す
func findFileName(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("%w: path = %s", ErrFileNotFound, path)
		}

		return "", fmt.Errorf("os.Stat: %w", err)
	}

	if info.IsDir() {
		return "", fmt.Errorf("%w: target is directory. path = %s", ErrFileNotFound, path)
	}

	return info.Name(), nil
}

// FindTaskfileName は、カレントディレクトリの Taskfile を探索し、ファイル名を返却する
func FindTaskfileName() (string, error) {
	taskfileName := ""
	for _, taskfile := range searchTaskfiles {
		found, err := findFileName(taskfile)
		if err != nil {
			// 期待するファイル名をループ探索しているので、ファイルが見つからなくてもスルーする
			if errors.Is(err, ErrFileNotFound) {
				continue
			}

			return "", fmt.Errorf("findFileName: %w", err)
		}

		taskfileName = found
		break
	}

	return taskfileName, nil
}

func SelectTaskName(taskfile string) (string, error) {
	_, err := findFileName(taskfile)
	if err != nil {
		if errors.Is(err, ErrFileNotFound) {
			return "", fmt.Errorf("%w: taskfile = %s", ErrTaskfileNotFound, taskfile)
		}

		return "", fmt.Errorf("findFileName: %w", err)
	}

	cmd := exec.Command("sh", "-c", selectTaskNameCommand(taskfile))
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("cmd.Output: %w", err)
	}

	o := strings.TrimSpace(string(output))
	// インクリメンタルサーチ中に Ctrl + C で中断された場合は、特定のエラーを返す
	if o == "" {
		return "", ErrCanceledIncrementalSearch
	}

	return o, nil
}

func ReadFile(path string) ([]byte, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}

	return file, nil
}

func Prompt(maxNameLen int, varName string) string {
	pad := maxNameLen + 2 // plus double quote
	if pad > maxVarPromptWidth {
		pad = maxVarPromptWidth
	}
	name := fmt.Sprintf(`"%s"`, varName)

	return fmt.Sprintf(`Enter %-*s: `, pad, name)
}

func Input(prompt string) string {
	fmt.Print(prompt)
	return readInput()
}

// readInput は値の入力を受け付ける
// Ctrl + C でキャンセルされた場合は、プログラムを正常終了する
// FIXME context & goroutine を使用した方法もあるので、検討する
func readInput() string {
	// Ctrl+C (SIGINT) を補足
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)

	// 入力がキャンセルされた場合、何もしない
	go func() {
		<-sigCh
		fmt.Println()
		os.Exit(0)
	}()

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		// 入力が完了したらシグナルハンドラーを停止
		signal.Reset(syscall.SIGINT)

		return scanner.Text()
	}

	err := scanner.Err()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: scanner.Err():", err)
	}

	return ""
}

func RunTask(taskfile string, name string, args ...string) error {
	fmt.Printf("run: task -t %s %s %s\n", taskfile, name, strings.Join(args, " "))

	cmdArgs := append([]string{"-t", taskfile, name}, args...)
	cmd := exec.Command("task", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("cmd.Run: %w", err)
	}

	return nil
}
