package internal

import (
	"fmt"
	"net"

	testerutils "github.com/codecrafters-io/tester-utils"
	testerutils_random "github.com/codecrafters-io/tester-utils/random"
)

func testReplGetaAckNonZero(stageHarness *testerutils.StageHarness) error {
	deleteRDBfile()
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
	}

	logger := stageHarness.Logger

	logger.Infof("Master is running on port 6379")

	replica := NewRedisBinary(stageHarness)
	replica.args = []string{
		"--port", "6380",
		"--replicaof", "localhost", "6379",
	}

	if err := replica.Run(); err != nil {
		return err
	}

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting: ", err.Error())
		return err
	}

	master := NewFakeRedisMaster(conn, logger)

	err = master.Handshake()
	if err != nil {
		return err
	}

	offset := 0
	err = master.GetAck(offset) // 37
	if err != nil {
		return err
	}
	offset += GetByteOffset([]string{"REPLCONF", "GETACK", "*"})

	cmd := []string{"PING"}
	master.Send(cmd) // 14
	offset += GetByteOffset(cmd)

	err = master.GetAck(offset) // 37
	if err != nil {
		return err
	}
	offset += GetByteOffset([]string{"REPLCONF", "GETACK", "*"})

	key := RandomAlphanumericString(testerutils_random.RandomInt(5, 20))
	value := RandomAlphanumericString(testerutils_random.RandomInt(5, 20))
	cmd = []string{"SET", key, value}
	master.Send(cmd) // 31
	offset += GetByteOffset(cmd)

	key = RandomAlphanumericString(testerutils_random.RandomInt(5, 20))
	value = RandomAlphanumericString(testerutils_random.RandomInt(5, 20))
	cmd = []string{"SET", key, value}
	master.Send(cmd) // 31
	offset += GetByteOffset(cmd)

	err = master.GetAck(offset)
	if err != nil {
		return err
	}

	listener.Close()
	return nil
}
