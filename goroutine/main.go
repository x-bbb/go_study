package main

import (
	"fmt"
	"math/rand"
)

type Job struct {
	Id     int
	Number int
}

type Result struct {
	Job *Job
	Sum int
}

func calc(job *Job, resultChan chan *Result) {
	number := job.Number
	var sum int
	for number != 0 {
		tmp := number % 10
		sum += tmp
		number /= 10
	}
	result := &Result{
		Job: job,
		Sum: sum,
	}

	resultChan <- result
}

func Worker(jobChan chan *Job, resultChan chan *Result) {
	for job := range jobChan {
		calc(job, resultChan)
	}
}

func StartWorkPool(num int, jobChan chan *Job, resultChan chan *Result) {
	for i := 0; i < num; i++ {
		go Worker(jobChan, resultChan)
	}
}
func PrintResult(resultChan chan *Result) {
	for result := range resultChan {
		fmt.Printf("Id:%d Number:%d Sum:%d\n", result.Job.Id, result.Job.Number, result.Sum)
	}
}

func main() {
	jobChan := make(chan *Job, 1000)
	resultChan := make(chan *Result, 1000)
	StartWorkPool(128, jobChan, resultChan)
	go PrintResult(resultChan)
	var id int
	for {
		id++
		number := rand.Int()
		job := &Job{
			Id:     id,
			Number: number,
		}
		jobChan <- job
	}

}
