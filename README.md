# Sqlettus

**Redis-like** service with SQLite as a backend. A playful fusion of SQLite and
Redis to offer in-memory data structure store capabilities with the persistent
storage of SQLite.

## Features

- Redis-compatible command support
- SQLite backend for data persistence
- Extensible command routing system

## Supported Commands (so far)

- `COMMAND DOCS`
- `CONFIG GET`
  - `save`
  - `appendonly`
- `FLUSHALL`
- `PING`
- `SET`
- `GET`

For a detailed list and updates on commands, see the handler package in the
code.

## Installation

As this is a Go-based project, you'd typically clone the repository and build it
using `go build`. However, specific installation instructions will depend on any
further setup you provide. For now:

```bash
git clone https://github.com/jtarchie/sqlettus.git
cd sqlettus
brew bundle
task
```

## Usage

After building, you can run the executable:

```bash
./sqlettus
```

For actual usage within applications, you'd likely use it similar to how you'd
use Redis, but you'll need to provide more specifics or examples.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to
discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)
