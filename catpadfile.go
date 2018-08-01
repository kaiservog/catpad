package main

import (
   "os"
   "os/user"
   "io/ioutil"
   "path/filepath"
   "bufio"
   "encoding/json"
)

const BASE_SW_NAME = "catpad"
const FILE_TEXT_NAME = "catpad.txt"
const FILE_DATA_NAME = ".data.json"
const FILE_CONTENT =
`
Welcome to CatPad (Beta)

Don't erase this text yet!!, it will help you to get use
to this notepad, when you become confortable with the shortcuts
be free to delete this text

be free to edit this text as well to make you own annotation about catpad

ctrl q   ->  Close this application
`
//TODO criar marcador para "cat tips"

func GetTxtContent() string {
  usr, _ := user.Current()
  path := filepath.Join(usr.HomeDir, "." + BASE_SW_NAME, FILE_TEXT_NAME)
  content, _ := ioutil.ReadFile(path)
  return string(content)
}

func setUpCatPadInDisk() error {
  usr, err := user.Current()

  if err != nil {
    return err
  }

  listDir, _ := ioutil.ReadDir(usr.HomeDir)

  hasDir := false
  for i:=0; i < len(listDir); i++ {
      if listDir[i].Name() == "."+BASE_SW_NAME {
        hasDir = true
        break
      }
  }

  catpadPath := filepath.Join(usr.HomeDir, "." + BASE_SW_NAME)
  if !hasDir {
    os.MkdirAll(catpadPath, os.ModePerm);
  }

  listCatPadDir, _ := ioutil.ReadDir(catpadPath)
  hasFile := false
  for i:=0; i < len(listCatPadDir); i++ {
      if listCatPadDir[i].Name() == FILE_TEXT_NAME {
        hasFile = true
        break
      }
  }

  if !hasFile {
    txtFilePath := filepath.Join(catpadPath, FILE_TEXT_NAME)
    f, err := os.Create(txtFilePath)
    if err != nil {
        return err
    }

    defer f.Close()

    w := bufio.NewWriter(f)
    _, err = w.WriteString(FILE_CONTENT)
    if err != nil {
        return err
    }

    w.Flush()
  }

  //.data.json
  hasFileData := false
  for i:=0; i < len(listCatPadDir); i++ {
      if listCatPadDir[i].Name() == FILE_DATA_NAME {
        hasFileData = true
        break
      }
  }

  if !hasFileData {
    dataFilePath := filepath.Join(catpadPath, FILE_DATA_NAME)
    f, err := os.Create(dataFilePath)
    if err != nil {
        return err
    }

    defer f.Close()

    w := bufio.NewWriter(f)
    _, err = w.WriteString("topline:100")
    if err != nil {
        return err
    }

    w.Flush()
  }


  return nil
}

func GetCatpadData() *CatpadData {
  usr, _ := user.Current()
  path := filepath.Join(usr.HomeDir, "." + BASE_SW_NAME, FILE_DATA_NAME)
  file, err := ioutil.ReadFile(path)

  if err != nil {
    os.Exit(1)
  }

  var catpadData CatpadData
  json.Unmarshal(file, &catpadData)

  return &catpadData
}

func HasToUpdateFile(ce *CatpadEditor) {
  //var lineBytes = make([]byte, 0)

  if ce.Fresh {
    updateFile(*ce)
  }
}

func updateFile(ce CatpadEditor) {
  usr, _ := user.Current()
  catpadPath := filepath.Join(usr.HomeDir, "." + BASE_SW_NAME)
  txtFilePath := filepath.Join(catpadPath, FILE_TEXT_NAME + ".tmp")
  f, err := os.Create(txtFilePath)
  if err != nil {
      os.Exit(1)
  }

  defer f.Close()

  w := bufio.NewWriter(f)

  _, err = w.WriteString(ce.Books[ce.Book])
  if err != nil {
    os.Exit(1)
  }

  w.Flush()

}

type CatpadData struct {
    Topline int
}
