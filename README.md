# stdroute
Routes Stdin to a specified endpoint via JSON-RPC

Example Usages:

    $ someProgram | stdroute "www.example.com:8000/jsonrpc"
    
or

    $ export STDROUTE_DEST=www.example.com:8000/jsonrpc
    $ someProgram | stdroute

Both of the above will pipe someProgram's output to stdroute, which in turn forwards the stream to www.example.com:8000/jsonrpc (1 RPC call per line):

    { "method": "write", "params": ["someProgram's output"], "id": 1}  -->  www.example.com:8000/jsonrpc
