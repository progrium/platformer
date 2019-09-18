package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
)

type ReplayServer struct {
	buffer  []ReplayFrame
	replays [][]ReplayFrame
}

func ListenAndServe() {
	rs := &ReplayServer{}
	pc, err := net.ListenPacket("udp", ":3333")
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	for {
		buffer := make([]byte, 1024)
		pc.ReadFrom(buffer)
		buf := bytes.NewBuffer(buffer)
		dec := gob.NewDecoder(buf)
		var f ReplayFrame
		err := dec.Decode(&f)
		if err != nil {
			log.Fatal(err)
		}
		if f.Tick == 1 {
			if len(rs.buffer) > 1 {
				rs.replays = append(rs.replays, rs.buffer)
			}
			rs.buffer = []ReplayFrame{}
			// TODO: start replaying
		}
		rs.buffer = append(rs.buffer, f)
	}

}
