package main

import (
	"bufio"
	"container/ring"
	"fmt"
	"net"
)

const maxRefNum = 100
const etx = 3

type Client struct {
	addr       string
	user       string
	password   string
	accessCode string

	ringCounter *ring.Ring

	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func (c *Client) initRefNum() {
	ringCounter := ring.New(maxRefNum)

	for i := 0; i < maxRefNum; i++ {
		ringCounter.Value = []byte(fmt.Sprintf("%02d", i))
		ringCounter = ringCounter.Next()
	}

	c.ringCounter = ringCounter
}

func (c *Client) nextRefNum() []byte {
	refNum := (c.ringCounter.Value).([]byte)

	c.ringCounter = c.ringCounter.Next()

	return refNum
}

func (c *Client) Connect() error {
	//init ring counter
	c.initRefNum()

	conn, _ := net.Dial("tcp", c.addr)
	c.conn = conn

	//create buffered reader and writer
	c.reader = bufio.NewReader(conn)
	c.writer = bufio.NewWriter(conn)

	c.writer.Write(createLoginReq(c.nextRefNum(), c.user, c.password))

	//fmt.Println("flushing that bitch")
	c.writer.Flush()
	resp, _ := c.reader.ReadString(etx)

	err := parseSessionResp(resp)

	fmt.Println(resp)

	return err
}

// func parseSessionResp(response string) error {
// 	/*  */
// }
