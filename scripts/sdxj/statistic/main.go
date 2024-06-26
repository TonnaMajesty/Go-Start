package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func main() {
	rootDir := "/Users/tonnamajesty/Downloads/siji"
	devices, err := getDeviceFolders(rootDir)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, device := range devices {
		devicePath := filepath.Join(rootDir, device)
		channelFolders, err := getChannelFolders(devicePath)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		for _, channel := range channelFolders {
			channelPath := filepath.Join(devicePath, channel)
			imageFiles, err := getImageFiles(channelPath)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			if len(imageFiles) <= 2 {
				continue
			}

			sort.Strings(imageFiles) // Sort files by name (time)

			avgInterval := calculateInterval(imageFiles)
			fmt.Printf("Device:  %s Interval: %v \n", channel, avgInterval)
		}
	}
}

func getDeviceFolders(root string) ([]string, error) {
	folders, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var devices []string
	for _, folder := range folders {
		if folder.IsDir() {
			devices = append(devices, folder.Name())
		}
	}
	return devices, nil
}

func getChannelFolders(devicePath string) ([]string, error) {
	folders, err := ioutil.ReadDir(devicePath)
	if err != nil {
		return nil, err
	}

	var channels []string
	for _, folder := range folders {
		if folder.IsDir() {
			channels = append(channels, folder.Name())
		}
	}
	return channels, nil
}

func getImageFiles(channelPath string) ([]string, error) {
	files, err := ioutil.ReadDir(channelPath)

	if err != nil {
		return nil, err
	}

	var imageFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".jpg") {
			imageFiles = append(imageFiles, file.Name())
		}
	}
	return imageFiles, nil
}

func parseImageFileName(filename string) (time.Time, error) {
	parts := strings.Split(strings.TrimSuffix(filename, ".jpg"), "_")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid file name format")
	}

	timestampStr := parts[1]
	timestamp, err := time.Parse("20060102150405", timestampStr)
	if err != nil {
		return time.Time{}, err
	}

	return timestamp, nil
}

func calculateInterval(imageFiles []string) time.Duration {
	if len(imageFiles) < 2 {
		return 0
	}

	var intervals []time.Duration
	prevTime, _ := parseImageFileName(imageFiles[0])

	for _, imageFile := range imageFiles[1:] {
		timestamp, err := parseImageFileName(imageFile)
		if err != nil {
			fmt.Println("Error parsing file name:", err)
			continue
		}

		interval := timestamp.Sub(prevTime)
		intervals = append(intervals, interval)

		prevTime = timestamp
	}

	// Calculate the average interval
	var totalInterval time.Duration
	for _, interval := range intervals {
		totalInterval += interval
	}
	averageInterval := totalInterval / time.Duration(len(intervals))

	return averageInterval
}
