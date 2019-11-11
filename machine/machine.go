package machine

import (
	"github.com/denisbrodbeck/machineid"
)

func GetMachineID() (string, error) {
	return machineid.ID()
}
