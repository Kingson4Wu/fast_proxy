# in proxy pprof
go tool pprof -http :8081 http://localhost:8033/debug/pprof/profile
go tool pprof -http :8081 http://localhost:8033/debug/pprof/heap

# out proxy pprof
go tool pprof -http :8082 http://localhost:8034/debug/pprof/profile
go tool pprof -http :8082 http://localhost:8034/debug/pprof/heap

