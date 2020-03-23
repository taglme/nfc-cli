# Taglme CLI Application 

[![release](https://badgen.net/github/tag/taglme/nfc-cli)](https://github.com/taglme/nfc-cli/releases)

Cross-platform CLI for reading NFC tags 

## Usage

```bash
nfc-cli [global options] command [command options] [arguments...]
```

### Commands

- `adapters` - Get adapters list
- `dump` - Dump tag memory
- `format` - Lock tag memory
- `lock` - Lock tag memory
- `read` - Read tag data with NDEF message
- `rmpwd` - Remove password for tag write acccess
- `run` - Load jobs from file and send them to server
- `setpwd` - Remove password for tag write acccess
- `transmit` - Transmit bytes to adapter or tag
- `version` - Application version
- `write` - Write NDEF message to the tag
- `help`, `h` - Shows a list of commands or help for one command

### Global options

- `--host` - Target host and port 
- `--adapter` - Adapter

## Development

- `make build-windows` – Build .exe for Windows platform   
- `make build-mac` – Build executable file for Mac platform      
- `make build-linux` – Build executable file for Linux platform     
- `make test` – Execute tests   
- `make deps` – Download and store dependencies
- `make lint` – Execute code linting
