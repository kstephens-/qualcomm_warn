package scraper

import (
    "testing"
    "time"
    "strings"

    "github.com/stretchr/testify/assert"
)

var new_date, _ = time.Parse("1/02/06", "10/04/16")

var result1 = QualcommEvent{date: new_date}
var result2 = QualcommEvent{date: new_date, name: "A cool birthday"}
var result3 = QualcommEvent{date: new_date, name: "A cool birthday", detail: "All day long!"}

var attributeTests = []struct {
    input string
    expected QualcommEvent
}{
    {"10/04/16", result1},
    {"A cool birthday", result2},
    {"All day long!", result3},
}

func TestSetAttributeInOrder(t *testing.T) {
    event := QualcommEvent{}
    for _, tt := range attributeTests {
        event.SetAttributeInOrder(tt.input)
        assert.Equal(t, tt.expected, event)
    }
}

var body = strings.NewReader(`
         <html>
         <body>
         <table>
         <tr>
         <td>10/04/16</td>
         <td>
             <a>A cool birthday</a>
             All day long!
        </td>
        </tr>
        </table>
        </body>
        </html>
`)

func TestExtract(t *testing.T) {
    expected := []QualcommEvent{result3}
    result := extract(body)
    assert.Equal(t, expected, result)
}
