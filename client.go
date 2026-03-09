package main

import (
	"encoding/json"
	"fmt"
	"net"
)

func daemonCmd(cmd Cmd) (*CmdResponse, error) {
	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		return nil, fmt.Errorf("daemon not running? %w", err)
	}
	defer conn.Close()

	if err := json.NewEncoder(conn).Encode(CmdRequest{Cmd: cmd}); err != nil {
		return nil, err
	}

	var resp CmdResponse
	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
