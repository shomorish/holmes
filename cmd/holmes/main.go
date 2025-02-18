package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func getLogoFromFile() (string, error) {
	content, err := os.ReadFile("logo.txt")
	if err != nil {
		return "", err
	}
	logo := strings.TrimSuffix(string(content), "\n")
	return logo, nil
}

func getHelpFromFile() (string, error) {
	content, err := os.ReadFile("help.txt")
	if err != nil {
		return "", err
	}
	help := strings.TrimSuffix(string(content), "\n")
	return help, nil
}

func containsString(path, searchString string) bool {
	content, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return strings.Contains(string(content), searchString)
}

func readDir(dir string, searchString string, wg *sync.WaitGroup) {
	defer wg.Done()

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("Failed to read `%s`", dir)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			wg.Add(1)
			go readDir(filepath.Join(dir, file.Name()), searchString, wg)
		} else if containsString(filepath.Join(dir, file.Name()), searchString) {
			fmt.Println(filepath.Join(dir, file.Name()))
		}
	}
}

func main() {
	logo, err := getLogoFromFile()
	if err != nil {
		log.Fatalf("Failed to get logo: %v", err)
	}
	fmt.Println(logo)
	fmt.Println()

	help, err := getHelpFromFile()
	if err != nil {
		log.Fatalf("Failed to get help: %v", err)
	}
	fmt.Println(help)
	fmt.Println()

	// 現在のディレクトリを取得
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}
	fmt.Printf("Current Directory: %s\n", currentDir)
	fmt.Println()

	// 入力された1行を読み込む。スペースも使えるように`bufio`パッケージを使用
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(">>> ")

quit: // コマンドでアプリを終了するために使用。コマンドの判別にswitchを使用しており、通常のbreakではループを抜けることができないためラベルが必要
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 1 {
			// 何も入力されていないためスキップ
			fmt.Print(">>> ")
			continue
		}

		// 文字列とコマンドで処理を分ける
		if len(line) > 1 && line[0] == '@' {
			// コマンド
			command := line[1]
			switch command {
			case 'c':
				splited := strings.Split(line, " ")
				fmt.Printf("splited(%d): %v\n", len(splited), splited)
				if len(splited) < 2 { // ディレクトリが入力されているか確認
					fmt.Println("Enter the directory name.")
				} else if err := os.Chdir(splited[1]); err != nil { // 指定したディレクトリが存在するか確認
					fmt.Println("It was an invalid directory.")
				} else { // 入力されたディレクトリが存在する
					currentDir, err = os.Getwd()
					if err != nil {
						log.Fatalf("Failed to get current directory: %v", err)
					}
				}
				fmt.Printf("Current Directory: %s\n", currentDir)
			case 'p':
				fmt.Printf("Current Directory: %s\n", currentDir)
			case 'h':
				fmt.Println(help)
			case 'q':
				// ループを抜けてアプリを終了
				break quit
			default:
				fmt.Println("Unexpected command. Type `@h` for help.")
			}
		} else {
			// 検索文字列
			var wg sync.WaitGroup
			wg.Add(1)
			go readDir(currentDir, line, &wg)
			wg.Wait()
		}
		fmt.Println()
		fmt.Print(">>> ")
	}

	fmt.Println()
	fmt.Println("See you!")
}
