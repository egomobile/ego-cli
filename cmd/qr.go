// This file is part of the ego-cli distribution.
// Copyright (c) Next.e.GO Mobile SE, Aachen, Germany (https://e-go-mobile.com/)
//
// ego-cli is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as
// published by the Free Software Foundation, version 3.
//
// ego-cli is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
// Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/briandowns/spinner"
	"github.com/skip2/go-qrcode"
	"github.com/thatisuday/commando"

	utils "github.com/egomobile/ego-cli/utils"
)

func qr_generate(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	recoveryLevel := qrcode.Medium

	extraHigh, err := flags["extra-high"].GetBool()
	if err == nil && extraHigh {
		recoveryLevel = qrcode.Highest
	} else {
		high, err := flags["high"].GetBool()
		if err == nil && high {
			recoveryLevel = qrcode.High
		} else {
			medium, err := flags["medium"].GetBool()
			if err == nil && medium {
				recoveryLevel = qrcode.Medium
			} else {
				low, err := flags["low"].GetBool()
				if err == nil && low {
					recoveryLevel = qrcode.Low
				}
			}
		}
	}

	size, err := flags["size"].GetInt()
	if err != nil {
		size = 1024
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	png, err := qrcode.Encode("https://example.org", recoveryLevel, size)
	if err != nil {
		log.Fatal(err)
	}

	var pngFile = filepath.Join(cwd, "qrcode.png")
	var i = 0
	var tryFindNextFilename func()
	tryFindNextFilename = func() {
		if _, err := os.Stat(pngFile); !os.IsNotExist(err) {
			i = i + 1
			pngFile = filepath.Join(cwd, "qrcode-"+strconv.Itoa(i)+".png")

			tryFindNextFilename()
		}
	}

	utils.WithSpinner("Try find unique name for output file", func(s *spinner.Spinner) {
		tryFindNextFilename()
	})

	utils.WithSpinner("Write QR code to "+path.Base(pngFile)+" ...", func(s *spinner.Spinner) {
		err = ioutil.WriteFile(pngFile, png, 0644)
		if err != nil {
			log.Fatal(err)
		}
	})
}

func Setup_qr_Command() {
	commando.
		Register("qr").
		SetShortDescription("generate QR code").
		SetDescription("Generates a QR code image in the current directory").
		AddFlag("size,s", "The size in pixels", commando.Int, 1024).
		AddFlag("medium,m", "Medium recovery level", commando.Bool, nil).
		AddFlag("high,h", "High recovery level", commando.Bool, nil).
		AddFlag("low,l", "Low recovery level", commando.Bool, nil).
		AddFlag("extra-high,x", "Highest recovery level", commando.Bool, nil).
		AddArgument("text", "The text for the QR code", "https://e-go-mobile.com").
		SetAction(qr_generate)
}
