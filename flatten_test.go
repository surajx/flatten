package main

import (
  "os"
  "fmt"
  "reflect"
  "testing"
)

func TestFileLister(t *testing.T) {
  fileList := []string{}
  expectedFileList := []string{
      "test/TestFileLister/dir1/dir2/file0",
      "test/TestFileLister/dir1/dir2/file9",
      "test/TestFileLister/dir1/dir3/abc",
      "test/TestFileLister/dir1/dir3/efg",
      "test/TestFileLister/dir1/dir3/hij",
      "test/TestFileLister/dir1/file6.py",
      "test/TestFileLister/dir1/file7.avi",
      "test/TestFileLister/dir1/file8.js",
      "test/TestFileLister/fil3.mp3",
      "test/TestFileLister/file1.py",
      "test/TestFileLister/file2.go",
      "test/TestFileLister/file4.mp4",
      "test/TestFileLister/file5.txt"}
  dirPath := "test/TestFileLister"
  err := fileLister(dirPath, &fileList)
  if err!=nil {
    t.Errorf(fmt.Sprintf("Error: %s", err))
  }
  if !reflect.DeepEqual(fileList, expectedFileList) {
    t.Errorf("Fail: FileList differs from Expected.")
  }
}

func TestFileCopy(t *testing.T){
  srcFileLoc := "test/TestFileCopy/source/sourceFile.txt"
  destFileLoc := "test/TestFileCopy/dest/sourceFile.txt"
  if _, err := os.Stat(destFileLoc); err==nil {
    os.Remove(destFileLoc)
  }
  err := fileCopy(srcFileLoc, destFileLoc)
  if err!=nil {
    t.Errorf(fmt.Sprintf("Error: %s", err))
  }
  destFileInfo, err := os.Stat(destFileLoc)
  if err!=nil {
    t.Errorf("Source File did not get copied into destination folder.")
  }
  srcFileInfo, err := os.Stat(srcFileLoc)
  if err!=nil {
    t.Errorf("Unable to stat source file!")
  }
  if destFileInfo.Size()!=srcFileInfo.Size() {
    t.Errorf("Source and Destination files are not of the same size (Dest: %d, Src: %d)", destFileInfo.Size(), srcFileInfo.Size())
  }
}
