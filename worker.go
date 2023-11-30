package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"server/stubs"
)

type Worker struct{}

func MakeMatrix(height, width int) [][]uint8 {
	matrix := make([][]uint8, height)
	for i := range matrix {
		matrix[i] = make([]uint8, width)
	}
	return matrix
}

func printBoard(world func(y, x int) uint8, sx int, sz int) {
	for i := 0; i < sx; i++ {
		for j := 0; j < sz; j++ {
			fmt.Printf("%3d", world(i, j))
		}
		fmt.Println()
	}
}
func mprintBoard(world [][]byte, sx int, sz int) {
	for i := 0; i < sx; i++ {
		for j := 0; j < sz; j++ {
			fmt.Printf("%3d", world[i][j])
		}
		fmt.Println()
	}
}
func CalculateNextWorld(req stubs.Request, startX int, endX int) [][]uint8 {
	world := MakeImmutableMatrix(req.World)

	startY := req.StartY
	fmt.Println("This is start y:", startY)

	height := req.EndY - req.StartY
	width := req.Params.Width

	type cell struct {
		x, y int
	}

	nextWorld := MakeMatrix(height, width)
	for i := req.StartY; i < req.EndY; i++ {
		fmt.Println("THIS IS ROW : ", i)
		for j := 0; j < endX; j++ {
			neighbours := [8]cell{{(i - 1 + req.Params.Height) % req.Params.Height, (j + width) % width},
				{(i + 1 + req.Params.Height) % req.Params.Height, (j + width) % width},
				{(i + req.Params.Height) % req.Params.Height, (j - 1 + width) % width},
				{(i + req.Params.Height) % req.Params.Height, (j + 1 + width) % width},
				{(i + 1 + req.Params.Height) % req.Params.Height, (j - 1 + width) % width},
				{(i + 1 + req.Params.Height) % req.Params.Height, (j + 1 + width) % width},
				{(i - 1 + req.Params.Height) % req.Params.Height, (j - 1 + width) % width},
				{(i - 1 + req.Params.Height) % req.Params.Height, (j + 1 + width) % width}}

			live := 0 // intialised to 0
			for _, cell := range neighbours {

				if world(cell.x, cell.y) != 0 {
					live++
				}

			}

			if world(i, j) == 255 {
				if live < 2 {
					nextWorld[i-startY][j] = 0
				} else if live == 2 || live == 3 {
					nextWorld[i-startY][j] = 255
				} else {
					nextWorld[i-startY][j] = 0
				}
			} else {
				if live == 3 {
					nextWorld[i-startY][j] = 255
				} else {
					nextWorld[i-startY][j] = 0
				}
			}
		}

	}
	return nextWorld

}
func MakeImmutableMatrix(matrix [][]uint8) func(y, x int) uint8 {
	return func(y, x int) uint8 {
		return matrix[y][x]
	}
}

func (s *Worker) WorkerProcess(req stubs.Request, res *stubs.Response) (err error) {

	fmt.Println(req.EndY)
	newData := CalculateNextWorld(req, 0, req.Params.Width)
	fmt.Println("In worker : ", req.StartY)
	//printBoard(newData, req.EndY-req.StartY, req.Params.Width)
	res.NextWorld = newData
	fmt.Println("TRUE THAT I am done")

	return
}

func main() {
	pAddr := flag.String("port", "8031", "Port to listen on")
	flag.Parse()
	rpc.Register(&Worker{})
	listener, _ := net.Listen("tcp", ":"+*pAddr)
	defer listener.Close()
	rpc.Accept(listener)
}
