package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

const DownloadedNameFile = "downloaded.txt"

var (
	cont                 int
	contAllVideos        int
	contVideosDownloaded int
	downloaded           = make(map[string]bool)
	releasesURL          = "https://github.com/yt-dlp/yt-dlp-nightly-builds/releases"
	stateFilePath        = "yt-dlp-version-check.json"
	dateLayout           = "2006-01-02"
	downloadBase         = "https://github.com/yt-dlp/yt-dlp-nightly-builds/releases/download"
	binaryFileName       = "yt-dlp.exe"
)

type VersionState struct {
	CurrentVersion string    `json:"current_version"`
	LastChecked    time.Time `json:"last_checked"`
}

type Root struct {
	Roots struct {
		BookmarkBar Node `json:"bookmark_bar"`
		Other       Node `json:"other"`
	} `json:"roots"`
}

type Node struct {
	Children []Node `json:"children"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	URL      string `json:"url"`
}

func sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func writeInDownloadedFile(url string) error {
	file, err := os.OpenFile(DownloadedNameFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(url + "\n")
	return err
}

func writeInErrorFile(errMsg string) {
	file, err := os.OpenFile("error.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(errMsg + "\n")
	if err != nil {
		fmt.Println(err)
	}
}

func createFile(item Node, path string) {
	fmt.Printf("Downloading %s in %s \n", item.Name, path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			writeInErrorFile(err.Error())
			cont--
			return
		}
	}

	cmd := exec.Command("./yt-dlp.exe", "-o", fmt.Sprintf("%s/%%(title)s.mp4", path), "--restrict-filenames", "--no-check-certificate", "-f", "bestvideo[height<=1080]+bestaudio/best[height<=1080]", item.URL)
	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := fmt.Sprintf("Error downloading %s. Error: %s. Output: %s \n", item.Name, err.Error(), string(output))
		log.Printf(msg)
		writeInErrorFile(msg)
		cont--
		return
	}

	contVideosDownloaded++
	fmt.Printf("%d/%d\n", contVideosDownloaded, contAllVideos)

	downloaded[item.URL] = true
	err = writeInDownloadedFile(item.URL)
	cont--
	if err != nil {
		writeInErrorFile(err.Error())
		return
	}
}

func tree(node []Node, path string) {
	for _, n := range node {
		for cont >= 3 {
			sleep(1000)
		}
		switch n.Type {
		case "folder":
			newPath := path + "/" + n.Name
			tree(n.Children, newPath)
		case "url":
			if strings.Contains(n.URL, "youtube.com") {
				if _, ok := downloaded[n.URL]; !ok {
					cont++
					go createFile(n, path)
				}
			}
		}
	}
}

func getBookmarksPath() (string, error) {
	var path string
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "windows":
		path = filepath.Join(homeDir, "AppData", "Local", "Google", "Chrome", "User Data", "Default", "Bookmarks")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			path = filepath.Join(homeDir, "AppData", "Local", "Google", "Chrome", "User Data", "Profile 1", "Bookmarks")
		}
	case "darwin":
		path = filepath.Join(homeDir, "Library", "Application Support", "Google", "Chrome", "Default", "Bookmarks")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			path = filepath.Join(homeDir, "Library", "Application Support", "Google", "Chrome", "Profile1", "Bookmarks")
		}
	case "linux":
		path = filepath.Join(homeDir, ".config", "google-chrome", "Default", "Bookmarks")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			path = filepath.Join(homeDir, ".config", "chromium", "Default", "Bookmarks")
			if _, err := os.Stat(path); os.IsNotExist(err) {
				path = filepath.Join(homeDir, ".config", "google-chrome", "Profile 1", "Bookmarks")
				if _, err := os.Stat(path); os.IsNotExist(err) {
					path = filepath.Join(homeDir, ".config", "chromium", "Profile 1", "Bookmarks")
				}
			}
		}

	}

	return path, nil
}

func getBookmarkJSON() (Root, error) {
	bookmarksPath, err := getBookmarksPath()
	if err != nil {
		return Root{}, err
	}

	bookmarksFile, err := os.Open(bookmarksPath)
	defer bookmarksFile.Close()
	if err != nil {
		return Root{}, err

	}

	byteValue, _ := io.ReadAll(bookmarksFile)

	var root Root
	err = json.Unmarshal(byteValue, &root)
	if err != nil {
		return Root{}, err
	}

	return root, nil
}

func loadDownloadedFiles() error {
	file, err := os.Open(DownloadedNameFile)
	defer file.Close()
	if err != nil {
		return nil
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		downloaded[line] = true
	}

	return scanner.Err()
}

func fetchLatestVersion() (string, error) {
	resp, err := http.Get(releasesURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`/yt-dlp/yt-dlp-nightly-builds/releases/tag/([0-9\.]+)`)
	match := re.FindStringSubmatch(string(body))
	if len(match) < 2 {
		return "", fmt.Errorf("no se encontró versión en el HTML")
	}

	return match[1], nil
}

func downloadLatestBinary(version string) error {
	url := fmt.Sprintf("%s/%s/%s", downloadBase, version, binaryFileName)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(binaryFileName)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	err = os.Chmod(binaryFileName, 0755)
	if err != nil {
		return err
	}

	return nil
}

func loadState() VersionState {
	var state VersionState
	file, err := os.Open(stateFilePath)
	if err != nil {
		return state
	}
	defer file.Close()

	json.NewDecoder(file).Decode(&state)
	return state
}

func saveState(state VersionState) error {
	file, err := os.Create(stateFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	json.NewEncoder(file).Encode(state)
	return nil
}

func downloadYTD() error {
	state := loadState()
	today := time.Now().Format(dateLayout)
	if state.LastChecked.Format(dateLayout) == today {
		return nil
	}

	latestVersion, _ := fetchLatestVersion()
	if latestVersion != state.CurrentVersion {
		fmt.Printf("Current yt-dlp version: %s\n", state.CurrentVersion)
		fmt.Printf("Downloading yt-dlp version: %s\n", latestVersion)

		err := downloadLatestBinary(latestVersion)
		if err != nil {
			return err
		}

		state.CurrentVersion = latestVersion
	}

	state.LastChecked = time.Now()
	return saveState(state)
}

func countVideos(node []Node) {
	for _, n := range node {
		switch n.Type {
		case "folder":
			countVideos(n.Children)
		case "url":
			if strings.Contains(n.URL, "youtube.com") {
				if _, ok := downloaded[n.URL]; !ok {
					contAllVideos++
				}
			}
		}
	}
}

func main() {
	err := downloadYTD()
	if err != nil {
		fmt.Println(err)
		bufio.NewReader(os.Stdin).ReadString('\n')
		return
	}

	err = loadDownloadedFiles()
	if err != nil {
		fmt.Println(err)
		bufio.NewReader(os.Stdin).ReadString('\n')
		return
	}

	root, err := getBookmarkJSON()
	if err != nil {
		fmt.Println(err)
		bufio.NewReader(os.Stdin).ReadString('\n')
		return
	}

	allChildren := root.Roots.BookmarkBar.Children
	allChildren = append(allChildren, root.Roots.Other.Children...)
	if len(allChildren) > 0 {
		countVideos(allChildren)
		tree(allChildren, "./bookmarks")
		for cont != 0 {
			sleep(1000)
		}
	}
}
