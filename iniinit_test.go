package main

import (
	"os"
	"reflect"
	"testing"
)

func TestIniInit(t *testing.T) {
	cases := []struct {
		dataBotIn  TBot
		dataConfig string
		wanTBotOut TBot
		wantResult bool
	}{
		{
			TBot{
				"",
				"stdout",
				os.Stdout,
				"",
				"localhost",
				"80",
				nil},
			`
			logfile = 
			token = 1199966707:AAFHcuGFvCC_zVh6r_XZ-bjXXXhwXM-Uw6k
			address = 127.0.0.1
			port = 8889
			`,
			TBot{
				"",
				"stdout",
				os.Stdout,
				"1199966707:AAFHcuGFvCC_zVh6r_XZ-bjXXXhwXM-Uw6k",
				"127.0.0.1",
				"8889",
				nil},
			true,
		},
		{
			TBot{
				"",
				"stdout",
				os.Stdout,
				"",
				"localhost",
				"80",
				nil},
			`

			`,
			TBot{"",
				"stdout",
				os.Stdout,
				"",
				"localhost",
				"80",
				nil},
			true,
		},
		{
			TBot{"",
				"stdout",
				os.Stdout,
				"",
				"localhost",
				"80",
				nil},
			`
			logfile = "file.log"
			blablabla = "Bla bla bla"
			token = 1199966707:AAFHcuGFvCC_zVh6r_XZ-bjXXXhwXM-Uw6k
			address = 127.0.0.1
			port = 8889
			`,
			TBot{"",
				"file.log",
				os.Stdout,
				"1199966707:AAFHcuGFvCC_zVh6r_XZ-bjXXXhwXM-Uw6k",
				"127.0.0.1",
				"8889",
				nil},
			true,
		},
	}
	for _, c := range cases {
		bot := c.dataBotIn
		gotResult := iniInit(&bot, []byte(c.dataConfig))
		if (gotResult != c.wantResult) || (!reflect.DeepEqual(bot, c.wanTBotOut)) {
			t.Errorf("For \nbotIn == %v \nand \nconfig == %v: \ngot \nbotOut == %v \nand \nresult == %v, \nvant \nbotOut == %v \nand \nresult == %v", c.dataBotIn, c.dataConfig, bot, gotResult, c.wanTBotOut, c.wantResult)
		}
	}
}
