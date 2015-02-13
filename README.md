# stdroute
Routes Stdin to a specified endpoint via JSON-RPC

Example Usage:

    someProgram | stdroute "www.example.com:8000/jsonrpc"

The above command will pipe someProgram's output to stdroute, which will in turn forward the stream to www.example.com:8000/jsonrpc (1 RPC call per line):

    --> { "method": "write", "params": ["someProgram's output"], "id": 1}
