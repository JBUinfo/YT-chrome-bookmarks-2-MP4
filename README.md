# YT-chrome-bookmarks-2-MP4
Download Youtube videos you have saved in your Chrome bookmarks. Bookmarks2MP4

# Install
You need install [Node](https://nodejs.org/es/) and Microsoft Visual C++ 2010 ("[vcredist_x86](https://github.com/JBUinfo/YT-chrome-bookmarks-2-file/blob/main/vcredist_x86.exe)" file).

# Instructions
1. Open index.html and upload your bookmarks.htm.

(index.html is a custom version of [chrome-bookmarks-converter-master](https://github.com/jsnelders/chrome-bookmarks-converter) )

2. It will download "downloaded.json" file. Put it in the folder where "script_Download_YT_Bookmarks.js" is.

3. Open CMD, go to the folder and execute "node script_Download_YT_Bookmarks.js".

# Result
You will have all the folders and videos you have in bookmarks.

# Notes
The script will download 3 videos simultaneously.

The file "downloaded.txt" will be created. It will have the URLs of the videos have been downloaded.

The file "errs.txt" will be created. It will have errors that "youtube-dl" throws.
