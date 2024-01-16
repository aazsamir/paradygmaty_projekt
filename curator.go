package main

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"strconv"
)

type Prolog struct {
	executable string
	database   string
}

func NewProlog() *Prolog {
	return &Prolog{
		executable: "swipl",
		database:   "curator.pl",
	}
}

func (p *Prolog) Query(id int) bool {
	cmd := exec.Command(p.executable, p.database, strconv.Itoa(id))
	out, err := cmd.Output()

	if err != nil {
		panic(err)
	}

	return string(out) == "1"
}

func (p *Prolog) save(id int) error {
	if p.Query(id) {
		return nil
	}

	// append line
	file, err := os.OpenFile(p.database, os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString("\nrelated(" + strconv.Itoa(id) + ", 1).\n")

	if err != nil {
		return err
	}

	return nil
}

func (p *Prolog) delete(id int) error {
	if !p.Query(id) {
		return nil
	}

	file, err := os.OpenFile(p.database, os.O_RDWR, 0644)

	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	var bs []byte
	buf := bytes.NewBuffer(bs)

	for scanner.Scan() {
		// read line
		line := scanner.Text()

		if line == "related("+strconv.Itoa(id)+", 1)." {
			continue
		}

		buf.WriteString(line + "\n")
	}
	file.Truncate(0)
	file.Seek(0, 0)
	_, err = buf.WriteTo(file)

	if err != nil {
		return err
	}

	return nil
}
