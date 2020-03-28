package ml

import (
    "fmt"
    "log"
    "strings" 
)

var Ctx context

type context struct {
    P bool
    V bool
    I bool
    N int
    B int
    E int
    T int
    F bool
    G bool
    VV bool
    Cmd string
    Reg string
    Filter []string
}


func Process() {
    filteredLines := 0
    passedLines := 0
    curFrame := ""

    filterFn := func (l scanline, li int) bool {
            // filter non-visibles
            if !Ctx.I && l.invisible() {
                return false
            }
            // filter fossil unless -v specified
            if !Ctx.VV && (strings.HasPrefix(l.src, "V") || strings.HasPrefix(l.src, "X")) {
                return false;
            }
            // filter out by tags
            if len(Ctx.Filter) > 0 {
                match := false
                for _, ptrn := range Ctx.Filter {
                    if strings.Contains(l.src, ptrn) {
                        match = true
                    }
                    if strings.Contains(l.frm.name, ptrn) {
                        match = true
                    }
                }
                if !match {
                    return false;
                }
            }

            filteredLines++;
            // limit number of lines -n if provided
            if Ctx.N > 0 && passedLines >= Ctx.N {
                return false
            }
            // start from -b if provided
            if filteredLines < Ctx.B {
                return false
            }
            // end with -e if provided
            if Ctx.E > 0 && filteredLines > Ctx.E {
                return false;
            }
            passedLines++;
            return true
    }

    scanlineCount := 0
    countFn := func(line scanline, li int) {
        scanlineCount++
    }

    iline := 0
    listFn := func(line scanline, li int) {
            iline ++;
            if !Ctx.P && curFrame != line.frm.name {
                curFrame = line.frm.name
                fmt.Println("\n#" + curFrame)
            }
            src := line.source()
            if Ctx.T > 0 && Ctx.T == iline {
                src = src + " <---"
            }
            if !Ctx.P {
                fmt.Printf("%2d: " + src + "\n", iline)
            } else {
                fmt.Println(src)
            }
    }

    aFn := func(line scanline, li int) (*frame, *scanline) {
            iline ++;
            if iline == 1 {
                line.frm.add(scanline{frm: line.frm, src: Ctx.Reg})
                return line.frm, nil
            }
            return nil, nil
    }

    fossilCnt := 0
    fossilFn := func(line scanline, li int) (*frame, *scanline) {
        fossilCnt++
        if Ctx.T == 0 || Ctx.T == fossilCnt {
            line.fossil()
            if !Ctx.P {
                fmt.Printf("%2d: " + line.source() + "\n", fossilCnt)
            } else {
                fmt.Println(line.source())
            }
            return line.frm, &line
        }
        return nil, nil
    }

    cancelCnt := 0
    cancelFn := func(line scanline, li int) (*frame, *scanline) {
        cancelCnt++
        if Ctx.T == 0 || Ctx.T == cancelCnt {
            line.cancel()
            if !Ctx.P {
                fmt.Printf("%2d: " + line.source() + "\n", fossilCnt)
            } else {
                fmt.Println(line.source())
            }
            return line.frm, &line
        }
        return nil, nil
    }

    scoreCnt := 0
    scoreFn := func(line scanline, li int) (*frame, *scanline) {
        scoreCnt++
        if Ctx.T == 0 || Ctx.T == scoreCnt {
            line.score()
            if !Ctx.P {
                fmt.Printf("%2d: " + line.source() + "\n", fossilCnt)
            } else {
                fmt.Println(line.source())
            }
            return line.frm, &line
        }
        return nil, nil
    }

    deleteCnt := 0
    deleteFn := func(line scanline, li int) (*frame, *scanline) {
        deleteCnt++
        if Ctx.T == 0 || Ctx.T == deleteCnt {
            line.frm.remove(li)
            if !Ctx.P {
                fmt.Printf("%2d: " + line.source() + "\n", fossilCnt)
            } else {
                fmt.Println(line.source())
            }
            return line.frm, &line
        }
        return nil, nil
    }

    switch cmd := Ctx.Cmd; {
        case cmd == "a":
            Apply(filterFn, aFn)
            filteredLines = 0
            passedLines = 0
            iline = 0
            curFrame = ""
            Trace(filterFn, listFn)
        case cmd == "v":
            if Ctx.T == 0 && !Ctx.F {
                Trace(filterFn, countFn)
                if scanlineCount > 1 {
                    log.Fatal("Use -f flag to force v/done/fossil command on multiple scanlines")
                }
            }
            Apply(filterFn, fossilFn)
        case cmd == "x":
            if Ctx.T == 0 && !Ctx.F {
                Trace(filterFn, countFn)
                if scanlineCount > 1 {
                    log.Fatal("Use -f flag to force x/cancel command on multiple scanlines")
                }
            }
            Apply(filterFn, cancelFn)
        case cmd == "d":
            if Ctx.T == 0 && !Ctx.F {
                Trace(filterFn, countFn)
                if scanlineCount > 1 {
                    log.Fatal("Use -f flag to force d/delete/remove command on multiple scanlines")
                }
            }
            Apply(filterFn, deleteFn)
        case cmd == "s":
            if Ctx.T == 0 && !Ctx.F {
                Trace(filterFn, countFn)
                if scanlineCount > 1 {
                    log.Fatal("Use -f flag to force x/score command on multiple scanlines")
                }
            }
            Apply(filterFn, scoreFn)
        case cmd == "l": Trace(filterFn, listFn)
        default: Trace(filterFn,listFn)
    }

    // write changes
    for _, f := range frames {
        f.save()
    }
}

