package utils

import (
	"fmt"
	"io/ioutil"
)

// CheckDirectoryContents는 주어진 경로의 디렉토리 및 파일 목록을 출력합니다.
func CheckDirectoryContents(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", path, err)
	}

	var contents []string
	fmt.Printf("Contents of %s:\n", path)
	for _, file := range files {
		if file.IsDir() {
			fmt.Printf("[DIR]  %s\n", file.Name())
			contents = append(contents, fmt.Sprintf("[DIR] %s", file.Name()))
		} else {
			fmt.Printf("[FILE] %s\n", file.Name())
			contents = append(contents, fmt.Sprintf("[FILE] %s", file.Name()))
		}
	}

	return contents, nil
}
