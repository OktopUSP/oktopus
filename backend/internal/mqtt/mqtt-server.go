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

// Get Mqtt Broker up and running
func StartMqttBroker() {

	//TODO: Start Container through Docker SDK for GO, eliminating docker-compose and shell comands.
	//TODO: Create Broker with user, password and CA certificate.
	//TODO: Set broker access control list to topics.
	//TODO: Set MQTTv5 CONNACK packet with topic for agent to use.

	cmd := exec.Command("sudo", "docker", "compose", "-f", "internal/mosquitto/docker-compose.yml", "up", "-d")

	err := cmd.Run()

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	log.Println("Broker Mqtt Up and Running!")
}
