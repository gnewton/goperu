package main

import (
	"fmt"
	"github.com/icza/gowut/gwu"
	"io/ioutil"
	//"io"
	"os"
	"strconv"
)

var expanderDir map[gwu.Expander]string

var handler = new(ZHandler)
var foobar = new(Foobar)

var baseDir = "/"

func main() {
	expanderDir = make(map[gwu.Expander]string)
	// Create and build a window
	win := gwu.NewWindow("main", "Test GUI Window")
	win.Style().SetFullWidth()
	win.SetHAlign(gwu.HA_LEFT)
	win.SetCellPadding(2)

	//win.SetTheme(gwu.THEME_DEBUG)

	p := gwu.NewPanel()
	p.Style().SetFullWidth()
	hostname, err := os.Hostname()
	if err != nil {
		hostname = ""
	}
	l := gwu.NewLabel(hostname + ": Explore           ----- goperu")
	l.Style().SetColor(gwu.CLR_GREEN)
	p.Add(l)
	p.AddVSpace(50)

	e := gwu.NewExpander()
	e.SetHeader(gwu.NewLabel(baseDir))

	e = makeExpander(baseDir)

	expanderDir[e] = baseDir

	p.Add(e)
	e.AddEHandler(handler, gwu.ETYPE_STATE_CHANGE)
	win.Add(p)

	// Create and start a GUI server (omitting error check)
	server := gwu.NewServer("guitest", "localhost:8081")
	server.SetText("Test GUI App")
	server.AddWin(win)
	server.Start("") // Also opens windows list in browser
}

type ZHandler struct {
}

type Foobar struct {
}

func makeExpander(dir string) gwu.Expander {
	e := gwu.NewExpander()
	e.SetHeader(gwu.NewLabel(dir))
	e.Style().SetPaddingLeft("60")
	expanderDir[e] = dir
	return e
}

const layout = "Jan 2, 2006 at 3:04pm (MST)"

var workingHeader = gwu.NewLabel("------Working------")

func (zz ZHandler) HandleEvent(e gwu.Event) {
	fmt.Println("----------------------------------")
	fmt.Println(e)

	fmt.Println(e.Src())
	switch t := e.Src().(type) {
	case gwu.Expander:

		if !t.Expanded() {

		} else {
			header := t.Header()
			e.MarkDirty(t)
			t.SetHeader(workingHeader)
			dir, _ := expanderDir[t]
			fmt.Println("Dir=" + dir)
			files, _ := ioutil.ReadDir(dir)
			var p gwu.Panel
			var table gwu.Table
			p = nil

			dirList := make([]gwu.Comp, 0, 0)
			fileLabelList := make([]gwu.Comp, 0, 0)
			fileTimeList := make([]gwu.Comp, 0, 0)
			fileSizeList := make([]gwu.Comp, 0, 0)

			for _, f := range files {
				if p == nil {
					p = gwu.NewPanel()
					t.SetContent(p)
					table = gwu.NewTable()
					table.SetCellPadding(5)
					p.Add(table)
				}
				if f.IsDir() {
					newExpander := makeExpander(f.Name())
					newExpander.SetToolTip("DIR - Click to view contents")
					newExpander.Style().SetColor("blue")
					if dir == "/" {
						expanderDir[newExpander] = dir + f.Name()
					} else {
						expanderDir[newExpander] = dir + "/" + f.Name()
					}
					newExpander.AddEHandler(handler, gwu.ETYPE_STATE_CHANGE)
					dirList = append(dirList, newExpander)
				} else {
					fileLabel := gwu.NewLabel(f.Name())
					fileLabel.SetToolTip("FILE")
					fileLabel.Style().SetColor("red")
					fileLabel.AddEHandler(foobar, gwu.ETYPE_MOUSE_DOWN)
					fileLabelList = append(fileLabelList, fileLabel)

					fileTimeLabel := gwu.NewLabel(f.ModTime().Format(layout))
					fileTimeList = append(fileTimeList, fileTimeLabel)
					fileSizeLabel := gwu.NewLabel(strconv.FormatInt(f.Size(), 10))
					fileSizeList = append(fileSizeList, fileSizeLabel)
				}
			}
			row := 0
			for _, dir := range dirList {
				table.Add(dir, row, 0)
				row = row + 1
			}
			for i, _ := range fileLabelList {
				table.Add(fileLabelList[i], row, 0)
				rfmt := table.CellFmt(row, 0)
				rfmt.SetHAlign(gwu.HA_LEFT)
				rfmt.SetVAlign(gwu.VA_TOP)
				table.Add(fileTimeList[i], row, 1)
				rfmt = table.CellFmt(row, 1)
				rfmt.SetHAlign(gwu.HA_LEFT)
				rfmt.SetVAlign(gwu.VA_TOP)

				row = row + 1
			}
			t.SetHeader(header)
		}
	default:
		// t is some other type that we didn't name.
	}

}

func (h Foobar) HandleEvent(e gwu.Event) {
	fmt.Println("******************************")
}
