package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	workigDirectory = filepath.Join(os.Getenv("HOMEPATH"), "cwork", "programming", "vim", "vim")
	buildDirectory  = filepath.Join(workigDirectory, "src")

	gitCommands = [][]string{
		{"git", "reset", "--hard"},
		{"git", "checkout", "master"},
		{"git", "pull", "origin", "master"},
		{"git", "checkout", "my"},
		{"git", "merge", "master"},
	}

	buildCommands = []string{
		"make", "-j", "8", "-f", "Make_ming.mak", "ARCH=x86-64",
	}

	buildCommands_cui = []string{
		"make", "-j", "8", "-f", "Make_ming.mak", "ARCH=x86-64", "GUI=no",
	}

	pluginUpdateComands = []string{
		"./gvim", "-c", "PlugUpgrade", "-c", "PlugUpdate", "-c", "qa",
	}
)

func print_proc(r io.Reader) {
	scanner := bufio.NewScanner(r)
	//scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
}

func execute(bin string, params ...string) int {
	fmt.Println(bin, params)
	cmd := exec.Command(bin, params...)
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	go print_proc(stderr)
	go print_proc(stdout)
	cmd.Wait()
	return cmd.ProcessState.ExitCode()
}

func run_build_cmd(cmd []string) int {
	return execute(cmd[0], cmd[1:]...)
}

func main() {
	fmt.Printf("cd to [%v]\n", workigDirectory)
	os.Chdir(workigDirectory)
	for _, c := range gitCommands {
		if ret := execute(c[0], c[1:]...); ret != 0 {
			fmt.Println("error. return code", ret)
			return
		}
	}
	fmt.Printf("cd to [%v]\n", buildDirectory)
	os.Chdir(buildDirectory)
	run_build_cmd(buildCommands)
	run_build_cmd(buildCommands_cui)
	run_build_cmd(pluginUpdateComands)
}
