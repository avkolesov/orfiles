package subfiles

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
)

// StFilesCnfg - config orolia files
type StFilesCnfg struct {
	Path  string
	Type  string
	Names []string
}

// StFile - orolia file data
type StFile struct {
	Name   string
	Values [][]string
}

func _checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// CheckFiles - check all orolia files in derictory and get names
func (fconf *StFilesCnfg) CheckFiles() {
	files, err := ioutil.ReadDir(fconf.Path)
	_checkErr(err)
	fconf.Names = []string{}
	for _, f := range files {
		name := f.Name()
		if strings.HasSuffix(name, fconf.Type) {
			fconf.Names = append(fconf.Names, name)
		}
	}
}

// AddFile - add file to derictory
func (fconf *StFilesCnfg) AddFile(file multipart.File, addName string) {
	if addName != "" {
		f, err := os.OpenFile(fconf.Path+"/"+addName, os.O_WRONLY|os.O_CREATE, 0666)
		_checkErr(err)
		defer f.Close()
		io.Copy(f, file)
	}
}

// DelFile - delete file from derictory
func (fconf *StFilesCnfg) DelFile(delName string) {
	os.Remove(fconf.Path + "/" + delName)
}

// GetFileValue - get values form orolia file
func (f *StFile) GetFileValue(path string) string {
	f.Values = [][]string{}
	content, err := os.Open(path + "/" + f.Name)
	_checkErr(err)
	defer content.Close()
	scanner := bufio.NewScanner(content)
	res := `[["X","Y"],`
	for scanner.Scan() {
		if !strings.HasPrefix(scanner.Text(), "#") {
			text := scanner.Text()
			text = strings.Replace(text, "\t", "", -1)
			if text != "" {
				spl := strings.Split(text, " ")
				x, _ := strconv.ParseFloat(spl[0], 64)
				y, _ := strconv.ParseFloat(spl[1], 64)
				res = res + "[" + strconv.FormatFloat(x, 'f', -1, 64) + "," + strconv.FormatFloat(y, 'f', -1, 64) + "],"
				f.Values = append(f.Values, strings.Split(text, " "))
			}
		}
	}
	res = res[0:len(res)-1] + "]"
	return res
}

// FileInit - initialization file struct
func (f *StFile) FileInit(name string) {
	f.Name = name
}
