package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"slices"
)

type PointCost struct {
	x    int
	y    int
	cost int
}

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

func PavingWayToFinish(distances [][]int, visited [][]bool, matrix [][]int, currPoint, finish Point) error {
	directions := [4][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	for {
		var points []PointCost
		for i := 0; i < 4; i++ {
			point := Point{currPoint.x + directions[i][0], currPoint.y + directions[i][1]}
			if CorrectPoint(point, matrix) && visited[point.y][point.x] == false {
				points = append(points, PointCost{point.x, point.y, int(matrix[point.y][point.x])})
			}
		}
		if len(points) == 0 {
			break
		}
		slices.SortFunc(points, func(i, j PointCost) int { return i.cost - j.cost })
		for _, point := range points {
			distances[point.y][point.x] = min(distances[point.y][point.x], distances[currPoint.y][currPoint.x]+int(matrix[point.y][point.x]))
			if point.x == finish.x && point.y == finish.y {
				return nil
			}
		}
		visited[currPoint.y][currPoint.x] = true
		currPoint = Point{points[0].x, points[0].y}
	}
	return errors.New("the way was not found")
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

func findingFastestWay(matrix [][]int, start, finish Point) error {
	distances := make([][]int, len(matrix))
	visited := make([][]bool, len(matrix))
	for i := range distances {
		distances[i] = make([]int, len(matrix[0]))
		visited[i] = make([]bool, len(matrix[0]))
		for j := range distances[i] {
			distances[i][j] = math.MaxInt
		}
	}
	distances[start.y][start.x] = int(matrix[start.y][start.x])
	err := PavingWayToFinish(distances, visited, matrix, start, finish)
	if err != nil {
		return err
	}
	ways := GetWay(distances, start, finish)
	for i := len(ways) - 1; i >= 0; i-- {
		fmt.Println(ways[i].y, ways[i].x)
	}
	fmt.Println(".")
	return nil
}

func main() {
	matrix, start, finish, err := ReadInput(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = findingFastestWay(matrix, start, finish)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
