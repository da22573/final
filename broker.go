package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"server/stubs"
	"strconv"
	"sync"
)

var WorldMutex sync.Mutex
var world [][]byte
var turn int

type GameOfLife struct{}

func MakeMatrix(height, width int) [][]uint8 {
	matrix := make([][]uint8, height)
	for i := range matrix {
		matrix[i] = make([]uint8, width)
	}
	return matrix
}

func (s *GameOfLife) ProcessAliveCells(req stubs.Request, res *stubs.Response) (err error) {

	res.InstantWorld = world
	res.CompletedTurn = turn

	return
}
func (s *GameOfLife) PauseGameProcess(req stubs.Request, res *stubs.Response) (err error) {

	WorldMutex.Lock()
	return
}
func (s *GameOfLife) UnPauseGameProcess(req stubs.Request, res *stubs.Response) (err error) {

	WorldMutex.Unlock()
	return
}

func (s *GameOfLife) ProcessAllTurns(req stubs.Request, res *stubs.Response) (err error) {
	fmt.Println("I am on the other end")
	turn = 1
	threads := 4
	world = req.World

	turns := req.Params.Turns

	var clients [4]*rpc.Client
	var responses [4]stubs.Response

	for turn <= turns {

		doneChannels := make([]chan *rpc.Call, 4)
		for i := 0; i < threads; i++ {
			doneChannels[i] = make(chan *rpc.Call, 1)
		}

		height := req.Params.Height / threads

		for i := range responses {
			responses[i] = stubs.Response{}
		}

		for index := range clients {

			request := stubs.Request{
				World: world,
				Params: stubs.Params{
					Turns:   req.Params.Turns,
					Threads: threads,
					Height:  req.Params.Height,
					Width:   req.Params.Width,
				},
				StartY: index * height,
				EndY:   (index + 1) * height,
			}
			clients[index], err = rpc.Dial("tcp", "127.0.0.1:803"+strconv.Itoa(index+1))

			if err != nil {
				fmt.Println("Error making RPC call", err)
			}
			if clients[index] != nil {
				//printBoard(request.World, 4, 16)
				clients[index].Go(stubs.WorkerProcess, request, &responses[index], doneChannels[index])
			} else {
				fmt.Println("Client is nil. Skipping client.Go call.")
			}

		}

		nextWorld := MakeMatrix(0, 0)
		for i := 0; i < threads; i++ {

			<-doneChannels[i]
			fmt.Println("YOU THINK I AM DONE")
			part := responses[i].NextWorld

			nextWorld = append(nextWorld, part...)

		}

		WorldMutex.Lock()
		world = nextWorld
		res.NextWorld = nextWorld
		turn++
		WorldMutex.Unlock()

	}
	return

}
func main() {
	pAddr := flag.String("port", "8030", "Port to listen on")
	flag.Parse()
	rpc.Register(&GameOfLife{})
	listener, _ := net.Listen("tcp", ":"+*pAddr)
	defer listener.Close()
	rpc.Accept(listener)
}
