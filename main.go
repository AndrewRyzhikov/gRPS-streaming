package main

import (
	"fmt"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"log"
	"net"
	"time"

	//hardwaremonitoring "github.com/percybolmer/grpcstreams/proto"
	hardwaremonitoring "streaming/internal/pb"
	// Dont forget this import :)
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Welcome to streaming HW monitoring")
	// Setup a tcp connection to port 7777
	lis, err := net.Listen("tcp", ":7777")
	if err != nil {
		panic(err)
	}
	// Create a gRPC server
	gRPCserver := grpc.NewServer()

	// Create a server object of the type we created in server.go
	s := &Server{}

	// Regiser our server as a gRPC server
	hardwaremonitoring.RegisterHardwareMonitorServer(gRPCserver, s)

	log.Println(gRPCserver.Serve(lis))

}

// Server is our struct that will handle the Hardware monitoring Logic
// It will fulfill the gRPC interface generated
type Server struct {
	hardwaremonitoring.UnimplementedHardwareMonitorServer
}

// Monitor is used to start a stream of HardwareStats
func (s *Server) Monitor(req *hardwaremonitoring.EmptyRequest, stream hardwaremonitoring.HardwareMonitor_MonitorServer) error {
	// Start a ticker that executes each 2 seconds
	timer := time.NewTicker(2 * time.Second)

	for {
		select {
		// Exit on stream context done
		case <-stream.Context().Done():
			fmt.Println("Client Done")
			return nil
		case <-timer.C:
			// Grab stats and output
			hwStats, err := s.GetStats()
			if err != nil {
				log.Println(err.Error())
			} else {

			}
			// Send the Hardware stats on the stream
			err = stream.Send(hwStats)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
}

// GetStats will extract system stats and output a Hardware Object, or an error
// if extraction fails
func (s *Server) GetStats() (*hardwaremonitoring.HardwareStats, error) {
	// Extarcyt Memory statas
	mem, err := memory.Get()
	if err != nil {
		return nil, err
	}
	// Extract CPU stats
	cpu, err := cpu.Get()
	if err != nil {
		return nil, err
	}
	// Create our response object
	hwStats := &hardwaremonitoring.HardwareStats{
		Cpu:        int32(cpu.Total),
		MemoryFree: int32(mem.Free),
		MemoryUsed: int32(mem.Used),
	}

	return hwStats, nil
}
