package kakfaconsumer

import (
	"bufio"
	"chatbot/pkg/utils"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

func ConsumedLoanFromKafka() {
	group_id := "prod_kafka_app"

	// do not change this value in production
	// topic := "prod.loan_response"
	// topic := "test.test_topic"

	topic := utils.GetEnv("topic")

	//print to check the topic
	fmt.Println("Consumer Topic: ", topic)

	// Read the client configuration from client.properties
	conf := ReadConfig()

	// Create a new kafka.Reader with the given configuration
	reader := CreateConsumer(conf["sasl.username"].(string), conf["sasl.password"].(string), conf["bootstrap.servers"].(string), group_id, topic)

	log.Println("Starting to read messages...")
	// Read messages from the topic
readMessageLoop:
	for {
		// Read a message from the topic
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			// log the error and continue reading the next message
			log.Printf("Reader encountered an error: %s", err)
			continue readMessageLoop
		}

		// fmt.Println(m.Value)
		var kafkaMsg map[string]any
		if err := json.Unmarshal(m.Value, &kafkaMsg); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		// status := sharedfunctions.GetBoolFromMap(kafkaMsg, "success")

		fmt.Println("Received Loan Status: ", kafkaMsg)

		result, err := StoreFailedDisbursement(kafkaMsg)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		fmt.Println("Saving Transaction Result: ", result)
		//print to check the topic
		fmt.Println("Consumer Topic: ", topic)
	}
}

func ReadConfig() map[string]any {
	// reads the client configuration from client.properties
	// and returns it as a key-value map
	m := make(map[string]any)

	file, err := os.Open("client.properties")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %s", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "#") && len(line) != 0 {
			kv := strings.Split(line, "=")
			parameter := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			m[parameter] = value
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Failed to read file: %s", err)
		os.Exit(1)
	}

	return m
}

func CreateConsumer(user, pass, host, group_id, topic string) *kafka.Reader {
	// Create a new plain SASL mechanism with the given username and password
	mechanism := plain.Mechanism{Username: user, Password: pass}

	// Create a new dialer with the SASL mechanism and other configuration
	dialer := &kafka.Dialer{
		// Set the timeout to 30 seconds
		Timeout: 30 * time.Second,
		// Enable dual stack
		DualStack: true,
		// Set the TLS configuration
		TLS: &tls.Config{
			// Set the minimum TLS version to 1.2
			MinVersion: tls.VersionTLS12,
		},
		// Set the SASL mechanism
		SASLMechanism: mechanism,
	}

	// Create a new kafka.Reader with the given configuration
	return kafka.NewReader(kafka.ReaderConfig{
		// Set the brokers to the given host
		Brokers: []string{host},
		// Set the group ID to the given value
		GroupID: group_id,
		// Set the topic to the given value
		Topic: topic,
		// Set the dialer to the created dialer
		Dialer: dialer,
	})
}
