go run server.go

# GOMAXPROCS=1 GODEBUG=schedtrace=1000 View the scheduler status
# GODEBUG='gctrace=1' Output trace information related to the GC each time the GC is executed.