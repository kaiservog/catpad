package main

import (
    "bufio"
    "fmt"
    "os"
    "log"
    "os/exec"
    "bytes"
    "strconv"
    //"strings"
    //"time"
//    "unicode/utf8",
    //"encoding/hex"
  )

//prints
const CURSOR_MOVE =  "\x1b[999C\x1b[999B"
const CURSOR_GET_POSITION = "%c[6n"
const SCREEN_CLEAR = "\x1b[2J"
const SCREEN_CLEAN_LINE = "\x1b[K"
const CURSOR_HIDE = "\x1b[?25l"
const CURSOR_SHOW = "\x1b[?25h"
const CURSOR_POS_START = "\x1b[H"
const CTRL_Q = 0x11
const SPECIAL = 0x1b

const ARROW_UP = "A"
const ARROW_DOWN = "B"
const ARROW_LEFT = "D"
const ARROW_RIGHT = "C"
const PAGE_UP = "5"
const PAGE_DOWN = "6"
const HOME = "1"
const END = "7"

const TITLE = "Catpad (0.0.0)"
const GOODBYE_MSG = "bye!"
const MAX_COLUMN = 80

type key int
func (k key) String() string {
  return strconv.Itoa(int(k))
}

const (
  TEXT_KEY key = iota
  ARROW_LEFT_KEY
  ARROW_RIGHT_KEY
  ARROW_UP_KEY
  ARROW_DOWN_KEY
  HOME_KEY
  END_KEY
  PAGE_UP_KEY
  PAGE_DOWN_KEY
  DELETE_KEY
  BACKSPACE_KEY
  ENTER_KEY
  ESC_KEY
  CTRL_Q_KEY
)

func fillKeyMap(keymap map[string]key) {
  keymap["\x1b[A\x00"] = ARROW_UP_KEY
  keymap["\x1b[B\x00"] = ARROW_DOWN_KEY
  keymap["\x1b[C\x00"] = ARROW_RIGHT_KEY
  keymap["\x1b[D\x00"] = ARROW_LEFT_KEY
  keymap["\x1b[1~"] =    HOME_KEY
  keymap["\x1b[7~"] =    HOME_KEY
  keymap["\x1bOH\x00"] = HOME_KEY
  keymap["\x1b[H\x00"] =    HOME_KEY
  keymap["\x1b[4~"] =    END_KEY
  keymap["\x1b[8~"] =    END_KEY
  keymap["\x1bOF\x00"] = END_KEY
  keymap["\x1b[F\x00"] =    END_KEY
  keymap["\x1b[5~"] =    PAGE_UP_KEY
  keymap["\x1b[6~"] =    PAGE_DOWN_KEY
  keymap["\x1b[3~"] =    DELETE_KEY
  keymap["\x08\x00\x00\x00"] =    BACKSPACE_KEY
  keymap["\x7f\x00\x00\x00"] =    BACKSPACE_KEY
  //keymap["\n\x00\x00\x00"] =    ENTER_KEY
  //keymap["\r\n\x00\x00"] =    ENTER_KEY
  //keymap["\r\x00\x00\x00"] =    ENTER_KEY
  keymap["\x11\x00\x00\x00"] = CTRL_Q_KEY
  keymap["\x1b\x00\x00\x00"] = ESC_KEY
}

var originalSttyState bytes.Buffer

func main() {
    err := getSttyState(&originalSttyState)
  	if err != nil {
  		log.Fatal(err)
  	}
    defer setSttyState(&originalSttyState)

    configTerminalAsRaw()
    ce := setupCatpad()
    keymap := make(map[string]key)
    fillKeyMap(keymap)

    var b []byte = make([]byte, 4)
    var c []byte = make([]byte, 4)
    screenBuffer := bytes.NewBufferString("")

    for {
        editorRefreshScreen(&ce, screenBuffer)

        _, err := os.Stdin.Read(b)
        if err != nil {
          QuitMessage(err.Error())
        }

        bs := string(b)
        processed := processCtrlKeypress(bs, &ce, keymap)

        if !processed {
            firstChar := string(b[0])
            if firstChar == "\x11" || firstChar == "\x1b"{
              //do nothing
            } else {
              pos := ce.BookPosition(ce.GetCurrentRowPos(),
                ce.GetCurrentEditorColPos())
              ce.InsertCharacter(firstChar, pos)
              moveCursor(ARROW_RIGHT, &ce)
            }
        }
        copy(b, c)                             //clean input
        screenBuffer.Reset()
    }
}


func interpretKey(bs string, keymap map[string]key) key {
  k := keymap[bs]
  if k == 0 {
    return TEXT_KEY
  }
  return k
}

func backspaceKey(ce *CatpadEditor) {
  if ce.Cx == 0 {
    return
  }

  pos := ce.BookPosition(ce.GetCurrentRowPos(), ce.GetCurrentEditorColPos())
  ce.DeleteCharacter(pos)
  ce.SetCursorPosition(ce.GetCurrentEditorRowPos()-1, ce.GetCurrentEditorColPos()-1)
}

func processCtrlKeypress(bs string, ce *CatpadEditor,
  keymap map[string]key) bool {

  k := interpretKey(bs, keymap)
  //quitMessage(originalSttyState, "teste " + k.String() + "--" + string(b[1:3]))
  switch (k) {
    case TEXT_KEY:
      return false
    case ARROW_UP_KEY, ARROW_DOWN_KEY, ARROW_LEFT_KEY, ARROW_RIGHT_KEY:
        moveCursor(string(bs[2]), ce)
        break
    case HOME_KEY:
        ce.SetCursorPosition(ce.GetCurrentEditorRowPos(), 0)
        break
    case END_KEY:
        endKey(ce)
        break
    case PAGE_UP_KEY, PAGE_DOWN_KEY:
        movePage(string(bs[2]), ce)
        break
    case DELETE_KEY:
        colPos := ce.GetCurrentEditorColPos()
        deleteChars(ce, colPos+1, colPos+1, ce.GetCurrentEditorRowPos())
      break
    case BACKSPACE_KEY:
        backspaceKey(ce)
      break
    case ESC_KEY, CTRL_Q_KEY:
      ce.SetCursorPosition(0, 0)
      QuitMessage(GOODBYE_MSG)
  }
  return true
}

func endKey(ce *CatpadEditor) {
  pos := ce.BookPosition(ce.GetCurrentRowPos(),
    ce.GetCurrentEditorColPos())

  for i:=0; i < ce.BookLen()-pos; i++ {
    if ce.GetCharacter(pos+i) != "\n" {
      moveCursor(ARROW_RIGHT, ce);
    } else {
      break
    }
  }
}

func deleteChars(ce *CatpadEditor, colBegin, colEnd, row int) {
  if colBegin < 1 || colEnd > MAX_COLUMN {
    return
  }
  pos := ce.BookPosition(ce.GetCurrentRowPos(), ce.GetCurrentEditorColPos())

  ce.DeleteCharacters(colBegin+pos, colEnd+pos)
  ce.SetCursorPosition(ce.GetCurrentEditorRowPos()-1, ce.GetCurrentEditorColPos()-1)
}

func quit() {
  fmt.Println(SCREEN_CLEAR)
  fmt.Println(GOODBYE_MSG)
  setSttyState(&originalSttyState)
  os.Exit(0)
}

func QuitMessage(msg string) {
  fmt.Println(SCREEN_CLEAR)
  setSttyState(&originalSttyState)
  fmt.Println(msg)
  os.Exit(0)
}

func isControlKey(b []byte) bool {
  if b[0] == 0x11 || b[0] == 0x1B {
    return true;
  } else {
    return false;
  }
}

func updateCursor(x, y int, sb *bytes.Buffer) {
  sb.WriteString(fmt.Sprintf("\x1b[%d;%dH", y+1, x+1))
}

func editorRefreshScreen(ce *CatpadEditor, screenBuffer *bytes.Buffer) {
  //clean screen (http://vt100.net/docs/vt100-ug/chapter3.html#ED)
  //fmt.Println(SCREEN_CLEAR)
  //screenBuffer.WriteString(SCREEN_CLEAR) optimal line by line
  editorDrawRows(ce, screenBuffer)

  updateCursor(ce.Cx, ce.Cy, screenBuffer)

  draw(screenBuffer)
}

func getTitle(catpadEditor *CatpadEditor) string {
  pad := (*catpadEditor.Cols /2) - (len(TITLE) /2)
  titlePadded := fmt.Sprintf("%" +strconv.Itoa(pad) +"v", TITLE)
  return titlePadded
}

func getRowDetails(currentRow int, isLineInFocus bool) string {
  details := fmt.Sprintf("% 10d", currentRow + 1)

  if isLineInFocus {
    return "\033[1m\033[4m"+details+"\033[0m\033[0m|"
  } else {
    return details + "|"
  }
}

func editorDrawRows(catpadEditor *CatpadEditor, sb *bytes.Buffer) {
  title := getTitle(catpadEditor)

  sb.WriteString(CURSOR_HIDE + CURSOR_POS_START + title+"\r\n" + SCREEN_CLEAN_LINE)

  book := catpadEditor.Books[catpadEditor.Book]
  currentRowPos := catpadEditor.GetCurrentEditorRowPos()
  currentRow := 0

  sb.WriteString(SCREEN_CLEAN_LINE + "\r\n")
  details := getRowDetails(currentRow+catpadEditor.RowsRolled, currentRow == currentRowPos)
  sb.WriteString(details)
  currentRow++
  bookLen := catpadEditor.BookLen()

  rollRow := 0
  for i:=0; i < bookLen; i++ {
    c := string(book[i])

    if rollRow < catpadEditor.RowsRolled {
      if c == "\n" || c == "\r" {
        rollRow++
      }

      continue
    }

    if c == "\n" || c == "\r" {
      if rollRow == catpadEditor.RowsRolled {
        details := getRowDetails(currentRow+rollRow, currentRow == currentRowPos)

        sb.WriteString(SCREEN_CLEAN_LINE + "\r\n" + details)
        currentRow++

        if currentRow >= *catpadEditor.Rows-3 {
          break
        }

      } else {
        rollRow++
        continue
      }

    } else {
      sb.WriteString(c)
    }
  }

  sb.WriteString("\r\n" + CURSOR_SHOW)
}

func draw(screenBuffer *bytes.Buffer) {
  fmt.Print(screenBuffer)
}

func configTerminalAsRaw() {
  setSttyState(bytes.NewBufferString("-ignbrk"))
  setSttyState(bytes.NewBufferString("-brkint"))
  setSttyState(bytes.NewBufferString("-ignpar"))
  setSttyState(bytes.NewBufferString("-parmrk"))
  setSttyState(bytes.NewBufferString("-inpck"))
  setSttyState(bytes.NewBufferString("-istrip"))
  setSttyState(bytes.NewBufferString("-inlcr"))
  setSttyState(bytes.NewBufferString("-igncr"))
  setSttyState(bytes.NewBufferString("-icrnl"))
  setSttyState(bytes.NewBufferString("-ixon"))
  setSttyState(bytes.NewBufferString("-ixoff"))
  setSttyState(bytes.NewBufferString("-iuclc"))
  setSttyState(bytes.NewBufferString("-ixany"))
  setSttyState(bytes.NewBufferString("-imaxbel"))
  setSttyState(bytes.NewBufferString("-opost"))
  setSttyState(bytes.NewBufferString("-echo"))
  setSttyState(bytes.NewBufferString("-isig"))//not working
  setSttyState(bytes.NewBufferString("cbreak"))

  //exec.Command("stty", "-isig").Run()
  exec.Command("stty", "-icanon -xcase min 1 time 1").Run()
}

func setSttyState(state *bytes.Buffer) (err error) {
  var stderr bytes.Buffer
	cmd := exec.Command("stty", state.String())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
  cmd.Stderr = &stderr
  err = cmd.Run()
  if err != nil {
    //fmt.Println(err, string(stderr.Bytes()))
  }

	return err
}

func getSttyState(state *bytes.Buffer) (err error) {
	cmd := exec.Command("stty", "-g")
	cmd.Stdin = os.Stdin
	cmd.Stdout = state
	return cmd.Run()
}

func setupCatpad() CatpadEditor {
  rows, cols := getWindowSize()
  setUpCatPadInDisk()
  inicialContent := GetTxtContent()

  ce := NewCatpadEditor(rows, cols, 1, inicialContent)
  return ce
  /** TODO
  fileUpdateTimer := time.NewTimer(5 * time.Second)
  go func() {
      <-fileUpdateTimer.C
      HasToUpdateFile(catpadEditor)
  }()
  */
}

func movePage(direction string, ce *CatpadEditor) {
  times := *ce.Rows - 4                   //header +1 and footer +1 terminal +2
  switch (direction) {
    case PAGE_UP:
      for i:= 0; i < times; i++ {
          moveCursor(ARROW_UP, ce)
      }
      break
    case PAGE_DOWN:
      for i:= 0; i < times; i++ {
          moveCursor(ARROW_DOWN, ce)
      }
      break
  }
}

func defineCursorPos(ce *CatpadEditor, x, y int) {
  ce.Cx = x
  ce.Cy = y
}

func moveCursor(direction string, ce *CatpadEditor) {
  switch (direction) {
    case ARROW_LEFT:
      if ce.Cx -1 == 10 {
          if ce.GetCurrentEditorRowPos() > 0 {
            pos := ce.BookPosition(ce.GetCurrentRowPos(),
              ce.GetCurrentEditorColPos())
            pos = ce.FindPrevious(pos, "\n")
            _, col := ce.ScreenPosition(pos)

            //QuitMessage("pos: " + strconv.Itoa(col))
            ce.SetCursorPosition(ce.GetCurrentEditorRowPos()-1, col)
          }
      } else {
        ce.Cx--;
      }
      break
    case ARROW_RIGHT:
      pos := ce.BookPosition(ce.GetCurrentRowPos(), ce.GetCurrentEditorColPos())
      c := ce.GetCharacter(pos)

      if c == "\n" {
        ce.SetCursorPosition(ce.GetCurrentEditorRowPos() +1, 0)
        break
      }

      if ce.Cx != *ce.Cols -1 {
        ce.Cx++;
      }

      break
    case ARROW_UP:
      if ce.Cy != 2 {
        ce.Cy--;
      } else {
        if ce.RowsRolled > 0 {
          ce.RowsRolled--
        }
      }

      row := ce.CurrentBookLines()[ce.RowsRolled + ce.GetCurrentEditorRowPos()]
      rowSize := len(row)
      if rowSize < ce.GetCurrentEditorColPos() {
        ce.SetCursorPosition(ce.GetCurrentEditorRowPos(), rowSize)
      }

      break
    case ARROW_DOWN:
      if ce.Cy != *ce.Rows -2 {
        if ce.Cy < len(ce.CurrentBookLines()) {
            ce.Cy++;
        }
      } else {
        if ce.RowsRolled < len(ce.CurrentBookLines()) {
          ce.RowsRolled++
        }
      }

      row := ce.CurrentBookLines()[ce.RowsRolled + ce.GetCurrentEditorRowPos()]
      rowSize := len(row)
      if rowSize < ce.GetCurrentEditorColPos() {
        ce.SetCursorPosition(ce.GetCurrentEditorRowPos(), rowSize)
      }

      break
  }
}

func getWindowSize() (int, int) {
  fmt.Print(CURSOR_MOVE)
  return getCursorPosition()
}

func getCursorPosition() (int, int) {
  cmd := exec.Command("echo", fmt.Sprintf(CURSOR_GET_POSITION, 27))
  buff := &bytes.Buffer{}
  cmd.Stdout = buff

  cmd.Run()
  reader := bufio.NewReader(os.Stdin)

  fmt.Print(buff)
  text, _ := reader.ReadSlice('R')
  r := string(text)[2:]

  var rows, cols int
  fmt.Sscanf(r, "%d;%d", &rows, &cols)

  return rows, cols
}
