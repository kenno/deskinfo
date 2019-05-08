package main

import (
    "bufio"
    "bytes"
    "flag"
    "fmt"
    "os"
    "regexp"
    "syscall"

    "github.com/fatih/color"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func grep(re, filename string) string {
    regex, err := regexp.Compile(re)
    check(err)

    fh, err := os.Open(filename)
    f := bufio.NewReader(fh)
    check(err)
    defer fh.Close()

    buf := make([]byte, 1024)
    for {
        buf, _, err = f.ReadLine()
        if err != nil {
            return string(buf)
        }

        s := string(buf)
        if regex.MatchString(s) {
            //fmt.Printf("%s\n", string(buf))
            return string(buf)
        }
    }
}

func getuptime() (float64, error) {
    info := syscall.Sysinfo_t{}
    if err := syscall.Sysinfo(&info); err != nil {
        return 0, err
    }
    //fmt.Printf("uptime: ", float64(info.Uptime))

    return float64(info.Uptime), nil
}

func printuptime(fieldName bool) {
    buf := new(bytes.Buffer)
    w := bufio.NewWriter(buf)

    u, _ := getuptime()
    uptime := uint64(u)

    days := uptime / (60 * 60 * 24)

    if days != 0 {
        s := ""
        if days > 1 {
            s = "s"
        }
        fmt.Fprintf(w, "%d day%s ", days, s)
    }

    minutes := uptime / 60
    hours := minutes / 60
    hours %= 24
    minutes %= 60

    fmt.Fprintf(w, "%2d:%02d", hours, minutes)

    w.Flush()

    if fieldName {
        color.Set(color.FgBlue, color.Bold)
        fmt.Print("Uptime: ")
        color.Unset()
    }
    fmt.Println(buf.String())
}

func printhostname(fieldName bool) {
    hostName, err := os.Hostname()
    check(err)

    if fieldName {
        color.Set(color.FgBlue, color.Bold)
        fmt.Print("Hostname: ")
        color.Unset()
    }
    fmt.Println(hostName)
}

func printos(fieldName bool) {
    // read /etc/lsb-release to get the codename
    codename := grep("DISTRIB_CODENAME", "/etc/lsb-release")

    deskCodename := ""

    switch codename {
    case "DISTRIB_CODENAME=bionic":
        deskCodename = "Desk 11"
    case "DISTRIB_CODENAME=xenial":
        deskCodename = "Desk 10"
    case "DISTRIB_CODENAME=trusty":
        deskCodename = "Desk 9"
    case "DISTRIB_CODENAME=precise":
        deskCodename = "Desk 8"
    default:
        deskCodename = "Unknown"
    }

    if fieldName {
        color.Set(color.FgBlue, color.Bold)
        fmt.Print("OS: ")
        color.Unset()
    }
    fmt.Println(deskCodename)
}

func main() {

    filterPtr := flag.String("filter", "all", "Filter output by field: all, os, host, uptime")
    fieldNamePtr := flag.Bool("name", true, "Field name: true, false")
    flagColor := flag.Bool("no-color", false, "Disable colour output")

    flag.Parse()

    if *flagColor {
        color.NoColor = true // disable colorized output
    }

    switch *filterPtr {
    case "os":
        printos(*fieldNamePtr)
    case "host":
        printhostname(*fieldNamePtr)
    case "uptime":
        printuptime(*fieldNamePtr)
    default:
        printos(*fieldNamePtr)
        printhostname(*fieldNamePtr)
        printuptime(*fieldNamePtr)
    }
}
