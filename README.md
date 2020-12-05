# GO URLs CHECKER CLI

**Description**: This is a basic URL checker written in GO. The app allows users to pass in a file name as a command-line argument, then it would extract all valid URLs and run status check on those urls.

## Installation
1. Clone the repo and build the binary file locally by running
```go
   go build linkDetector.go
```
2. At the folder you want the binary file to be and install it directly through `go get`
```go
   go get github.com/tonyvugithub/GoURLsCheckerCLI
```
**Notice**: This will create a binary name GoURLsCheckerCLI as opposed to linkDetector in the first option. 
## Features
1. Display help panel by not include any argument
   ```go
    ./linkDetector
   ```
2. Display current version using -v or -version flag
   ```go
    ./linkDetector -v or -version
   ```
3. Check multiple files:
   ```go
    //For paths that contain space, simply wrap the path in quotes
    ./linkDetector check -f [file-path-1] [file-path-2] ... [file-path-nth]
   ```
4. Check all files in single or multiple directories:
   ```go
    //For paths that contain space, simply wrap the path in quotes
    ./linkDetector check -d [directory-path-1] [directory-path-2] ... [directory-path-nth]
   ```
5. Allows user to pass glob pattern as argument:<br/>
   Example: Uses as a standalone, The command would be applied to the current directory of the executable file
   ```go
    ./linkDetector -g *.html
   ```
   Example: Uses with -d flag
   ```go
    //The glob pattern needs to be the last argument
    ./linkDetector check -d -g "Absolute\Path\To\Your\Directory" *.txt
   ```
   
6. Create report file by adding -r flag:
   ```go
    ./linkDetector check -f -r [file-path]
   ```
   ```go
    ./linkDetector check -d -r [directory-path]
   ```

7. Use custom User-Agents with the -u flag:
   ```go
    ./linkDetector check -u [User-Agent string] -f [file-path]
   ```
   an example using a Firefox User-Agent:
   ```go
    ./linkDetector check -u "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:81.0) Gecko/20100101 Firefox/81.0" -f [file-path]
   ```

8. Ignore a list of urls:
   ```go
    ./linkDetector check -i [file-path-ignore-list] [file-path]
   ```
   A valid ignore list will have a # for comments and http/https links
   For example:
   ```
   # ignore the canadian google links
   https://google.ca
   # this link wont be ignore https://netflix.com
   https://www.spotify.com/
   ```
9. Check the latest 10 post of [Telescope](https://github.com/Seneca-CDOT/telescope) posts
  ```go
    ./linkDetector check -t 
  ```
**Note: the order of flags does not matter but all flags need to follow main command and before any file path / directory path arguments**
