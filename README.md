

## Writing and Compiling Cap'n Proto Schema for Go

1. Read the installation instructions from [Cap'n Proto Go](https://github.com/capnproto/go-capnp/blob/main/docs/Getting-Started.md#remote-calls-using-interfaces).
   **Note for Windows users:** Follow the installation steps only up to step 3.

2. Install the Cap'n Proto compiler.

3. Install go-capnp:
    ```
    go get capnproto.org/go/capnp/v3
    ```
4. Clone the go-capnp repository:
    ```
    git clone https://github.com/capnproto/go-capnp.git
    ```
5. Write your schema file (e.g., `yourschema.capnp`).

6. Compile the schema using the following command:
```
capnp compile -I /path/to/go-capnp/std -ogo /path/to/yourschema.capnp
```
For more detailed information refer to the [Cap'n Proto documentation](https://capnproto.org/index.html).


### This will mainly be a place to write thigs down so that I dont forget them.

schema been created , capnp.go file been created, need to update and properly use in structs 