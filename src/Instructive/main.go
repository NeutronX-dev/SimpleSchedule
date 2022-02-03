package Inst

import (
	"fmt"
	"strings"

	"github.com/gen2brain/dlgs"
	"github.com/pkg/browser"
)

var PREFIX string = "i:"

var ACTION_INVALID = 0
var ACTION_OPEN_WINDOW int8 = 1
var ACTION_DIALOG int8 = 2

type Instruction struct {
	Type       int8
	Parameters []string
}

func (inst *Instruction) Execute() error {
	switch inst.Type {
	case ACTION_OPEN_WINDOW:
		if len(inst.Parameters) >= 1 {
			return browser.OpenURL(inst.Parameters[0])
		}
		return fmt.Errorf("e: Not enough parameters to execute OPEN instruction")
	case ACTION_DIALOG:
		if len(inst.Parameters) >= 2 {
			_, err := dlgs.MessageBox(inst.Parameters[0], inst.Parameters[1])
			return err
		}
		return fmt.Errorf("e: Not enough parameters to execute DIALOG instruction")
	}
	return fmt.Errorf("e: Did not execute any action; unrecognized Instruction type")
}

func ParseInstructions(toParse string) ([]*Instruction, error) {
	var res []*Instruction = make([]*Instruction, 0)
	var insts []string = strings.Split(strings.ReplaceAll(toParse, "\n", ""), ";")

	for _, v := range insts {
		toParse = v[len(PREFIX):]
		open_bracket_index := strings.Index(v, "[")
		closing_bracket_index := strings.LastIndex(v, "]")
		if open_bracket_index != -1 || closing_bracket_index != -1 {
			parameters := strings.Split(toParse[open_bracket_index-1:closing_bracket_index-2], "<|>")
			toParse = toParse[0 : open_bracket_index-2]
			var inst_type int8
			switch strings.ToLower(toParse) {
			case "open":
				inst_type = ACTION_OPEN_WINDOW
			case "dialog":
				inst_type = ACTION_DIALOG
			}
			res = append(res, &Instruction{
				Type:       inst_type,
				Parameters: parameters,
			})
		}
	}

	return res, nil
}
