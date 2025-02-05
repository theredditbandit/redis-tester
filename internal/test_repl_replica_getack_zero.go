package internal

import (
	"fmt"
	"net"

	testerutils "github.com/codecrafters-io/tester-utils"
)

func testReplGetaAckZero(stageHarness *testerutils.StageHarness) error {
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

	err = master.GetAck(0)
	if err != nil {
		return err
	}

	conn.Close()
	listener.Close()
	return nil
}
