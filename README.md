# Sqlettuce

**Redis-like** service with SQLite as a backend. A playful fusion of SQLite and
Redis offers in-memory data structure store capabilities with the persistent
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

```bash
git clone https://github.com/jtarchie/sqlettuce.git
cd sqlettuce
brew bundle
task
```

## Usage

After building, you can run the executable:

```bash
./sqlettuce
```

## Contributing

Pull requests are welcome. For significant changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)
