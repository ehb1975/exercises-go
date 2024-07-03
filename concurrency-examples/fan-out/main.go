package main

import (
 "fmt"
 "sync"
)

func generator(nums ...int) <-chan int {
 myChannel := make(chan int) //declare a channel

 go func() {
  //iterate the nums data and sends it to channel
  for _, val := range nums {
   myChannel <- val
  }
  close(myChannel)
 }()

 return myChannel
}

func main() {
 data1 := []int{1, 2, 3, 4, 5}
 data2 := []int{10, 20, 30, 40, 50}
 var wg sync.WaitGroup

 //it receives a "receive-only" directional channel
 ch1 := generator(data1...)
 ch2 := generator(data2...)
 wg.Add(2)

 //we will loop through both the channels till all data is sent and marked as close
 go func() {
  for val := range ch1 {
   fmt.Printf("Channel1 data: %v\n", val)
  }
  wg.Done()
 }()

 go func() {
  for val := range ch2 {
   fmt.Printf("Channel2 data: %v\n", val)
  }
  wg.Done()
 }()

 wg.Wait() //will wait till the above goroutines are marked as done
}