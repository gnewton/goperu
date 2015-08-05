package main

import (
	//"fmt"
	"github.com/icza/gowut/gwu"
	"io/ioutil"
	//"io"
	"os"
	"strconv"
	"time"
)

var expanderDir map[gwu.Expander]string

var handler = new(ExpanderHandler)
var khandler = new(KHandler)
var foobar = new(Foobar)

var baseDir = "/"

func main() {
	expanderDir = make(map[gwu.Expander]string)
	// Create and build a window
	win := gwu.NewWindow("main", "goperu - peruse files with browser")
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

	link := gwu.NewLink("goperu", "https://github.com/gnewton/goperu")
	p.Add(link)
	p.AddVSpace(10)

	l := gwu.NewLabel(hostname + ": Explore")
	l.Style().SetColor(gwu.CLR_GREEN)
	p.Add(l)

	p.AddVSpace(20)

	e := makeExpander(baseDir)

	expanderDir[e] = baseDir

	p.Add(e)
	e.AddEHandler(handler, gwu.ETYPE_STATE_CHANGE)
	//e.AddEHandler(handler, gwu.ETYPE_MOUSE_DOWN)
	win.Add(p)

	// Create and start a GUI server (omitting error check)
	server := gwu.NewServer("goperu", "localhost:8081")
	server.SetText("goperu")
	server.AddWin(win)
	server.Start("") // Also opens windows list in browser
}

type ExpanderHandler struct {
}

type KHandler struct {
}

type Foobar struct {
}

func makeExpander(dir string) gwu.Expander {
	e := gwu.NewExpander()
	e.SetHeader(gwu.NewLabel(dir))
	e.Style().SetPaddingLeft("20")
	expanderDir[e] = dir
	return e
}

const layout = time.RFC822

func (kk KHandler) HandleEvent(e gwu.Event) {

}

func (zz ExpanderHandler) HandleEvent(e gwu.Event) {

	switch t := e.Src().(type) {
	case gwu.Expander:
		dir, _ := expanderDir[t]
		//fmt.Println("Dir=" + dir)
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
				// owned by root // only works on linux
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
				fileLabel.Style().SetPaddingLeft("20")
				fileLabel.SetToolTip("FILE")
				fileLabel.AddEHandler(foobar, gwu.ETYPE_MOUSE_DOWN)
				fileLabelList = append(fileLabelList, fileLabel)

				fileTimeLabel := gwu.NewLabel(f.ModTime().Format(layout))
				fileTimeList = append(fileTimeList, fileTimeLabel)
				fileSizeLabel := gwu.NewLabel(strconv.FormatInt(f.Size(), 10) + "k")
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

			table.Add(fileSizeList[i], row, 2)
			rfmt = table.CellFmt(row, 2)
			rfmt.SetHAlign(gwu.HA_RIGHT)
			rfmt.SetVAlign(gwu.VA_TOP)

			row = row + 1
		}

	default:
		// t is some other type that we didn't name.
	}

}

func (h Foobar) HandleEvent(e gwu.Event) {

}
