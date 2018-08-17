package main

import (
  "strings"
//  "fmt"
  //"strconv"
)

const LB = "\n"

func formatText(nt string) string {
  if len(nt) == 0 {
    return LB
  } else {
    tr := strings.Replace(nt, "\n", "", -1)
    tr = tr + LB
    return tr
  }
}

type CatpadEditor struct {
    //screen
    Rows *int
    Cols *int

    //editor
    Cx, Cy int

    //file
    RowsRolled int //first row in editor showing
    FirstBookRow int

    Book    int
    bookLen int
    Books []string
    CatpadData *CatpadData
    Fresh bool

    //bookmark
    Bookmarks []Bookmark
}

func NewCatpadEditor(rows, cols, firstBookRow int, inicialContent string) CatpadEditor {
  ce := CatpadEditor{
    Rows: &rows,
    Cols: &cols,

    Cx: 11, //row details size
    Cy: 1+1,//title + empty row

    RowsRolled: 0,
    FirstBookRow: firstBookRow,

    Book: 0,
    Fresh: false }

  ce.Books = make([]string, 1)
  //inicialContent = UpdateLineBreaks(inicialContent)
  ce.Books[0] = inicialContent
  ce.bookLen = len([]rune(ce.CurrentBook()))

  return ce
}

func (ce *CatpadEditor) GetCurrentEditorRowPos() int {
  return ce.Cy -2
}

func (ce *CatpadEditor) GetCurrentEditorColPos() int {
  return ce.Cx -11
}

func (ce *CatpadEditor) GetCurrentRowPos() int {
  return ce.GetCurrentEditorRowPos() + ce.RowsRolled
}

func (ce *CatpadEditor) SetCursorPosition(row, col int) {
  ce.Cx = col+11
  ce.Cy = row+1+1
}

func (ce *CatpadEditor) GetCharacter(pos int) string {
  return string(ce.CurrentBook()[pos])
}

func (ce *CatpadEditor) VersionBook(newBook string) {
  ce.Book++
  ce.Books = append(ce.Books, newBook)
  ce.bookLen = len([]rune(ce.CurrentBook()))
}

func (ce *CatpadEditor) InsertCharacter(char string, pos int) {
  currentBook := ce.CurrentBook()
  leftBook := currentBook[:pos]
  rightBook := currentBook[pos:]

  newBook :=  leftBook + char + rightBook
  newBook = UpdateLineBreaks(newBook)
  ce.VersionBook(newBook)
}

func overflow(text string, size int) string {
  if len(text) <= size {
    return ""
  }

  return text[size:]
}

func UpdateLineBreaks(book string) string {
  books := strings.Split(book, "\n")

  for i:=0; i < len(books); i++ {
      of := overflow(books[i], 60)
      ofLen := len(of)

      if ofLen > 0 {
        books[i] = books[i][:60]

        if len(books) > i+1 {
          books[i+1] = of + books[i+1]
        } else {
          books = append(books, of)
        }
      }
  }

  return strings.Join(books, "\n")
}

func (ce *CatpadEditor) DeleteCharacters(begin, end int) {
  currentBook := ce.CurrentBook()
  if end > len(currentBook) {
    return
  }

  leftBook := currentBook[:begin]
  rightBook := currentBook[end+1:]

  newBook :=  leftBook + rightBook
  newBook = UpdateLineBreaks(newBook)

  ce.VersionBook(newBook)
}

func (ce *CatpadEditor) DeleteCharacter(pos int) {

  currentBook := ce.CurrentBook()
  leftBook := currentBook[:pos]
  rightBook := currentBook[pos+1:]

  newBook :=  leftBook + rightBook

  ce.Book++
  ce.Books = append(ce.Books, newBook)
  ce.bookLen = len([]rune(ce.CurrentBook()))
}

func (ce *CatpadEditor) CurrentBook() string {
  return ce.Books[ce.Book]
}

func (ce *CatpadEditor) CurrentBookLines() []string {
  return strings.Split(ce.CurrentBook(), "\n")
}

func (ce *CatpadEditor) FindPrevious(pos int, c string) int {
  bookTextSearch := ce.CurrentBook()[:pos]
  size := len([]rune(bookTextSearch))

  for i:=size; i>0; i-- {
    if string(bookTextSearch[i-1]) == c {
      return i-1
    }
  }
  return -1
}

func (ce *CatpadEditor) BookLen() int {
    return ce.bookLen
}

func (ce *CatpadEditor) ScreenPosition(pos int) (row int, col int) {
  for i:=0; i < ce.BookLen(); i++ {
    c := string(ce.CurrentBook()[i])
    if c == "\n" {
      // just jump row when its not the breakline that caller wants
      if i == pos {
        return row, col
      }
      row++
      col=0
    } else {
      col++
    }
  }
  return row, col
}

func (ce *CatpadEditor) BookPosition(row, col int) int {
  if row == 0  && col == 0 {
    return 0
  }

  currentBook := ce.CurrentBook()
  rowsPassed := 0
  colsPassed := 0

  for i:=0; i < ce.BookLen(); i++ {
    if rowsPassed == row && colsPassed == col {
      return i
    }

    c := string(currentBook[i])
    if c == "\n" {
      rowsPassed++
      colsPassed=0
    } else {
      colsPassed++
    }
  }

  return -1
}
