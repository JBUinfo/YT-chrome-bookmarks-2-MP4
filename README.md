# YT-chrome-bookmarks-2-MP4
Download Youtube videos you have saved in your Google Chrome bookmarks

# Build
1. Install [Golang](https://go.dev/dl/)
2. Open a terminal, go to the project folder and execute: 

    `go build ytcb2mp4.go`

# Install
Download your version from [Releases](https://github.com/JBUinfo/YT-chrome-bookmarks-2-MP4/releases)

# Instructions
1. Run the downloaded/built binary

# Result
You will have all the folders and videos you have in bookmarks.

# Notes
- The script will download 3 videos simultaneously.

- The file "downloaded.txt" will be created. It will have the URLs of the videos have been downloaded.

- The file "errs.txt" will be created. It will have errors that "yt-dlp.exe" throws.

- Under the hood, the script looks for:

#### Windows
    - "$HOME\AppData\Local\Google\Chrome\User Data\Default\Bookmarks"
    OR
    - "$HOME\AppData\Local\Google\Chrome\User Data\Profile 1\Bookmarks"

#### Darwin
    - "$HOME/Library/Application Support/Google/Chrome/Default/Bookmarks"
    OR
    - "$HOME/Library/Application Support/Google/Chrome/Profile1/Bookmarks"

#### Linux
    - "$HOME/.config/google-chrome/Default/Bookmarks"
    OR
    - "$HOME/.config/chromium/Default/Bookmarks"
    OR
    - "$HOME/.config/google-chrome/Profile 1/Bookmarks"
    OR
    - "$HOME/.config/chromium/Profile 1/Bookmarks"
