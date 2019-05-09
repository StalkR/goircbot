package url

import (
  "testing"
)

func TestTrim(t *testing.T) {
  for _, tt := range []struct {
    text string
    want string
  }{
    {text: "\n\t\r  foo  \n\t\r\n\t\r  bar   \n\t\r", want: "foo bar"},
  } {
    if got := trim(tt.text); got != tt.want {
      t.Errorf("trim(%#v): got %v; want %v", tt.text, got, tt.want)
    }
  }
}

func TestStripTags(t *testing.T) {
  for _, tt := range []struct {
    text string
    want string
  }{
    {text: "<a href='xxx'>foo <b>bar</b></a>", want: "foo bar"},
  } {
    if got := stripTags(tt.text); got != tt.want {
      t.Errorf("stripTags(%#v): got %v; want %v", tt.text, got, tt.want)
    }
  }
}
