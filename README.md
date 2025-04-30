# HTTP Server from Scratch

A simple HTTP server implemented in Go from scratch for learning purpose.

## Features

- **No net/http** ✅
- **HTTP Request Parsing** ✅
- **Static File Serving** ✅
- **Dynamic Routes** ✅
- **File Uploads** ✅
- **Gzip Compression** ✅
- **Connection Handling** ✅

## Getting Started

### Prerequisites

- Go 1.24 or later installed on your system.

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/codecrafters-io/http-server-starter-go.git
   cd http-from-scratch
   ```

2. Build the project:
   ```bash
   go build -o http-server ./app/main.go
   ```

### Usage

Run the server with the following command:

```bash
./http-server -directory <path_to_static_files>
```


## License

This project is released under the [Unlicense](https://unlicense.org/).
