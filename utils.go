package main

import (
	"fmt"
	"log"
	"strings"
	"time"
  "strconv"
)

func (app *application) printSumarry() {
  fmt.Println("-------------------------------")
  fmt.Printf("Duplicated records: %d\n", app.DupRecords)
  fmt.Printf("New records: %d\n", app.NewRecords)
  fmt.Println("-------------------------------")
}

func check(err error) {
  if err != nil {
    log.Fatal(err)
  }
}

// extractIdAndCatigorie extract the catigorieId and id form a morocannonces post link
//
// the return: cat, id, error
func extractIdAndCatigorie(url string) (int, int, error) {
  base1 := 10
  if base1 > len(url) {
    return 0, 0, fmt.Errorf("expected a valid url found: %s", url)
  }
  var id []byte
  var cat []byte
  for i := base1; i < len(url); i++ {
    if url[i] >= byte('0') && url[i] <= byte('9') {
      cat = append(cat, url[i])
    } else if url[i] == '/' {
      break;
    } else {
      fmt.Printf("error in the url: %s\n", url)
      return 0, 0, fmt.Errorf("expected / found: %d", url[i])
    }
  }
  base2 := base1 + len(cat) + 23

  for i := base2; i < len(url); i++ {
    if url[i] >= '0' && url[i] <= '9' {
      id = append(id, url[i])
    } else if url[i] == '/' {
      break;
    } else {
      fmt.Printf("error in the url: %s\n", url)
      return 0, 0, fmt.Errorf("expected / found: %d", url[i])
    }
  }

  return toInt(cat),toInt(id), nil
}

// toInt convert arr to an integer value
func toInt(arr []byte) int {
  multiplier := 1
  result := 0
  for i := 0; i < len(arr); i++ {
    multiplier *= 10
  }
  for _, v := range arr {
    multiplier /= 10
    result += (int(v) - 48) * multiplier
  }
  return result
}

var monthsArr [12]string = [12]string{"Jan", "Feb", "Mar", "Apr", "May", "June", "July", "Aug", "Sept", "Oct", "Nov", "Dec"}

// extractDateAndTime extract time and date from s, example: "Publiée le: 4 Mar-9:58" -> "04/03/2023", "09:58:00"
// should only be used in collectPage
func extractDateAndTime(s string) (string, string) {
  parts := strings.Split(strings.ReplaceAll(s, "Publiée le: ", ""), "-")
  day := strings.Split(parts[0], " ")[0]
  strMonth := strings.Split(parts[0], " ")[1]

  var intMonth int
  for i, v := range monthsArr {
    if v == strMonth {
      intMonth = i + 1
      break
    }
  }

  if intMonth > 9 {
    strMonth = strconv.Itoa(intMonth)
  } else {
    strMonth = "0" + strconv.Itoa(intMonth)
  }

  if len(day) == 1 {
    day = "0" + day
  }

  year := time.Now().Year()

  time_ := parts[1] + ":00"

  return time_, (strings.Join([]string{strconv.Itoa(year), strMonth, day}, "-"))
}
