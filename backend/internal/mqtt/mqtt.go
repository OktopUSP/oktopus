/*
	Runs MQTT broker trough a Docker container.
	Better approach would be to use docker api to Go language, but os/exec lib is already enough for our purpose,
	since it's more convenient and easier to use docker shell commands, and it's already a start point.
*/
package mqtt

import (
	"log"
	"os/exec"
)

func StartMqttBroker() {

	//TODO: Start Container through Docker SDK for GO, eliminating docker-compose and shell comands.

	cmd := exec.Command("sudo", "docker", "compose", "-f", "internal/mqtt/docker-compose.yml", "up", "-d")

	err := cmd.Run()

	if err != nil {
		log.Fatal(err.Error())
		return
	}
}
