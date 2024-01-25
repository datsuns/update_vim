package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	cp "github.com/otiai10/copy"
)

var (
	rootDirectory = filepath.Join(os.Getenv("HOME"), "cwork", "programming", "vim")

	workigDirectory  = filepath.Join(rootDirectory, "vim")
	runtimeDirectory = filepath.Join(workigDirectory, "runtime")
	buildDirectory   = filepath.Join(workigDirectory, "src")

	installRoot = filepath.Join(rootDirectory, "install")
	installDir  = filepath.Join(installRoot, "vim")

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

func copy(src, dest string) error {
	fmt.Printf("copy [%v] -> [%v]\n", src, dest)
	err := cp.Copy(src, dest)
	if err != nil {
		panic(err)
	}
	return err
}

func run_install(srcRoot, installRoot string) {
	copy(filepath.Join(srcRoot, "runtime"), installRoot)
	list, _ := filepath.Glob(filepath.Join(srcRoot, "src", "*.exe"))
	for _, l := range list {
		f := filepath.Base(l)
		copy(l, filepath.Join(installRoot, f))
	}
	copy(filepath.Join(srcRoot, "src", "tee", "tee.exe"), filepath.Join(installRoot, "tee.exe"))
	copy(filepath.Join(srcRoot, "src", "xxd", "xxd.exe"), filepath.Join(installRoot, "xxd.exe"))
	os.MkdirAll(filepath.Join(installRoot, "GvimExt32"), 0770)
	os.MkdirAll(filepath.Join(installRoot, "GvimExt64"), 0770)
	extlib := filepath.Join(srcRoot, "src", "GvimExt", "gvimext.dll")
	copy(extlib, filepath.Join(installRoot, "GvimExt32", "gvimext.dll"))
	copy(extlib, filepath.Join(installRoot, "GvimExt64", "gvimext.dll"))
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
	run_install(workigDirectory, installDir)
	run_build_cmd(pluginUpdateComands)
}
