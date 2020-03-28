package ml

import (
    "os"
    "log"
    "fmt"
    "bufio"
    "io/ioutil"
    "strings"
)

var paths []string

var frames []frame

type frame struct {
    name string
    path string
    tags []string
    lines []scanline
    modified bool
}

type scanline struct {
    frm *frame
    src string
    tags []string
}

func isLoaded(path string) bool {
    for _, p := range paths {
        if p == path {
            return true
        }
    }
    return false
}

func (f *frame) add(nl scanline) {
    nls := make([]scanline, len(f.lines)+1)
    nls[0] = nl
    copy(nls[1:], f.lines)
    f.lines = nls
}

func (f *frame) remove(li int) {
    fmt.Printf("why are we deletingt the wrong line? #%2d", li)
    //f.lines = f.lines[:li+copy(f.lines[li:], f.lines[li+1:])]
}

func (f *frame) write() {
    of, err := os.Create(f.path)
    defer of.Close()
    if err != nil {
        log.Fatal(err)
    }
    for _, l := range f.lines {
        of.WriteString(l.source() + "\n")
    }
    of.Sync()
}

func (f *frame) backup() {
    bf := read(f.name, f.path)
    bf.path = bf.path + ".bak"
    bf.write()
}

func (f *frame) save() {
    if f.modified {
        if Ctx.V {
            fmt.Println("saving", f.path)
        }
        f.backup()
        f.write()
    }
}

func (l *scanline) invisible () bool {
    if strings.TrimSpace(l.src) == "" {
        return true
    }
    return false
}

func (l *scanline) fossil () {
    if l.src[1] == ' ' {
        l.src = "V" + l.src[1:]
    } else {
        l.src = "V " + l.src
    }
}

func (l *scanline) cancel () {
    if l.src[1] == ' ' {
        l.src = "X" + l.src[1:]
    } else {
        l.src = "X " + l.src
    }
}

func (l *scanline) score() {
    l.src = l.src + " x"
}

func (l *scanline) source() string {
    return l.src
}


func Trace(filter func(l scanline, li int) bool, trace func(l scanline, li int)) {
    for _, f := range frames {
        for li, l := range f.lines {
            if filter(l, li) {
                trace(l, li)
            }
        }
    }
}

func Apply(filter func(l scanline, li int) bool,
            probe func(l scanline, li int) (*frame, *scanline)) {
    for fi, f := range frames {
        for li, l := range f.lines {
            if filter(l, li) {
                var nl *scanline
                var nf *frame
                nf, nl = probe(l, li)
                if nf != nil {
                    nf.modified = true
                    frames[fi] = *nf
                }
                if nl != nil {
                    f.modified = true
                    f.lines[li] = *nl
                }
            }
        }
    }
}

func Load(dir string) int {
    filesLoaded := 0
    if isLoaded(dir) {
        return filesLoaded
    }
    paths = append(paths, dir)
    files, _ := ioutil.ReadDir(dir)
    for _, f := range files {
        if !f.IsDir() && strings.HasSuffix(f.Name(), ".ml") {
            if Ctx.V {
                fmt.Println("== reading " + dir + "/" + f.Name())
            }
            frames = append(frames, read(f.Name(), dir + "/" + f.Name()))
            filesLoaded++
        }
    }
    return filesLoaded
}

func read(name string, path string) frame {
	file, err := os.Open(path)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    // got new .ml file to read
    nextFrame := frame{ name: name[0:len(name)-3], path: path, modified: false }

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        var text = scanner.Text()
        sl := scanline{frm: &nextFrame, src: text}
        nextFrame.lines = append(nextFrame.lines, sl)
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
    return nextFrame 
}
