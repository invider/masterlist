package main

import (
    "os"
    "os/user"
    "flag"
    "fmt"
    "log"
    "strings"
	"main/ml"
)

func main() {

    // process options
    var pFlag = flag.Bool("p", false, "(plain) hide line numbers and file tags")
    var vFlag = flag.Bool("v", false, "verbose output")
    var iFlag = flag.Bool("i", false, "show invisible scanlines")
    var nFlag = flag.Int("n", 0, "limit number of scanlines shown")
    var bFlag = flag.Int("b", 0, "start showing from number")
    var eFlag = flag.Int("e", 0, "end showing at number")
    var tFlag = flag.Int("t", 0, "target for command")
    var fFlag = flag.Bool("f", false, "force command")
    var gFlag = flag.Bool("g", false, "work with global $HOME and $MLPATH .ml")
    var vvFlag = flag.Bool("V", false, "show fossil scanlines (V/X ...)")
    var usageFunc = flag.Usage
    flag.Usage = func() {
        usageFunc()
        fmt.Println(
`
=== Available Commands ===
    a, add - add new line on top
    v, done, fossil - mark as done (V)
    x, cancel, kill - mark as canceled (X)
    d, delete, remove - remove the line completely
    ., score, cycle - score some cycles for the line/tag
`)
        fmt.Println("\nType ml h or ml help for more information\n")
    }
    flag.Parse()

    ml.Ctx.P = *pFlag
    ml.Ctx.V = *vFlag
    ml.Ctx.I = *iFlag
    ml.Ctx.N = *nFlag
    ml.Ctx.B = *bFlag
    ml.Ctx.E = *eFlag
    ml.Ctx.T = *tFlag
    ml.Ctx.F = *fFlag
    ml.Ctx.G = *gFlag
    ml.Ctx.VV = *vvFlag

    if (ml.Ctx.V) {
        fmt.Println("Welcome to MasterList v0.1")
    }

    args := flag.Args()
    for i := 0; i < len(args); i++ {
        switch arg := strings.ToLower(args[i]); {
        case arg == "h" || arg == "help":
            flag.Usage()
            return
        case arg == "l" || arg == "list":
            ml.Ctx.Cmd = "l"
        case arg == "a" || arg == "add":
            ml.Ctx.Cmd = "a"
            i++
            ml.Ctx.Reg = args[i]
        case arg == "v" || arg == "done" || arg == "fossil":
            ml.Ctx.Cmd = "v"
        case arg == "x" || arg == "cancel" || arg == "kill":
            ml.Ctx.Cmd = "x"
        case arg == "d" || arg == "delete" || arg == "remove":
            ml.Ctx.Cmd = "d"
        case arg == "." || arg == "s" || arg == "score" || arg == "cycle":
            ml.Ctx.Cmd = "s"
        default:
            ml.Ctx.Filter = append(ml.Ctx.Filter, args[i])
        }
    }

    // load local
    localMode := false
    globalMode := false
    localPath, err := os.Getwd()
    if err != nil {
        log.Println(err)
    } else {
        if (ml.Load(localPath) > 0) {
            localMode = true
        }
    }
    // load home
    if ml.Ctx.G || !localMode {
        globalMode = true
        usr, err := user.Current()
        if err != nil {
            log.Println(err)
        } else {
            ml.Load(usr.HomeDir)
        }
    }
    // load env
    if ml.Ctx.G || !localMode {
        globalMode = true
        path := os.Getenv("MLPATH")
        if ml.Ctx.V {
            fmt.Println("=== MLPATH: " + path)
        }
        if path != "" {
            ml.Load(path)
        }
    }

    if ml.Ctx.V {
        if localMode {
            fmt.Println("=== Working with local context ===")
        }
        if globalMode {
            fmt.Println("=== Working with global context ===")
        }
    }

    // execute context
    ml.Process()
}
