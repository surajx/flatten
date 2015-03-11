package main

import (
  "io"
  "os"
  "fmt"
  "flag"
  "errors"
  "strings"
  "path/filepath"
)
var verbose bool

//walk through the directories and generate a list of absolute paths of all files.
func fileLister(dirPath string, fileListPtr *[]string) error {
  fileList := *fileListPtr
  curLocation := func(path string, locInfo os.FileInfo, err error) error {
    if err!=nil {
      return err
    }
    if !locInfo.IsDir() {
      fileList = append(fileList, path)
    }
    return err
  }
  err := filepath.Walk(dirPath, curLocation)
  *fileListPtr = fileList
  return err
}

// Check pre-requesited and initiate copy
func fileCopy(src, dest string) (err error) {
  srcFileInfo, err := os.Stat(src)
  if err==nil {
    if !srcFileInfo.Mode().IsRegular() {
      return fmt.Errorf("[COPY ERROR] %s is not a regular file, mode: (%q)", srcFileInfo.Name(), srcFileInfo.Mode().String())
    }
  } else {return}
  if _, err = os.Stat(dest); err==nil {
    return fmt.Errorf("[COPY ERROR] File %s already exists at %s, copy aborted.", src, filepath.Dir(dest))
  }
  err = doFileCopy(src, dest)
  return
}

//Actually copy file content to destination folder
func doFileCopy(src, dest string) (err error) {
  inFile, err := os.Open(src)
  if err!=nil {return}
  defer inFile.Close()

  outFile, err := os.Create(dest)
  if err!=nil {return}
  defer outFile.Close()

  if _, err = io.Copy(outFile, inFile); err!=nil {return}
  err = outFile.Sync()
  if verbose && err==nil {fmt.Println("[COPIED] ", src)}
  return
}

func flatten(dirPath string) (err error) {
  fileList := []string{}
  err = fileLister(dirPath, &fileList)
  if err!=nil {
    return
  }
  for _, aFileLoc := range fileList {
    _, aFileName := filepath.Split(aFileLoc)
    if filepath.Clean(filepath.Dir(aFileLoc))==dirPath {
      continue
    }
    cerr := fileCopy(aFileLoc, filepath.Clean(dirPath + string(os.PathSeparator) + aFileName))
    if cerr!=nil{
      existingErr := ""
      if err!=nil{
        existingErr = err.Error() + "\n"
      }
      err = errors.New(existingErr + cerr.Error())
    }
  }
  return
}

func main() {

  deleteDirsPtr := flag.Bool("delete", false, "Enables deletion of directories after flatten. Use cautiously.")
  verbosePtr := flag.Bool("verbose", false, "Verbose, duh!")
  flag.Parse()
  dirPath := strings.TrimSpace(flag.Arg(0))
  verbose = *verbosePtr

  cmdLineError := func(errMsg string) {
    fmt.Println("Error: " + errMsg)
    fmt.Println("Usage: flatten [-delete] [-verbose] <Directory to flatten>")
    os.Exit(1)
  }

  if dirPath == "" || len(flag.Args())>1 {
    cmdLineError("Invalid number of arguments.")
  }

  dirPath = filepath.Clean(dirPath)

  if dirPathInfo, err:= os.Stat(dirPath); err!=nil || !dirPathInfo.IsDir() {
    cmdLineError("Directory path provided is not a directory.")
  }
  fmt.Println("Copying files to base directory... ")
  err := flatten(dirPath)
  if err!=nil {
    fmt.Println(err)
  }
  fmt.Println("Done!")
  if *deleteDirsPtr {
    if err!=nil{
      fmt.Println("\nEncountered copy errors, review errors before proceeding with delete.")
    }
    var resp string
    fmt.Printf("\nDelete all sub-directories? (yes/NO)> ")
    _, err := fmt.Scanln(&resp)
    if err==nil && strings.ToLower(resp)=="yes" {
      fmt.Println("Deleting all sub-directories...")
      dirPathFD, _ := os.Open(dirPath)
      if err != nil {
        fmt.Println("[DELETE ERROR] Unable to open ", dirPath)
        os.Exit(1)
      }
      items, _ := dirPathFD.Readdirnames(0)
      for _, item := range items {
          itemPath := filepath.Clean(dirPath + string(os.PathSeparator) + item)
        if itemPathInfo, err:= os.Stat(itemPath); err==nil && itemPathInfo.IsDir() {
          err := os.RemoveAll(itemPath)
          if err!=nil {
            fmt.Println("[DELETE ERROR] Unable to detele ", itemPath)
          }
          if verbose && err==nil {fmt.Println("[DELETED] ", itemPath)}
        }
      }
      fmt.Println("Done!")
    } else{
      fmt.Println("Skipping delete.")
    }
  }
}
