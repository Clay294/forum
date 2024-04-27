package common

import (
	"math/rand"
	"time"
)

var Rander = rand.New(rand.NewSource(time.Now().UnixNano()))
