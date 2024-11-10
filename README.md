# RINEX File Processing Project

A tool to process RINEX (Receiver Independent Exchange Format) files and output data in an easily consumable format.

## Installation Requirements

### Cap'n Proto Setup
1. Install Cap'n Proto compiler
2. Install go-capnp:
```bash
go get capnproto.org/go/capnp/v3
```
3. Add Go bin to PATH:
```bash
export PATH=$PATH:$HOME/go/bin
```

### Windows Users Note
Follow Cap'n Proto Go installation steps only up to step 3. The Cap'n Proto compiler works in conjunction with go-capnp to compile schema files into Go code.

## Schema Development
1. Create schema file (.capnp)
2. Compile schema:
```bash
capnp compile -I /path/to/go-capnp/std -ogo /path/to/yourschema.capnp
```

## Project Status

### Completed
- Initial schema creation
- Generated capnp.go file

### Current Focus
- RINEX file processing implementation
- Data structure optimization
- Schema refinement for RINEX data representation

### Future Development
- Support for various RINEX formats
- Data output standardization
- Performance optimization

## Documentation
For more detailed information, refer to the [Cap'n Proto documentation](https://capnproto.org/documentation.html).
