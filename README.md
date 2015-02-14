# stdroute
Routes Stdin to a specified endpoint via JSON-RPC

Example Usages:

    $ someProgram | stdroute "www.example.com:8000/jsonrpc"
    
or

    $ export STDROUTE_DEST=www.example.com:8000/jsonrpc
    $ someProgram | stdroute

Both of the above will pipe someProgram's output to stdroute, which in turn will forward the output to www.example.com:8000/jsonrpc (1 RPC call per line):

    { "method": "Stdin.Write", "params": ["someProgram's output"], "id": 1}  -->  www.example.com:8000/jsonrpc
