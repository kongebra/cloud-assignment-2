package main

import (
	"strconv"
	"testing"
	"time"
)

func Test_ApiInit(t *testing.T) {
	info := "Some information"
	version := "v1"

	var api = API{
		Info: info,
		Version: version,
	}

	if api.Version != version {
		t.Error("Version is incorrect, should be: " + version + ", is: " + api.Version)
	}

	if api.Info != info {
		t.Error("Info is incorrect, should be: " + info + ", is: " + api.Info)
	}
}

func TestAPI_CalculateUptime(t *testing.T) {
	info := "Some information"
	version := "v1"

	var api = API{
		Info: info,
		Version: version,
	}

	t0 := time.Now().Unix()

	time.Sleep(3 * time.Second)

	t1 := time.Now().Unix()

	elapsed := int(t1 - t0)

	api.CalculateUptime(elapsed)

	str := "PT" + strconv.Itoa(elapsed) + "S"

	if api.Uptime != str {
		t.Error("Uptime is incorrect, should be: " + str + ", is: " + api.Uptime)
	}


}
