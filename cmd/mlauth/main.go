package main

import (
	"log"
	"mlauth/pkg/api"
	"mlauth/pkg/rpc"
	"net"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		r := api.Route()
		err := r.Run()
		if err != nil {
			log.Fatalln("Failed to start API server:", err.Error())
		}
	}()

	wg.Add(1)
	go func() {
		l, err := net.Listen("tcp", ":8001")
		if err != nil {
			log.Fatalln("Failed to listen:", err.Error())
		}

		s := rpc.Register()
		err = s.Serve(l)
		if err != nil {
			log.Fatalln("Failed to start RPC server:", err.Error())
		}
	}()

	wg.Wait()
}
