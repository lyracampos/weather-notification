package worker

import (
	"fmt"
	"time"

	"weather-notification/configs"
)

func Run(config *configs.Config) {
	keepRunning := true
	count := 1
	for keepRunning {
		time.Sleep(2 * time.Second)
		fmt.Println("...")
		fmt.Println("worker is running...%i:", count)
		count++
	}
}
