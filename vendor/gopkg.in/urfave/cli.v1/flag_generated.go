package cli

import (
	"flag"
	"strconv"
	"time"
)

// WARNING: This file is generated!

// BoolFlag is a flag with type bool
type BoolFlag struct {
	Name        string
	Usage       string
	EnvVar      string
	Hidden      bool
	Destination *bool
}

// String returns a readable representation of this value
// (for usage defaults)
func (f BoolFlag) String() string {
	return FlagStringer(f)
}

// GetName returns the name of the flag
func (f BoolFlag) GetName() string {
	return f.Name
}

// Bool looks up the value of a local BoolFlag, returns
// false if not found
func (c *Context) Bool(name string) bool {
	return lookupBool(name, c.flagSet)
}

// GlobalBool looks up the value of a global BoolFlag, returns
// false if not found
func (c *Context) GlobalBool(name string) bool {
	if fs := lookupGlobalFlagSet(name, c); fs != nil {
		return lookupBool(name, fs)
	}
	return false
}

func lookupBool(name string, set *flag.FlagSet) bool {
	f := set.Lookup(name)
	if f != nil {
		parsed, err := strconv.ParseBool(f.Value.String())
		if err != nil {
			return false
		}
		return parsed
	}
	return false
}

// BoolTFlag is a flag with type bool that is true by default
type BoolTFlag struct {
	Name        string
	Usage       string
	EnvVar      string
	Hidden      bool
	Destination *bool
}

// String returns a readable representation of this value
// (for usage defaults)
func (f BoolTFlag) String() string {
	return FlagStringer(f)
}

// GetName returns the name of the flag
func (f BoolTFlag) GetName() string {
	return f.Name
}

// BoolT looks up the value of a local BoolTFlag, returns
// false if not found
func (c *Context) BoolT(name string) bool {
	return lookupBoolT(name, c.flagSet)
}

// GlobalBoolT looks up the value of a global BoolTFlag, returns
// false if not found
func (c *Context) GlobalBoolT(name string) bool {
	if fs := lookupGlobalFlagSet(name, c); fs != nil {
		return lookupBoolT(name, fs)
	}
	return false
}

func lookupBoolT(name string, set *flag.FlagSet) bool {
	f := set.Lookup(name)
	if f != nil {
		parsed, err := strconv.ParseBool(f.Value.String())
		if err != nil {
			return false
		}
		return parsed
	}
	return false
}

// DurationFlag is a flag with type time.Duration (see https://golang.org/pkg/time/#ParseDuration)
type DurationFlag struct {
	Name        string
	Usage       string
	EnvVar      string
	Hidden      bool
	Value       time.Duration
	Destination *time.Duration
}

// String returns a readable representation of this value
// (for usage defaults)
func (f DurationFlag) String() string {
	return FlagStringer(f)
}

// GetName returns the name of the flag
func (f DurationFlag) GetName() string {
	return f.Name
}

// Duration looks up the value of a local DurationFlag, returns
// 0 if not found
func (c *Context) Duration(name string) time.Duration {
	return lookupDuration(name, c.flagSet)
}

// GlobalDuration looks up the value of a global DurationFlag, returns
// 0 if not found
func (c *Context) GlobalDuration(name string) time.Duration {
	if fs := lookupGlobalFlagSet(name, c); fs != nil {
		return lookupDuration(name, fs)
	}
	return 0
}

func lookupDuration(name string, set *flag.FlagSet) time.Duration {
	f := set.Lookup(name)
	if f != nil {
		parsed, err := time.ParseDuration(f.Value.String())
		if err != nil {
			return 0
		}
		return parsed
	}
	return 0
}

// Float64Flag is a flag with type float64
type Float64Flag struct {
	Name        string
	Usage       string
	EnvVar      string
	Hidden      bool
	Value       float64
	Destination *float64
}

// String returns a readable representation of this value
// (for usage defaults)
func (f Float64Flag) String() string {
	return FlagStringer(f)
}

// GetName returns the name of the flag
func (f Float64Flag) GetName() string {
	return f.Name
}

// Float64 looks up the value of a local Float64Flag, returns
// 0 if not found
func (c *Context) Float64(name string) float64 {
	return lookupFloat64(name, c.flagSet)
}

// GlobalFloat64 looks up the value of a global Float64Flag, returns
// 0 if not found
func (c *Context) GlobalFloat64(name string) float64 {
	if fs := lookupGlobalFlagSet(name, c); fs != nil {
		return lookupFloat64(name, fs)
	}
	return 0
}

func lookupFloat64(name string, set *flag.FlagSet) float64 {
	f := set.Lookup(name)
	if f != nil {
		parsed, err := strconv.ParseFloat(f.Value.String(), 64)
		if err != nil {
			return 0
		}
		return parsed
	}
	return 0
}

// GenericFlag is a flag with type Generic
type GenericFlag struct {
	Name   string
	Usage  string
	EnvVar string
	Hidden bool
	Value  Generic
}

// String returns a readable representation of this value
// (for usage defaults)
func (f GenericFlag) String() string {
	return FlagStringer(f)
}

// GetName returns the name of the flag
func (f GenericFlag) GetName() string {
	re