package app

import "fmt"

func Uptime(uptime uint64) string {
	d := uint64(uptime/86400)
	h := uint64((uptime - (d * 86400))/3600)
	m := uint64((uptime - (d * 86400 + h * 3600))/60)

	return fmt.Sprintf("%v days, %v hours, %v mins", d, h, m)
}