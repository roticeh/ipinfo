package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	cInf  = color.New(color.FgCyan, color.Bold).SprintFunc()
	cWarn = color.New(color.FgYellow, color.Bold).SprintFunc()
	cErr  = color.New(color.FgRed, color.Bold).SprintFunc()
	cSucc = color.New(color.FgGreen, color.Bold).SprintFunc()
	cFatl = color.New(color.BgRed, color.FgWhite, color.Bold).SprintFunc()
	cTime = color.New(color.FgHiBlack).SprintFunc()

	cTag = color.New(color.FgHiMagenta, color.Bold).SprintFunc()
)

func init() {
	log.SetFlags(0)
}

func timeStamp() string {
	return cTime(time.Now().Format("2006-01-02 15:04"))
}

// highlightSubTag: Instantly finds and highlights the [ALT_CODE] structure at the very beginning of the message.
// Example: “[ENV_CRITICAL] Config failed” -> [ENV_CRITICAL] in purple + “Config failed” in normal text
func highlightSubTag(msg string) string {
	msg = strings.TrimSpace(msg)
	if strings.HasPrefix(msg, "[") {
		endIdx := strings.Index(msg, "]")
		if endIdx > 0 {
			tag := msg[:endIdx+1]
			rest := msg[endIdx+1:]
			return cTag(tag) + rest
		}
	}
	return msg
}

func Log(format string, v ...interface{}) {
	custom_prefix := v[0].(string)
	vobject := v[1:]
	msg := fmt.Sprintf(format, vobject...)
	msg = highlightSubTag(msg)
	fmt.Printf("%s %s %s\n", timeStamp(), cInf("["+custom_prefix+"]"), msg)
}

func LogInfo(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	msg = highlightSubTag(msg)
	fmt.Printf("%s %s %s\n", timeStamp(), cInf("[INFO]"), msg)
}

func LogSuccess(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	msg = highlightSubTag(msg)
	fmt.Printf("%s %s %s\n", timeStamp(), cSucc("[OK]"), msg)
}

func LogWarn(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	msg = highlightSubTag(msg)
	fmt.Printf("%s %s %s\n", timeStamp(), cWarn("[WARN]"), msg)
}

func LogError(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	msg = highlightSubTag(msg)
	fmt.Fprintf(os.Stderr, "%s %s %s\n", timeStamp(), cErr("[ERR]"), msg)
}

func LogFatal(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	msg = highlightSubTag(msg)
	fmt.Fprintf(os.Stderr, "%s %s %s\n", timeStamp(), cFatl("[FATAL]"), msg)
	os.Exit(1)
}

func LogServerStart(port int, baseURL string) {
	fmt.Println()
	fmt.Printf("   %s  %s\n", cSucc("⚡ Server is Active"), cTime("waiting for requests..."))
	fmt.Printf("   %s  %s\n", cInf("➜ Local:"), fmt.Sprintf("http://localhost:%d", port))
	fmt.Printf("   %s  %s\n", cInf("➜ Public:"), color.New(color.FgHiBlue, color.Underline).Sprint(baseURL))
	fmt.Println()
}
