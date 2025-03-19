# v6d-fuse

A FUSE filesystem implementation for accelerate object storage accessing, using [vineyard](https://v6d.io/docs.html) as a cache layer

## Features

- Mount a FUSE filesystem to read object storage
- Readonly for object storage
- Distributed-caching mechanism for improved performance

## Prerequisites

- Go 1.21 or later
- FUSE system requirements:
  - Linux: `fuse` package installed
  - MacOS: haven't been tested

## Quick Start

```bash
# Mount v6d filesystem
v6d-fuse <mountpoint>

# Unmount
fusermount -u <mountpoint>
```

## Status

This project is currently in early development stage. Implementing basic functionalities.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 