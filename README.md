# GO URLs CHECKER CLI

**Description**: This is a basic URL checker written in GO. The app allows users to pass in a file name as a command-line argument, then it would extract all valid URLs and run status check on those urls.

## How to use? 
1. Build binary file:
   ```go
      go build linkDetector.go
   ``` 
2. Run the file:
   ```go
      ./linkDetector check -f [file-name]
   ```

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
    ./linkDetector check -f [file-name-1] [file-name-2] ... [file-name-nth]
   ```
4. Check all files in single or multiple directories:
   ```go
    ./linkDetector check -d [directory-path-1] [directory-path-2] ... [directory-path-nth]
   ```
5. Create report file by adding -r flag:
   ```go
    ./linkDetector check -f -r [file-name]
   ```
   ```go
    ./linkDetector check -d -r [directory-path]
   ```

**Note: the order of flags does not matter but all flags need to follow main command and before any file path / directory path arguments**