package telnet

import (
	"bufio"
	"errors"
	"net"
	"strings"
	"time"
	"fmt"
)

const (
	TIME_DELAY_AFTER_WRITE = 200
)

type Client struct {
	Address string
	Conn    net.Conn
	buf     [4096]byte
}

func (c *Client) Write(conn net.Conn, bufs []byte) (n int, err error) {
	n, err = conn.Write(bufs)
	if err != nil {
		return n, err
	}
	time.Sleep(time.Millisecond * TIME_DELAY_AFTER_WRITE)
	return n, err
}

func (c *Client) Connect(address string) (err error) {
	c.Conn, err = net.DialTimeout("tcp", address, 1 * time.Second)
	if err != nil {
		return err
	}
	c.Conn.SetDeadline(time.Now().Add(3 * time.Second))

	n, err := c.Conn.Read(c.buf[0:])
	if err != nil {
		return err
	}

	c.buf[1] = 252
	c.buf[4] = 252
	c.buf[7] = 252
	c.buf[10] = 252

	//n,err = Write(conn,buf[0:n])
	if err != nil {
		return err
	}

	n, err = c.Conn.Read(c.buf[0:])
	if err != nil {
		return err
	}

	c.buf[1] = 252
	c.buf[4] = 251
	c.buf[7] = 252
	c.buf[10] = 254
	c.buf[13] = 252
	n, err = c.Write(c.Conn, c.buf[0:n])
	if err != nil {
		return err
	}

	n, err = c.Conn.Read(c.buf[0:])
	if err != nil {
		return err
	}

	c.buf[1] = 252
	c.buf[4] = 252
	n, err = c.Write(c.Conn, c.buf[0:n])
	if err != nil {
		return err
	}

	n, err = c.Conn.Read(c.buf[0:])
	if err != nil {
		return err
	}

	c.Conn.SetDeadline(time.Now().Add(15 * time.Second))

	return err
}

func (c *Client) Login(username string, password string) error {
login:	
	n, err := c.Write(c.Conn, []byte("\n"))
	if err != nil {
		return err
	}
	n, err = c.Write(c.Conn, []byte("\n"))
	if err != nil {
		return err
	}
	n, err = c.Write(c.Conn, []byte("\n"))
	if err != nil {
		return err
	}
	n, err = c.Write(c.Conn, []byte(username+"\n"))
	if err != nil {
		return err
	}

	n, err = c.Conn.Read(c.buf[0:])
	if err != nil {
		return err
	}

	n, err = c.Write(c.Conn, []byte(password+"\n"))
	if err != nil {
		return err
	}

	testCase := 0

	for {
		if testCase > 2 {
			return errors.New("error login")
		}
		testCase++
		n, err = c.Write(c.Conn, []byte("\n"))
		if err != nil {
			return err
		}
		n, err = c.Write(c.Conn, []byte("\n"))
		if err != nil {
			return err
		}
		n, err = c.Write(c.Conn, []byte("\n"))
		if err != nil {
			return err
		}
		n, err = c.Conn.Read(c.buf[0:])
		if err != nil {
			return err
		}

		fmt.Println(string(c.buf[0:n]))
		if strings.HasSuffix(string(c.buf[0:n]), ">") {
			break
		}
		if strings.HasSuffix(string(c.buf[0:n]), "Username(1-64 chars):") {
			goto login
			break
		}
		if strings.HasSuffix(string(c.buf[0:n]), "Password(1-96 chars):") {
			n, err = c.Write(c.Conn, []byte(password+"\n"))
			if err != nil {
				return err
			}
		}
	}
	n, err = c.Write(c.Conn, []byte("enable\n"))
	if err != nil {
		return err
	}

	n, err = c.Conn.Read(c.buf[0:])
	if err != nil {
		return err
	}
	fmt.Println(string(c.buf[0:n]))

	n, err = c.Write(c.Conn, []byte(password+"\n"))
	if err != nil {
		return err
	}

	n, err = c.Conn.Read(c.buf[0:])
	if err != nil {
		return err
	}
	fmt.Println(string(c.buf[0:n]))

	n, err = c.Write(c.Conn, []byte("terminal length 0\n"))
	if err != nil {
		return err
	}

	n, err = c.Conn.Read(c.buf[0:])
	if err != nil {
		return err
	}
	fmt.Println(string(c.buf[0:n]))
	return err
}

func (c *Client) Cmd(shell string) (context string, err error) {
	_, err = c.Write(c.Conn, []byte(shell+"\n\n"))
	if err != nil {
		return "", err
	}
	
	//
	//for {
	//	n, err = conn.Read(buf[0:])
	//	if err != nil {
	//		return "",err
	//	}
	//	context += string(buf[0:n])
	//	if strings.HasSuffix(string(buf[0:n]), "#") {
	//		break
	//	}
	//}

	reader := bufio.NewReader(c.Conn)

	if reader == nil {
		return "", errors.New("Create reader failed.")
	}

	for {
		n, err := reader.Read(c.buf[0:])
		if err != nil {
			return "", err
		}
		context += string(c.buf[0:n])
		if strings.HasSuffix(string(c.buf[0:n]), "#") {
			break
		}
	}

	return context, err
}
