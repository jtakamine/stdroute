# stdroute
Routes Stdin to a specified endpoint via JSON-RPC

Example Usages:

    $ someProgram | stdroute -dest "www.example.com:8000/jsonrpc" -method "Log.Write"

The above will pipe someProgram's output to stdroute, which in turn will forward the output to www.example.com:8000/jsonrpc (1 Log.Write RPC call per line):

    { "method": "Log.Write", "params": ["someProgram's output"], "id": 1}  -->  www.example.com:8000/jsonrpc
