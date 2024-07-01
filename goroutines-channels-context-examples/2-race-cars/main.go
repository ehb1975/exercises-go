package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type car struct {
	name  string
	tempo int
}

type corrida struct {
	voltas int
}

func main() {

	wg := sync.WaitGroup{}
	wg.Add(2)

	car1 := &car{name: "Ferrari", tempo: 0}
	car2 := &car{name: "Lamborghini", tempo: 0}

	fmt.Println("Race starts")

	// Create a Goroutine for each car

	go func() {
		race(car1)
		wg.Done()
	}()

	go func() {
		race(car2)
		wg.Done()
	}()

	wg.Wait()

	fmt.Println("Race over!")
	fmt.Println("Tempo", car1.name, car1.tempo)
	fmt.Println("Tempo", car2.name, car2.tempo)

	if car1.tempo < car2.tempo {
		fmt.Println(car1.name, "- Carro vencedor")
	} else if car1.tempo > car2.tempo {
		fmt.Println(car2.name, "- Carro vencedor")
	} else {
		fmt.Println("Empate")
	}

}

func race(carro *car) {
	fmt.Println(carro.name, "starts racing...")
	tempo := 0
	for i := 1; i < 3; i++ {
		tempo = rand.Intn(10-1) + 1
		carro.tempo = carro.tempo + tempo
		fmt.Println(carro.name, " - Volta", i, " - Tempo:", tempo)
		time.Sleep(time.Duration(tempo) * time.Second)
	}
}
