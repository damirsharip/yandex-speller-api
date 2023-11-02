package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

func main() {
	// читаем из файла
	body, err := os.ReadFile("file.json")
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	// анмаршлим в слайс структуры pair
	var pairs []pair
	err = json.Unmarshal(body, &pairs)
	if err != nil {
		log.Fatalf("failed to unmarshal JSON: %v", err)
	}

	// записываем кол-во горутин
	fmt.Print("Количество горутин: ")
	var numOfRoutines int
	_, err = fmt.Scan(&numOfRoutines)
	if err != nil {
		log.Fatalf("failed to read number of goroutines: %v", err)
	}

	// Метод 1. Использование буферизированного канала
	totalSumBuffered := SumWithBufferedChan(numOfRoutines, pairs)
	fmt.Printf("Общая сумма чисел c буферизированным каналом: %d\n", totalSumBuffered)

	// Метод 2. Использование небуферизированного канала
	totalSumUnbuffered := SumWithUnBufferedChan(numOfRoutines, pairs)
	fmt.Printf("Общая сумма чисел с не буферизированным каналом: %d\n", totalSumUnbuffered)
}

type pair struct {
	A int `json:"a"`
	B int `json:"b"`
}

func SumWithBufferedChan(numOfRoutines int, pairs []pair) int {
	var totalSum int

	// делим общее кол-во элементов на горутины
	// в случае если у нас много горутин чем элементов, ограничиваем кол-во горутин на len(pairs)
	// partsize это количество элементов которая каждая горутина должна вычитать в свой sum чтобы отправить в канал
	partSize := len(pairs) / numOfRoutines
	if partSize == 0 {
		partSize = 1
		numOfRoutines = len(pairs)
	}

	var mu sync.Mutex // создаем мьютекс чтобы асинхронно обновлять totalsum

	// создаем буферизированным канал с кол-во горутин, от которых в дальнейшем будем брать данные
	resultChan := make(chan int, numOfRoutines)

	// создаем waitgroup, чтобы потом ждать окончания всех записей от горутин в канал
	var wg sync.WaitGroup
	for i := 0; i < numOfRoutines; i++ {
		wg.Add(1)
		start := i * partSize
		end := (i + 1) * partSize

		if i == numOfRoutines-1 { // если разделение было не целым числом, то в последнем круге end может быть меньше len(pairs), из за этого его end = len(pairs)
			end = len(pairs)
		}

		go func(start, end int) {
			defer wg.Done()
			sum := 0
			for j := start; j < end; j++ {
				sum += pairs[j].A + pairs[j].B
			}
			resultChan <- sum
			mu.Lock()       // Заблокируйем мьютекс перед обновлением totalSum
			totalSum += sum // меняем totalSum
			mu.Unlock()     // разблокируем чтобы другие горутины могли писать
		}(start, end)
	}
	// ждем и закрываем канал
	wg.Wait()
	close(resultChan)

	// возвращаем ответ
	return totalSum
}

func SumWithUnBufferedChan(numOfRoutines int, pairs []pair) int {
	var totalSum int
	resultChan := make(chan int)
	// это сигнальный канал который информирует нас о том что все горутины закончили
	doneChan := make(chan struct{})

	// делим общее кол-во элементов на горутины
	// в случае если у нас много горутин чем элементов, ограничиваем кол-во горутин на len(pairs)
	// partsize это количество элементов которая каждая горутина должна вычитать в свой sum чтобы отправить в канал1
	partSize := len(pairs) / numOfRoutines
	if partSize == 0 {
		partSize = 1
		numOfRoutines = len(pairs)
	}

	// создаем waitgroup, чтобы потом ждать окончания всех записей от горутин в канал
	var wg sync.WaitGroup
	for i := 0; i < numOfRoutines; i++ {
		wg.Add(1)
		start := i * partSize
		end := (i + 1) * partSize

		if i == numOfRoutines-1 { // если разделение было не целым числом, то в последнем круге end может быть меньше len(pairs), из за этого его end = len(pairs)
			end = len(pairs)
		}

		go func(start, end int) {
			defer wg.Done()
			sum := 0
			for j := start; j < end; j++ {
				sum += pairs[j].A + pairs[j].B
			}
			resultChan <- sum
		}(start, end)
	}

	// ждем в отдельной горутин окончания всех считывающих горутин и затем закрываем канал doneChan, информируя о том что все закончилось и нужно возвращать результат
	go func() {
		wg.Wait()
		close(doneChan)
		close(resultChan)
	}()

	// Используйем select для получения значений из небуферизованного канала
	for {
		select {
		case sum := <-resultChan:
			totalSum += sum
		case <-doneChan:
			// возвращаем ответ
			return totalSum
		}
	}
}
