# GO URLs CHECKER CLI ENVIRONMENT SET UP

## How to use? 
1. Build binary file:
   ```go
      go build linkDetector.go
   ``` 
2. Run the file:
   ```go
      ./linkDetector check -f [file-name]
   ```

## Formatting Source 
### This project used gofmt for formatting source code and golint for linting visit `.vscode/settings.json` for details. VSCode Go extension is required for the format setting to work.  

### To format the whole project from command line
  At the root, run:
  ```go
    go fmt ./...
  ```

### To lint the whole project from command line
  At the root, run:
  ```
    golint ./...
  ```