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

## Testing
### To run tests on this project. You can `cd` to the folder that includes the code you want to test and run `go test -v`. The test files should be in the same folder with the files including the code being tested. 

### To write tests
  At the directory of files including code being tests, create the new test file. The naming convention for test file is the name of the file being tested plus a suffix of `_test`. For example, if you are testing `parsers.go`, your test file should be `parsers_test.go`

### To get code coverage
  At the directory of files that you want to run code coverage on, run `go test -cover`. This will give you the percentage of codes have been covered in a package.
  