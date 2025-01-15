package main

import (
	"bufio"
	"container/heap"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
)

type Point struct {
	x int
	y int
}

func CorrectPoint(p Point, matrix [][]int) bool {
	return 0 <= p.x && p.x < len(matrix[0]) && 0 <= p.y && p.y < len(matrix) && matrix[p.y][p.x] != 0
}

func ReadInput(input io.Reader) ([][]int, Point, Point, error) {
	var length, width int
	var start, finish Point
	reader := bufio.NewReader(input)

	_, err := fmt.Fscan(reader, &length, &width)
	if err != nil || length == 0 || width == 0 {
		return nil, start, finish, errors.New("incorrect dimensions of the maze")
	}

	matrix := make([][]int, length)
	for i := range matrix {
		matrix[i] = make([]int, width)
		for j := range matrix[i] {
			_, err = fmt.Fscan(reader, &matrix[i][j])
			if err != nil || matrix[i][j] < 0 || 9 < matrix[i][j] {
				return nil, start, finish, errors.New("incorrect maze structure")
			}
		}
	}

	_, err = fmt.Fscan(reader, &start.y, &start.x, &finish.y, &finish.x)
	if err != nil || !CorrectPoint(start, matrix) || !CorrectPoint(finish, matrix) || start == finish {
		return nil, start, finish, errors.New("incorrect data about the beginning or end of the way")
	}

	return matrix, start, finish, nil
}

func PavingWayToFinish(distances [][]int, visited [][]bool, matrix [][]int, start, finish Point) error {
	pq := &PriorityQueue{}
	heap.Push(pq, &PointCost{point: Point{start.x, start.y}, cost: 0})
	directions := [4][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	for pq.Len() > 0 {
		currPointCost := heap.Pop(pq).(*PointCost)
		if visited[currPointCost.point.y][currPointCost.point.x] {
			continue
		}
		visited[currPointCost.point.y][currPointCost.point.x] = true

		for _, direction := range directions {
			x := currPointCost.point.x + direction[0]
			y := currPointCost.point.y + direction[1]
			if CorrectPoint(Point{x, y}, matrix) {
				newCost := currPointCost.cost + matrix[y][x]
				if newCost < distances[y][x] {
					distances[y][x] = newCost
					heap.Push(pq, &PointCost{point: Point{x, y}, cost: newCost})
				}
			}
		}
	}
	if distances[finish.y][finish.x] == math.MaxInt {
		return errors.New("the way was not found")
	}
	return nil
}

func GetWay(distances [][]int, start, finish Point) []Point {
	ways := []Point{finish}
	directions := [4][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	for {
		minCost := math.MaxInt
		var way Point
		for _, value := range directions {
			x, y := finish.x+value[0], finish.y+value[1]
			point := Point{x, y}
			if point == start {
				return append(ways, point)
			}
			if CorrectPoint(point, distances) && distances[y][x] < minCost {
				minCost = distances[y][x]
				way = point
			}
		}
		ways = append(ways, way)
		finish = way
	}
}

func FindingFastestWay(matrix [][]int, start, finish Point) ([]Point, error) {
	distances := make([][]int, len(matrix))
	visited := make([][]bool, len(matrix))
	for i := range distances {
		distances[i] = make([]int, len(matrix[0]))
		visited[i] = make([]bool, len(matrix[0]))
		for j := range distances[i] {
			distances[i][j] = math.MaxInt
		}
	}
	distances[start.y][start.x] = matrix[start.y][start.x]
	err := PavingWayToFinish(distances, visited, matrix, start, finish)
	if err != nil {
		return nil, err
	}
	return GetWay(distances, start, finish), nil
}

func main() {
	matrix, start, finish, err := ReadInput(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	way, err := FindingFastestWay(matrix, start, finish)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for i := len(way) - 1; i >= 0; i-- {
		fmt.Println(way[i].y, way[i].x)
	}
	fmt.Println(".")
}
