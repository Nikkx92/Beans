package coffee

import (
	"bufio"
	"bytes"
	"encoding/json"
	"image"
	_ "image/png"
	"log/slog"
	"os"
	"strings"
	"voicer/images"
	"voicer/internal/logger"
	"voicer/internal/storage"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget"
)

func (c *Components) HandleEvents(gtx layout.Context) {
	if c.Settings.addCountry.Clicked(gtx) {
		block := addItemBlock()
		block.Sorts = append(block.Sorts, addSortBlock())
		c.Blocks = append(c.Blocks, block)
		c.ListBlocks.Position.First = len(c.Blocks) - 5
		c.animatingFunc = func() {
			if c.ListBlocks.Position.BeforeEnd == false && c.ListBlocks.Position.OffsetLast == 0 {
				c.isAnimating = false
			}
		}
		c.isAnimating = true
	}

	if c.Settings.confirmDeleting.Clicked(gtx) {
		c.deleteFunc()
		c.approveDeleteModal.Disappear(gtx.Now)
	}

	if c.Settings.rejectDeleting.Clicked(gtx) {
		c.approveDeleteModal.Disappear(gtx.Now)
	}

	for i := len(c.Blocks) - 1; i >= 0; i-- {
		b := &c.Blocks[i]

		if b.seekFlag.Clicked(gtx) {
			country := b.fieldCountry.Text()
			if country == "" {
				//modal with caution
				return
			}

			b.flag = c.flags[country]

			if b.flag == nil {
				fl := images.LoadImage(country)

				if fl == nil {
					b.loader = true

					go func() {
						defer func() { b.loader = false }()

						fBytes, err := c.a.Client.GetFlag(country)
						if err != nil {
							slog.Info(err.Error())
							return
						}
						img, _, err := image.Decode(bytes.NewReader(fBytes))
						if err != nil {
							slog.Info(err.Error())
							return
						}

						fl = &widget.Image{
							Src: paint.NewImageOp(img),
							Fit: widget.Fill,
						}

						b.flag = fl
					}()
				}

				b.flag = fl
			}

		}

		if b.deleteCountry.Clicked(gtx) {
			c.delModal(gtx, b.fieldCountry.Text())
			c.deleteFunc = func() {
				c.Blocks = append(c.Blocks[:i], c.Blocks[i+1:]...)
			}
			continue
		}

		if b.addSort.Clicked(gtx) {
			b.Sorts = append(b.Sorts, addSortBlock())
			currentItem := c.ListBlocks.Position.Count
			c.animatingFunc = func() {
				if c.ListBlocks.Position.Count != currentItem {
					c.isAnimating = false
				}
			}
			c.isAnimating = true
		}

		if b.upCountryBtn.Clicked(gtx) {
			if i != 0 {
				c.Blocks[i], c.Blocks[i-1] = c.Blocks[i-1], c.Blocks[i]
			}
		}

		if b.downCountryBtn.Clicked(gtx) {
			if i != len(c.Blocks)-1 {
				c.Blocks[i], c.Blocks[i+1] = c.Blocks[i+1], c.Blocks[i]
			}
		}

		//upperCase country title
		for {
			e, ok := b.fieldCountry.Update(gtx)
			if !ok {
				break
			}

			if _, ok := e.(widget.ChangeEvent); ok {
				currentText := b.fieldCountry.Text()
				upperText := strings.ToUpper(currentText)
				if currentText != upperText {
					start, end := b.fieldCountry.Selection()
					b.fieldCountry.SetText(upperText)
					b.fieldCountry.SetCaret(start, end)
				}
			}
		}

		for j := len(b.Sorts) - 1; j >= 0; j-- {
			s := &b.Sorts[j]

			if s.deleteSort.Clicked(gtx) {
				c.delModal(gtx, s.fieldSort.Text())
				c.deleteFunc = func() {
					c.Blocks[i].Sorts = append(c.Blocks[i].Sorts[:j], c.Blocks[i].Sorts[j+1:]...)
					if len(c.Blocks[i].Sorts) == 0 {
						c.Blocks = append(c.Blocks[:i], c.Blocks[i+1:]...)
					}
				}
				continue
			}

			if s.up.Clicked(gtx) {
				if j != 0 {
					c.Blocks[i].Sorts[j], c.Blocks[i].Sorts[j-1] = c.Blocks[i].Sorts[j-1], c.Blocks[i].Sorts[j]
				}
			}

			if s.down.Clicked(gtx) {
				if j != len(c.Blocks[i].Sorts)-1 {
					c.Blocks[i].Sorts[j], c.Blocks[i].Sorts[j+1] = c.Blocks[i].Sorts[j+1], c.Blocks[i].Sorts[j]
				}
			}

			//upperCase sort
			for {
				e, ok := s.fieldSort.Update(gtx)
				if !ok {
					break
				}

				if _, ok := e.(widget.ChangeEvent); ok {
					currentText := s.fieldSort.Text()
					upperText := strings.ToUpper(currentText)
					if currentText != upperText {
						start, end := s.fieldSort.Selection()
						s.fieldSort.SetText(upperText)
						s.fieldSort.SetCaret(start, end)
					}
				}
			}

			//first letter is up
			for _, v := range []*widget.Editor{
				s.fieldBlend,
				s.fieldDrying,
				s.fieldLocation,
				s.fieldTaste,
			} {
				for {
					e, ok := v.Update(gtx)
					if !ok {
						break
					}

					if _, ok := e.(widget.ChangeEvent); ok {
						if v.Len() == 0 {
							continue
						}

						currentText := v.Text()
						r := []rune(currentText)
						upperText := strings.ToUpper(currentText)
						ru := []rune(upperText)
						if string(r[0]) != string(ru[0]) {
							start, end := v.Selection()
							v.SetText(string(ru[0]))
							v.SetCaret(start, end)
						}
					}
				}
			}
		}

	}

	if c.Settings.closeSettings.Clicked(gtx) {
		c.a.AutoSaveCancelHelper()
		c.a.DevelopMode.Store(false)
	}

	if c.Settings.logBtn.Clicked(gtx) {
		messages := getLogMessages()

		c.logModal(gtx, messages)
	}

	for _, val := range c.Settings.SettingItems {
		for _, val1 := range val.Parameters {
			for {
				ev, ok := val1.Editor.Update(gtx)
				if !ok {
					break
				}

				m := map[string]map[string]*storage.Type{}
				if c.a.MobileMode.Load() {
					m = c.Settings.FontSzM
				} else {
					m = c.Settings.FontSz
				}

				if _, ok := ev.(widget.ChangeEvent); ok {
					d := aToI(val1.Editor.Text())
					m[val.Title][val1.Title].Str = val1.Editor.Text()
					m[val.Title][val1.Title].Num = d
				}
			}
		}
	}

	if c.Settings.closeSavePDF.Clicked(gtx) {
		c.savePdf.Disappear(gtx.Now)
	}

}

func addItemBlock() ItemBlock {
	return ItemBlock{
		seekFlag:       new(widget.Clickable),
		loader:         false,
		flag:           nil,
		fieldCountry:   new(widget.Editor),
		addSort:        new(widget.Clickable),
		deleteCountry:  new(widget.Clickable),
		upCountryBtn:   new(widget.Clickable),
		downCountryBtn: new(widget.Clickable),
		Sorts:          []SortBlock{},
	}

}

func addSortBlock() SortBlock {
	return SortBlock{
		fieldSort:  new(widget.Editor),
		fieldBlend: new(widget.Editor),
		fieldGrade: &widget.Editor{
			Filter: "0,1,2,3,4,5,6,7,8,9,`,`,`.` ",
		},
		fieldLocation:    new(widget.Editor),
		fieldHeightAbove: new(widget.Editor),
		fieldDrying:      new(widget.Editor),
		fieldCornSize:    new(widget.Editor),
		fieldTaste:       new(widget.Editor),
		fieldPrice: &widget.Editor{
			Filter:    "0,1,2,3,4,5,6,7,8,9",
			Alignment: text.Alignment(layout.End),
		},
		fieldQuantity: &widget.Editor{
			Alignment: text.Alignment(layout.End),
		},
		fieldStockDesc: &widget.Editor{
			Alignment: text.End,
		},
		deleteSort: new(widget.Clickable),
		up:         new(widget.Clickable),
		down:       new(widget.Clickable),
		inStock:    new(widget.Bool),
		limited:    new(widget.Bool),
	}
}

func getLogMessages() []logger.Message {
	path, err := storage.Path()
	if err != nil {
		slog.Info(err.Error())
	}
	path = path + "/app.log"

	file, err := os.Open(path)
	if err != nil {
		slog.Info(err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var messages []logger.Message
	for scanner.Scan() {
		line := scanner.Bytes()
		var message logger.Message
		if err := json.Unmarshal(line, &message); err != nil {
			message = logger.Message{Message: string(line)}
		}
		messages = append(messages, message)
	}
	return messages
}
