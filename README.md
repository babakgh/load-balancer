I chose golang for this state but I might implement the C++ version later if there was enough time.

- Lets start simple
  we need a routine to genearate fake request. I am not going to actually create a http request, I will just create a struct that represents a request. The struct should have a string that will later be used as the partition key to find the corsponding backend.
  Also I am going to limit the partition keys and for each key we have a corresponding backend struct that will be used to handle the request.

So there will be a map of partition keys to backend structs.

Benchmarking

- When engine is doing nothing
  BenchmarkEngine-10 1000000000 0.0000002 ns/op
