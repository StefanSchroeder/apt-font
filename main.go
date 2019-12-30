/*

  Searching the internal font database for the font provided
  on the cmdline.
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/StefanSchroeder/odtfontfind"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var data map[string]string

func find_package_name(p string) {
	path, err1 := exec.LookPath("dpkg")
	pathapt, err2 := exec.LookPath("apt")
	if err1 != nil || err2 != nil {
		log.Printf("dpkg or apt is not available. This is not a supported system.")
		return
	}

	d, err := ioutil.TempDir("", "apt-font-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(d)

	cmd := exec.Command(pathapt, "download", p)
	cmd.Dir = d
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	glog_result, err := filepath.Glob(d + "/*.deb")
	if err != nil {
		log.Fatal(err)
	}

	cmd = exec.Command(path, "-x", glog_result[0], d)
	cmd.Dir = d
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	h := os.Getenv("HOME")
	fontdir := filepath.Join(h, ".fonts")
	os.MkdirAll(fontdir, 0700)

	err = filepath.Walk(d, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error while walking tempdir.")
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".ttf") {
			log.Printf("=> %v\n", info.Name())

			copy(path, info.Name(), fontdir)
		}
		return nil
	})
	if err != nil {
		log.Printf("Error while walking tempdir\n")
	}
}

// copy is a generic file copier.
func copy(from string, basename string, to string) {
	fh, err := os.Open(from)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	th, err := os.OpenFile(filepath.Join(to, basename), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer th.Close()
	_, err = io.Copy(th, fh)
	if err != nil {
		log.Fatal(err)
	}
}

// install_named_font will lookup the provided name
// in the configuration and install it locally.
func install_named_font(s string) {
	log.Println("Searching for: ", s)
	if value, ok := data[s]; ok {
		log.Println("OK: ", s, value)
		find_package_name(value)
	} else {
		log.Println("Not found: ", s)
	}
}

func process_document(s string) {
	fmt.Printf("Processing " + s + "\n")

	fi, err := os.Lstat(s)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch mode := fi.Mode(); {
	case mode.IsRegular():
		r := odtfontfind.LibreofficeFontReader(s)
		for _, i := range r {
			fmt.Println(i)
		}
	case mode.IsDir():
		err = filepath.Walk(s, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("Error while walking tempdir.")
				return err
			}
			for _, suf := range []string{".odg", ".odt", ".ods", ".odp"}  {
				if !info.IsDir() && strings.HasSuffix(info.Name(), suf) {
					fmt.Println("=>" + path)
					r := odtfontfind.LibreofficeFontReader(path)
					for _, i := range r {
						install_named_font(i)
					}
				}
			}
			return err
		})
		if err != nil {
			log.Printf("Error while walking directory.\n")
		}

	}
}

func main() {
	configFile := flag.String("c", "fonts.json", "Configuration file")
	wantShow := flag.Bool("s", false, "Display configuration")
	wantHelp := flag.Bool("h", false, "Help.")
	wantFont := flag.Bool("f", false, "Local font install.")
	wantDocs := flag.Bool("i", false, "Local font install from documents.")

	flag.Parse()

	file, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	_ = json.Unmarshal([]byte(file), &data)

	if *wantHelp {
		flag.PrintDefaults()
		return
	}

	if *wantShow {
		for i := range data {
			fmt.Println(i, "=>", data[i])
		}
		return
	}

	if *wantFont {
		for _, i := range flag.Args() {
			install_named_font(i)
		}
		return
	}

	if *wantDocs {
		for _, i := range flag.Args() {
			process_document(i)
		}
		return
	}
	fmt.Println("Nothing to do. Use -h for help")
}
