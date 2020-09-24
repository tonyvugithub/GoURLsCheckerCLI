# GO URLs CHECKER CLI

**Description**: This is a basic URL checker written in GO. The app allows users to pass in a file name as an command-line argument, then it would extract all valid URLs and run status check on those links.

## How to use? 
1. Build binary file:
   ```go
      go build linkDetector.go
   ``` 
2. Run the file:
   ```go
      ./linkDetector check [file-name]
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