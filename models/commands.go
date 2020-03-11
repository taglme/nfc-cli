package models

type Command = string

const (
	CommandVersion  Command = "version"
	CommandAdapters Command = "adapters"
	CommandWrite    Command = "write"
	CommandRead     Command = "read"
	CommandDump     Command = "dump"
	CommandTransmit Command = "transmit"
	CommandLock     Command = "lock"
	CommandSetpwd   Command = "setpwd"
	CommandRmpwd    Command = "rmpwd"
	CommandFormat   Command = "format"
	CommandRun      Command = "run"
)
