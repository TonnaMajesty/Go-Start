package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// 指定要读取的文件夹路径
	folderPath := "/Users/tonnamajesty/Downloads/sichuan_batch1_site659_image659_segBase_20240103_label"

	// 读取指定文件夹的子文件夹列表
	subDirs, err := ioutil.ReadDir(folderPath)
	if err != nil {
		fmt.Printf("Error reading directory: %s\n", err)
		return
	}

	for _, subDir := range subDirs {
		// 检查是否是子文件夹
		if subDir.IsDir() {
			// 获取子文件夹名称
			dirName := subDir.Name()

			// 修改子文件夹名称格式
			newDirName := strings.ReplaceAll(dirName, "_", "--")

			// 构建新的子文件夹路径
			newPath := filepath.Join(folderPath, newDirName)

			// 重命名子文件夹
			err := os.Rename(filepath.Join(folderPath, dirName), newPath)
			if err != nil {
				fmt.Printf("Failed to rename folder %s: %s\n", dirName, err)
			} else {
				fmt.Printf("Renamed folder: %s -> %s\n", dirName, newDirName)
			}
		}
	}
}
