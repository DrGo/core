package args

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const errWrongCommand = "must specify a valid command. For a list of commands, type %s help"

//Options stores package run options. If needed, must be set before calling any other package funcs
type Options struct {
	// help is the function called when an user requests help.
	Help func()
	//noSuchCommand is the function when the user enters the wrong command
	// NoSuchCommand func()
}

var opts *Options

//SetOptions sets package run options. If needed, must be called before calling any other package funcs
func SetOptions(options *Options) {
	opts = options
}

//Flag holds command line flag info
type Flag struct {
	dest   interface{}
	name   string
	letter string
	value  interface{}
}

// NewFlag returns a new Flag struct
func NewFlag(destination interface{}, name, letter string) Flag {
	return Flag{
		dest:   destination,
		name:   name,
		letter: letter,
	}
}

//NewCommand creates a new command
func NewCommand(name string, args []Flag) *flag.FlagSet {
	cmd := flag.NewFlagSet(name, flag.ContinueOnError)
	//add duplicate args with both name and letter specified, so the command
	//can be invoked by either
	var expanded []Flag
	for _, arg := range args {
		expanded = append(expanded, arg)
		if arg.letter != "" {
			expanded = append(expanded, arg)
			expanded[len(expanded)-1].name = arg.letter
		}
	}
	for _, arg := range expanded {
		if arg.value == nil {
			arg.value = arg.dest
		}
		switch p := arg.dest.(type) {
		case *string:
			cmd.StringVar(p, arg.name, *p, "")
		case *bool:
			cmd.BoolVar(p, arg.name, *p, "")
		case *int:
			cmd.IntVar(p, arg.name, *p, "")
		default:
			continue
		}
	}
	return cmd
}

// ParseCommandLine parses command line arguments for the appropriate subcommandparses arguments.
// The first command is the default command and can be nil.
func ParseCommandLine(top *flag.FlagSet, subs ...*flag.FlagSet) (*flag.FlagSet, error) {
	exeName := os.Args[0]
	flg, err := ParseArguments(os.Args[1:], top, subs...)
	if err != nil {
		s := err.Error()
		switch {
		case strings.Contains(s, "flag provided but not defined"):
			s = strings.Replace(s, "provided but not defined", "does not exist", 1)
			return nil, fmt.Errorf(s)
		case strings.Contains(s, "help requested"):
			if opts.Help != nil {
				opts.Help()
			}
			//FIXME: prevent treating this as an error
			// return nil, nil
		default:
			return nil, err
		}
	}
	if flg == nil || flg.Name() == "" {
		return nil, fmt.Errorf(errWrongCommand, exeName)
	}
	return flg, nil
}

//ParseArguments parses arguments (passed as a string array) for the appropriate subcommand
func ParseArguments(args []string, top *flag.FlagSet, subs ...*flag.FlagSet) (*flag.FlagSet, error) {
	if top == nil {
		top = flag.NewFlagSet("", flag.ContinueOnError)
	}
	if err := top.Parse(args); err != nil {
		return nil, err
	}
	args = top.Args()
	if len(args) == 0 || len(subs) == 0 { //nothing left to parse
		return top, nil
	}
	cmdTable := make(map[string]*flag.FlagSet)
	for _, cmd := range subs {
		if cmd != nil {
			cmdTable[cmd.Name()] = cmd
		}
	}
	flagSet, found := cmdTable[args[0]] //retrieve the FlagSet for this subcommand
	if !found {
		return nil, fmt.Errorf("command %v is not found", args[0])
	}
	if len(args) == 1 { //nothing left to parse
		return flagSet, nil
	}
	args = args[1:] //skip over the subcommand name
	//move (positional) arguments to their own array
	posArgs := []string{}
	for len(args[0]) > 1 && args[0][0] != '-' { //loop while the first argument is not a flag
		posArgs = append(posArgs, args[0]) //add it to the positional
		//skip to the next arg if any
		if len(args) == 1 {
			break
		}
		args = args[1:]

	}
	//parse the flags
	if err := flagSet.Parse(args); err != nil {
		return nil, err
	}
	//parse the positional arguments
	if err := flagSet.Parse(posArgs); err != nil {
		return nil, err
	}
	return flagSet, nil
}
