package altsrc

import (
	"flag"

	"gopkg.in/urfave/cli.v1"
)

// WARNING: This file is generated!

// BoolFlag is the flag type that wraps cli.BoolFlag to allow
// for other values to be specified
type BoolFlag struct {
	cli.BoolFlag
	set *flag.FlagSet
}

// NewBoolFlag creates a new BoolFlag
func NewBoolFlag(fl cli.BoolFlag) *BoolFlag {
	return &BoolFlag{BoolFlag: fl, set: nil}
}

// Apply saves the flagSet for later usage calls, then calls the
// wrapped BoolFlag.Apply
func (f *BoolFlag) Apply(set *flag.FlagSet) {
	f.set = set
	f.BoolFlag.Apply(set)
}

// ApplyWithError saves the flagSet for later usage calls, then calls the
// wrapped BoolFlag.ApplyWithError
func (f *BoolFlag) ApplyWithError(set *flag.FlagSet) error {
	f.set = set
	return f.BoolFlag.ApplyWithError(set)
}

// BoolTFlag is the flag type that wraps cli.BoolTFlag to allow
// for other values to be specified
type BoolTFlag struct {
	cli.BoolTFlag
	set *flag.FlagSet
}

// NewBoolTFlag creates a new BoolTFlag
func NewBoolTFlag(fl cli.BoolTFlag) *BoolTFlag {
	return &BoolTFlag{BoolTFlag: fl, set: nil}
}

// Apply saves the flagSet for later usage calls, then calls the
// wrapped BoolTFlag.Apply
func (f *BoolTFlag) Apply(set *flag.FlagSet)