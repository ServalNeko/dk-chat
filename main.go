package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	textView := tview.NewTextView().
		SetChangedFunc(func() {
			app.Draw()
		})
	textView.SetBorder(true)
	textView.SetTitle("ちゃっと")

	inputField := tview.NewInputField()
	inputField.SetTitle("にゅうりょく").
		SetBorder(true)
	cfg := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", "koukoku.shadan.open.ad.jp:992", cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	go func() {
		r := bufio.NewReader(conn)
		for {
			line, _, err := r.ReadLine()
			if err != nil {
				log.Fatal(err)
			}
			if strings.HasPrefix(string(line), ">>") {
				fmt.Fprintf(textView, "%s\n", line)
			}
		}

	}()

	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			_, err := fmt.Fprintf(conn, "%s\n", inputField.GetText())
			if err != nil {
				log.Fatal(err)
			}

			inputField.SetText("")
			return nil
		}
		return event
	})

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow).
		AddItem(inputField, 3, 0, true).
		AddItem(textView, 0, 1, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}

}
