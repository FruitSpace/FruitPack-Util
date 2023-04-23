package main

import (
	"archive/zip"
	"fmt"
	"golang.org/x/exp/slices"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

var Deep bool

var Blacklist = []string{".psd",".jpg", ".txt"}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Drag \"Resources\" folder to this executable")
		nop()
	}
	srcDir := os.Args[1]
	info, err := os.Stat(srcDir)
	if err != nil {
		fmt.Println(err)
		nop()
	}
	if !info.IsDir() {
		fmt.Println("This is NOT A FOLDER... Please drag \"Resources\" folder to this executable")
		nop()
	}
	srcDir = strings.ReplaceAll(srcDir, "\\", "/")

	deltafiles:=LDir(srcDir, "")
	t:= time.Now()
	fname:=fmt.Sprintf("TexturePack_%d%d%d_%d%d%d.fpack",t.Year(),t.Month(),t.Day(),t.Hour(),t.Minute(),t.Second())

	zipFile, err := os.Create(fname)
	if err!=nil {
		fmt.Println(err)
		nop()
	}
	defer zipFile.Close()


	zipW := zip.NewWriter(zipFile)
	defer zipW.Close()
	//zipW.RegisterCompressor(zip.Deflate, func(w io.Writer) (io.WriteCloser, error) {
	//	return zstd.NewWriter(w)
	//})

	for _, pack := range deltafiles {
		cPack, err := zipW.Create(pack)
		uPack, err := os.Open(path.Join(srcDir,pack))
		if err!=nil {
			fmt.Println(err)
			nop()
		}

		_, err = io.Copy(cPack, uPack)
		if err!=nil {
			fmt.Println(err)
			nop()
		}
	}
	fmt.Printf("Texture Pack \"%s\" was successfully created\n", fname)

	nop()
}

func LDir(srcDir string, prefix string) []string {
	files, err := os.ReadDir(srcDir)
	if err!=nil {return nil}
	var dlist []string
	for _, file := range files{
		pr := path.Join(srcDir, file.Name())
		if file.IsDir() {
			flist := LDir(pr, prefix+file.Name()+"/")
			dlist = append(dlist, flist...)
		} else {
			// file
			p := prefix+file.Name()
			if slices.Contains(Blacklist, path.Ext(file.Name())) {
				fmt.Println("Ignore:", p)
				continue
			}
			target, ok := FileList[p]
			if !ok {
				// New file
				dlist = append(dlist, p)
				fmt.Println("Add:",p)
				continue
			}
			if !target.Compare(pr, Deep) {
				dlist = append(dlist, p)
				fmt.Println("Add:",p)
			}
			//fh := *NewFileHold(path.Join(srcDir, file.Name()))
			//fmt.Printf("\"%s\": FileHold{Size:%d, AccessTime:%d, Hash:\"%s\"},\n",
			//	prefix+file.Name(), fh.Size, fh.AccessTime, fh.Hash)
		}
	}
	return dlist
}


func nop() {
	fmt.Println("(Ctrl+C to Exit)")
	for {}
	os.Exit(1)
}