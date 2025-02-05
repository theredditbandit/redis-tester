package internal

import (
	"fmt"

	testerutils "github.com/codecrafters-io/tester-utils"
)

func testReplMultipleReplicas(stageHarness *testerutils.StageHarness) error {
	deleteRDBfile()
	master := NewRedisBinary(stageHarness)
	master.args = []string{
		"--port", "6379",
	}

	if err := master.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	var replicas []*FakeRedisReplica

	for j := 0; j < 3; j++ {
		conn, err := NewRedisConn("", "localhost:6379")
		if err != nil {
			fmt.Println("Error connecting to TCP server:", err)
		}
		defer conn.Close()
		replica := NewFakeRedisReplica(conn, logger)
		replicas = append(replicas, replica)
		replica.LogPrefix = fmt.Sprintf("[replica-%v] ", j+1)
	}

	conn, err := NewRedisConn("", "localhost:6379")
	if err != nil {
		fmt.Println("Error connecting to TCP server:", err)
	}
	defer conn.Close()
	client := NewFakeRedisMaster(conn, logger)
	client.LogPrefix = "[client] "

	for i := 0; i < len(replicas); i++ {
		replica := replicas[i]
		err = replica.Handshake()
		if err != nil {
			return err
		}
	}

	kvMap := map[int][]string{
		1: {"foo", "123"},
		2: {"bar", "456"},
		3: {"baz", "789"},
	}
	for i := 1; i <= len(kvMap); i++ {
		// We need order of commands preserved
		key, value := kvMap[i][0], kvMap[i][1]
		client.Log(fmt.Sprintf("Setting key %s to %s", key, value))
		client.Send([]string{"SET", key, value})
	}

	for j := 0; j < 3; j++ {
		replica := replicas[j]
		logger.Infof("Testing Replica : %v", j+1)

		// Redis will send SELECT, but not expected from Users.
		err, _ = replica.readAndAssertMessagesWithSkip([]string{"SET", "foo", "123"}, "SELECT", true)
		if err != nil {
			return err
		}

		for i := 2; i <= len(kvMap); i++ {
			// We need order of commands preserved
			key, value := kvMap[i][0], kvMap[i][1]

			err, _ = replica.readAndAssertMessages([]string{"SET", key, value}, true)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
