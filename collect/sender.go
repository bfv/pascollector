package collect

import (
	"fmt"
	"time"

	"github.com/bfv/pascollector/types"
)

func SendData(config types.ConfigFile) {
	fmt.Println("send: " + time.Now().String())
}
