# Distributed Failure Detection

Architecture:

- Central node opens a port for registration using a seperate thread.
- Worker nodes can connect to central node and register themselves in the topology using this connection.
- Central node requests for liveness from every registered worker node in the topology one by one.
- Worker node receives the request and announces the liveness to central node.
- If a worker node does not respond to central nodes request, then central node keeps track of it.
- If worker node does not respond to the central node's request for more than 3 times then it is deregistered from the topology, worker node needs to reregister in topology when it's back online
- This process is repeated by central node every 5 seconds
