[![progress-banner](https://backend.codecrafters.io/progress/shell/d9f7394a-7cab-496f-b7c5-c9076d532c12)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)

This is a starting point for Go solutions to the
["Build Your Own Shell" Challenge](https://app.codecrafters.io/courses/shell/overview).

# GoShell

GoShell is a custom POSIX-compliant shell written in Go. Built as part of the CodeCrafters Shell Challenge, GoShell provides an interactive command-line environment that can interpret shell commands, run external programs, and handle built-in commands with a modern, efficient design.

## Overview

GoShell is designed to give users a familiar Unix-like command-line interface. Along the way, the project explores low-level terminal operations, custom command parsing, and advanced features like autocomplete. The development process deepened my understanding of terminal raw mode, and data structures in Go.
**Usage:** To run the shell, simply execute the provided `run.sh` script:

```sh
./run.sh
```

## Features

### Core Capabilities

- **Command Execution**: Run external programs and capture their output.
- **Input/Output Redirection**: Support for `>`, `>>`, and `<` operators.
- **Autocompletion**: autocomplete commands with `\t`.

### Built-in Commands

- `exit`: Terminate the shell.
- `echo`: Display text to stdout.
- `type`: Show command information.
- `pwd`: Print current working directory.
- `cd`: Change the current directory.

### Interactive Enhancements

- **Raw Mode Terminal**: Direct control over terminal input by manually putting it in raw mode.
- **Trie-Based Autocomplete**: Efficiently suggest completions for commands and file names.
- **Command History Navigation**: Browse and reuse previous commands.
- **Dynamic Cursor Control**: Real-time handling of cursor positions, insertions, and key events.

### Advanced Parsing

- **Quote and Escape Handling**: Parse single and double quotes, along with escaped characters.
- **Token Recognition**: Break down input into meaningful commands and arguments.
- **Redirection Parsing**: Detect and handle redirection operators.

## Implementation Details

### Terminal Management

GoShell manually sets the terminal into raw mode, allowing it to process each keystroke individually. This enables:

- Fine-grained control of input/output.
- Custom key bindings.
- Immediate display updates and precise cursor management.
- Signal handling (e.g., interrupts via Ctrl+C).

### Command Parser

The shell’s parser tokenizes user input and handles:

- Environment variable and home directory expansion.
- Special tokens (quotes, redirection symbols, etc.).
- Conversion of raw input into a structured format for execution.

### Autocomplete with Trie

The autocomplete system uses a Trie data structure:

- **Efficient Lookup**: Quickly suggests completions based on the current input prefix.
- **Memory Efficiency**: Stores command and file names in a compact format.
- **Seamless Integration**: Enhances the interactive experience by suggesting completions as you type.

### Process Management

GoShell launches external programs by:

- Creating and managing subprocesses.
- Setting up proper I/O redirection.
- Propagating signals and handling exit codes correctly.

## Future Enhancements

While GoShell currently supports many core features, future improvements may include:

- **Piping and Job Control**: Advanced features to handle pipelines and background tasks.
- **History Navigation**: More intuitive history search and manipulation.

## Acknowledgments

I would like to extend my gratitude to [CodeCrafters](https://app.codecrafters.io/catalog) for designing this challenge. It pushed me to explore deep system-level programming in Go, improved my skills significantly, and sparked new ideas for future projects.
