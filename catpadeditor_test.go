package main

import (
  "testing"
  "encoding/hex"
  "encoding/base64"
  "strconv"
  "strings"
  "math/rand"
  "bytes"
  "fmt"
)

func TestFormatText(t *testing.T) {
  text := formatText("\nteste")
  if text != "teste\n" {
    t.Error("must have just last \\n :\n", hex.Dump([]byte(text)))
  }
}

func TestBreakLine(t *testing.T) {
  lineOne := strings.Repeat("a", 61) + "\n"
  lineTwo := strings.Repeat("b", 60) + "\n"
  lineThree := strings.Repeat("c", 25) + "\nx"

  bookExpected := strings.Repeat("a", 60) + "\na" +
                  strings.Repeat("b", 59) + "\nb" +
                  strings.Repeat("c", 25) + "\nx"

  book := lineOne + lineTwo + lineThree
  newBook := UpdateLineBreaks(book)

  fmt.Println("newBook")
  fmt.Println(newBook)
  fmt.Println("")

  fmt.Printf("%x\n\n", newBook)

  fmt.Println("bookExpected")
  fmt.Println(bookExpected)
  fmt.Println("")
  fmt.Printf("%x\n\n", bookExpected)

  if bookExpected != newBook {
    t.Error("book not updated correctly")
  }
}

func TestCharacterInsert(t *testing.T) {
  ce := NewCatpadEditor(10, 90, 1, "Tete\n")
  ce.InsertCharacter("s", 2)
  ce.InsertCharacter("s", 5)

  book := ce.CurrentBook()
  if book != "Testes\n" {
    t.Error("Book should be [Testes\\n]")
  }
}

func TestCharacterInsert2(t *testing.T) {
  inicial := "The Project Gutenberg EBook of Plain Tales from the Hills, b\ny Rudyard Kipling\nxx"
  expected := "aThe Project Gutenberg EBook of Plain Tales from the Hills, \nby Rudyard Kipling\nxx"
  ce := NewCatpadEditor(10, 90, 1, inicial)
  ce.InsertCharacter("a", 0)

  book := ce.CurrentBook()
  if book != expected {
    t.Error("Book should be \n" + expected + " but was \n" + book)
  }
}


func TestCharacterInsertWithLineBreaker(t *testing.T) {
  lineText := randStr(60)
  ce := NewCatpadEditor(10, 90, 1, lineText)
  ce.InsertCharacter("x", 60)

  book := ce.CurrentBook()
  if book != lineText+"\nx" {
    t.Error("should return line came " + lineText)
  }
}

func BenchmarkCharacterInsert(b *testing.B) {
    str := randStr(70000*40)
    ce := NewCatpadEditor(10, 90, 1, str)

    for n:=0; n<b.N; n++ {
      ce.InsertCharacter("a", 100)
    }
}

func randStr(len int) string {
    buff := make([]byte, len)
    rand.Read(buff)
    str := base64.StdEncoding.EncodeToString(buff)
    // Base 64 can be longer than len
    return str[:len]
}

func randBook(rows int) string {
    buff := bytes.NewBufferString("")
    for i:=0; i<rows; i++ {
      buff.WriteString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\n")
    }

    return buff.String()
}

func TestCharacterDelete(t *testing.T) {
  ce := NewCatpadEditor(10, 90, 1, "Tesite\n")
  ce.DeleteCharacter(3)

  book := ce.CurrentBook()
  if book != "Teste\n" {
    t.Error("Book should be [Teste\\n] but is " + book)
  }
}

func TestCharactersDelete(t *testing.T) {
  ce := NewCatpadEditor(10, 90, 1, "Texxxxxste\n")
  ce.DeleteCharacters(2, 6)

  book := ce.CurrentBook()
  if book != "Teste\n" {
    t.Error("Book should be [Teste\\n] but is " + book)
  }
}

func TestBookPosition(t *testing.T) {
  text := "The\nquick\nbrown\nfox\njumps\nover\nthe\nlazy\ndog"
  ce := NewCatpadEditor(10, 90, 1, text)

  row := 4
  col := 1

  ap := ce.BookPosition(row, col)

  if ap != 21 {
    t.Error("should return position 21 from 4,0 but was " + strconv.Itoa(ap))
  }
}

func TestFindCharacter(t *testing.T) {
  text := "The\nquick\nbrown\nfox\njumps\nover\nthe\nlazy\ndog"
  ce := NewCatpadEditor(10, 90, 1, text)
  pos := ce.FindPrevious(22, "\n") //j[u]mps,

  if pos != 19 {
    t.Error("Wrong position of \\n should be 19 and was " + strconv.Itoa(pos))
  }
}
