// a tiny little pre-processor that reads a Sensu event JSON from
// the standard input and if the last state is new (indicates a change),
// returns the return codes 0) for a resolution, 1) for a warning, and 2) for an error,
// in addition to a summary sentence and the whole JSON structure to be included
// in the email. If the last state does not indicate a change, the return code
// will be 3.

// The usage in a sensu handler will be something like:
//    JSON="$(time2sendmail)"
//    RC=$?
//    if [ RC == 2 ]
//    then
//      echo $JSON | mail -s "Sensu Error Report" admin@website.com
//    fi
//
// This is my first Go program ever. Any feedback is more than welcome! :-)

// (c) 2015 Roozbeh Farahbod

package main
import "encoding/json"
import "fmt"
import "bufio"
import "os"

type SensuCheck struct {
  Name string
  History []string
}

type SensuEvent struct {
  Id string
  Check SensuCheck
}

func main() {
  exitCode := 0
  defer func() {
    os.Exit(exitCode)
  }()

  scanner := bufio.NewScanner(os.Stdin)
  jsonstr := ""

  for scanner.Scan() {
    jsonstr = jsonstr + "\n" + scanner.Text()
  }

  byt := []byte(jsonstr)
  var event SensuEvent

  if err := json.Unmarshal(byt, &event); err != nil {
     panic(err)
  }

  var history = event.Check.History
  var hlen = len(history)

  var time2sendOnError = (hlen > 0 && history[hlen-1] == "2") && (hlen == 1 || (history[hlen-1] != history[hlen-2]))
  var time2sendOnWarn = (hlen > 0 && history[hlen-1] == "1") && (hlen == 1 || (history[hlen-1] != history[hlen-2]))
  var time2sendOnRecovery = hlen > 1 && (history[hlen-1] == "0") && (history[hlen-2] != "0")

  if time2sendOnError {
    fmt.Println("Check '" + event.Check.Name + "' failed.")
    fmt.Println()
    fmt.Println(jsonstr)
    exitCode = 2
  } else {
    if time2sendOnWarn {
      fmt.Println("Check '" + event.Check.Name + "' has warning.")
      fmt.Println()
      fmt.Println(jsonstr)
      exitCode = 1
    } else {
      if time2sendOnRecovery {
        fmt.Println("Check '" + event.Check.Name + "' is resolved.")
        fmt.Println()
        fmt.Println(jsonstr)
        exitCode = 0
      } else {
        exitCode = 3
      }
    }
  }

}
