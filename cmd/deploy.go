package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

const (
	WORK_DIR     = "./work/"
	INPUT_DIR    = "./cmd/"
	OUTPUT_DIR   = "./app/static/js/"
	TEMPLATE_DIR = "./app/templates/entry/"
	INPUT        = "editor.go"
	OUTPUT_GO    = "editor.go.map"
	OUTPUT_JS    = "editor.js"
	OUTPUT_MAP   = "editor.js.map"
)

func init() {
}

func main() {

	if err := os.RemoveAll(WORK_DIR); err != nil {
		fmt.Println(err)
	}

	defer os.RemoveAll(WORK_DIR)
	err := run()
	if err != nil {
		fmt.Println(err)
	}
}

func createGoFile() error {

	article_tmpl := TEMPLATE_DIR + "entry.tmpl"
	article, err := ioutil.ReadFile(article_tmpl)
	if err != nil {
		fmt.Println("article template read error")
		return err
	}
	action_tmpl := TEMPLATE_DIR + "action.tmpl"
	action, err := ioutil.ReadFile(action_tmpl)
	if err != nil {
		fmt.Println("action template read error")
		return err
	}

	embed := struct {
		ARTICLE_TMPL string
		ACTION_TMPL  string
	}{string(article), string(action)}

	file, err := os.Create(WORK_DIR + INPUT)
	if err != nil {
		fmt.Println("WorkFile create error")
		return err
	}
	defer file.Close()

	tpl := template.Must(template.ParseFiles(INPUT_DIR + INPUT))
	if err := tpl.Execute(file, embed); err != nil {
		fmt.Println("WorkFile Execute error")
		return err
	}
	return nil
}

func generateJSFile() error {

	err := Command("gopherjs", "build", "-m", WORK_DIR+INPUT, "-o", WORK_DIR+OUTPUT_JS)
	if err != nil {
		fmt.Println("GopherJS build error")
		return err
	}
	return nil
}

func rename() error {

	// editor.go -> editor.go.map
	err := os.Rename(WORK_DIR+INPUT, WORK_DIR+OUTPUT_GO)
	if err != nil {
		fmt.Println("Rename error")
		return err
	}

	//Change map file
	file, err := os.Open(WORK_DIR + OUTPUT_MAP)
	if err != nil {
		fmt.Println("Open Map file error")
		return err
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Read Map file error")
		return err
	}

	err = file.Close()
	if err != nil {
		fmt.Println("Close Map file error")
		return err
	}

	buf := string(b)
	slice := strings.Split(buf, INPUT)
	newBuf := slice[0] + OUTPUT_GO + slice[1]

	err = ioutil.WriteFile(WORK_DIR+OUTPUT_MAP, []byte(newBuf), os.ModePerm)
	if err != nil {
		fmt.Println("Write Map file error")
		return err
	}

	return nil
}

func deploy() error {
	err := os.Rename(WORK_DIR+OUTPUT_GO, OUTPUT_DIR+OUTPUT_GO)
	if err != nil {
		fmt.Println("Go File Move error")
		return err
	}
	err = os.Rename(WORK_DIR+OUTPUT_JS, OUTPUT_DIR+OUTPUT_JS)
	if err != nil {
		fmt.Println("JS File Move error")
		return err
	}
	err = os.Rename(WORK_DIR+OUTPUT_MAP, OUTPUT_DIR+OUTPUT_MAP)
	if err != nil {
		fmt.Println("Map File Move error")
		return err
	}
	return nil
}

func run() error {

	fmt.Println("Create Work directory")

	err := os.Mkdir(WORK_DIR, 0777)
	if err != nil {
		fmt.Println("Create work error")
		return err
	}

	fmt.Println("Create Go File")
	err = createGoFile()
	if err != nil {
		return err
	}

	fmt.Println("Generate JS File")
	err = generateJSFile()
	if err != nil {
		return err
	}

	fmt.Println("Rename")
	err = rename()
	if err != nil {
		return err
	}

	fmt.Println("Deploy")
	err = deploy()
	if err != nil {
		return err
	}

	return nil
}

func Command(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	return err
}
