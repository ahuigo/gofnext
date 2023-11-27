package decorator

import (
	"math/rand"
	"time"
)

func sleepRandom(min time.Duration, max time.Duration) {
    randomDuration := time.Duration(rand.Int63n(int64(max-min)) + int64(min))

    time.Sleep(randomDuration)

}