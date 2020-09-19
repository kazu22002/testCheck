package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
)

func main(){
	paths := inputFile()

	for _, path := range paths {
		data := useBufioScanner(path)
		fmt.Printf("input: %s\n", data)

		// テスト実行
		out := stdinPipe(data)
		fmt.Printf("結果: %s", out)

		err := outputFile(path, out)
		if err != nil {
			fmt.Printf("エラー: %s", err.Error())
		}

		if checkTest(path) {
			fmt.Println("OK")
		} else {
			fmt.Println("NG")
		}

		fmt.Println()
	}
}

func stdinPipe(in []string) []byte {
	cmd := exec.Command("./main")
	stdin, _ := cmd.StdinPipe()
	for _, text := range in {
		io.WriteString(stdin, text)
	}
	stdin.Close()
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("エラー: %s", err.Error())
	}
	return out
}

func inputFile() []string {
	paths := dirwalk("./tests")
	ret := []string{}
	for _, path := range paths {
		e := filepath.Ext(path)
		if e == ".in" {
			ret = append(ret, path)
		}
	}
	return ret
}

func outputFile(fileName string, out []byte) error {
	outFile := outTestFilePath(fileName)

	if !exists("out") {
		if err := os.Mkdir("out", 0777); err != nil {
			fmt.Println(err)
		}
	}
	outPath := "out/" + outFile
	if exists(outPath) {
		if err := os.Remove(outPath); err != nil {
			fmt.Println(err)
		}
	}

	err := ioutil.WriteFile(outPath, out, 0666)
	if err != nil {
		fmt.Println(os.Stderr, err)
//		os.Exit(1)
		return err
	}
	return nil
}

func exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func checkTest(fileName string) bool {
	outFile := outTestFilePath(fileName)
	checkPath := "tests/" + outFile
	outPath := "out/" + outFile

	checkBuf := useBufioScanner(checkPath)
	outBuf := useBufioScanner(outPath)

	if reflect.DeepEqual(checkBuf,outBuf) {
		return true
	}

	return false
}

func dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths
}

func useBufioScanner(fileName string) []string{
	fp, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	ret := []string{}
	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}
	return ret
}

func outTestFilePath(fileName string) string {
	fileBase := baseFileName(fileName)
	return fileBase + ".out"
}

func baseFileName(fileName string) string {
	_, f := filepath.Split(fileName)
	for i := len(f) - 1; i >= 0; i-- {
		if f[i] == '.' {
			return f[0:i]
		}
	}
	return ""
}