package options

import (
	"os"

	"github.com/docopt/docopt-go"
	"github.com/mattn/go-isatty"
)

//nolint:gochecknoglobals
var (
	args        []string
	command     string
	interactive bool
	monitor     bool
	script      string
	terminal    int
	version     bool

	usage = `penishell

Usage:
  penishell [-m] SCRIPT [ARGUMENTS...]
  penishell [-m] -c COMMAND [NAME [ARGUMENTS...]]
  penishell [-im] [-s [ARGUMENTS...]]
  penishell -h
  penishell -v

Arguments:
  ARGUMENTS  Positional parameters.          
  SCRIPT     Path to penishell script. Also used as the value for $0.
  NAME       Override $0. Otherwise, $0 is set to name used to invoke penishell.

Options:
  -c, --command=COMMAND  Run the specified command.
  -m, --monitor          Invert job control mode.
  -i, --interactive      Disable interactive mode.
  -s, --stdin            Read commands from stdin.
  -h, --help             Display this help.
  -v, --version          Print penishell version.

If penishell's stdin is a TTY, and penishell was invoked with no non-option operands or
penishell was explicitly directed to evaluate commands from stdin, interactive and
job control features are enabled. Otherwise, these features are disabled.
`
)

// Args returns positional parameters (if any).
func Args() []string {
	return args
}

// Command returns the command specified (if any).
func Command() string {
	return command
}

// Interactive returns true if penishell should run in interactive mode.
func Interactive() bool {
	return interactive
}

// Parse parses the command line options for this invocation of penishell.
func Parse() {
	docopt.DefaultParser.OptionsFirst = true

	opts, err := docopt.ParseDoc(usage)
	if err != nil {
		// Error in the usage doc. This should never happen.
		panic(err.Error())
	}

	script = ""

	command, _ = opts.String("--command")

	name, _ := opts.String("NAME")
	if name == "" {
		name = os.Args[0]
	}

	path, _ := opts.String("SCRIPT")
	if path != "" {
		command = "source " + path
		name = path
		script = path
	} else if command == "" && isatty.IsTerminal(os.Stdin.Fd()) {
		interactive = true
		monitor = true
		terminal = int(os.Stdin.Fd())
	}

	args, _ = opts["ARGUMENTS"].([]string)
	args = append([]string{name}, args...)

	invertInteractive, _ := opts.Bool("--interactive")
	interactive = interactive != invertInteractive

	invertMonitor, _ := opts.Bool("--monitor")
	monitor = monitor != invertMonitor

	version, _ = opts.Bool("--version")
}

// Script returns the script name (if any).
func Script() string {
	return script
}

// Monitor returns true if job control features should be enabled.
func Monitor() bool {
	return monitor
}

// Terminal returns the terminal's integer file descriptor.
func Terminal() int {
	return terminal
}

// Version returns true if penishell's version was request.
func Version() bool {
	return version
}
