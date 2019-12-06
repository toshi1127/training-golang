package main // main
import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Conn struct {
	rootDir    string
	workdir    string
	reqUser    string
	user       string
	currentDir string
	granted    bool
	dataType   string
	ctrlConn   *net.TCPConn
	dataConn   *net.TCPConn
}

func (conn *Conn) writeMessage(code int, message string, v ...interface{}) error {
	msg := fmt.Sprintf(message, v...)
	_, err := fmt.Fprintf(conn.ctrlConn, "%d %s\r\n", code, msg)
	return err
}

func (conn *Conn) handleCommand(line string) {
	tokens := strings.Fields(line)
	opc := strings.ToUpper(tokens[0])
	opr := tokens[1:]
	switch opc {
	// ACCESS CONTROL COMMANDS
	case "USER":
		conn.handleUserCommand(opc, opr)
		return
	case "PASS":
		conn.handlePassCommand(opc, opr)
		return
	case "CWD":
		conn.handleCwdCommand(opc, opr)
		return
	case "PWD":
		conn.handlePwdCommand(opc, opr)
		return
	case "QUIT":
		conn.handleQuitCommand(opc, opr)
		return
	// TRANSFER PARAMETER COMMANDS
	case "PORT":
		conn.handlePortCommand(opc, opr)
		return
	case "TYPE":
		conn.handleTypeCommand(opc, opr)
		return
	case "STRU":
		conn.handleStruCommand(opc, opr)
		return
	case "MODE":
		conn.handleModeCommand(opc, opr)
		return
	// FTP SERVICE COMMANDS
	case "RETR":
		conn.handleRetrCommand(opc, opr)
		return
	default:
		conn.writeMessage(500, "%s not understood", opc)
		return
	}
}

func (conn *Conn) handleUserCommand(opc string, opr []string) {
	conn.reqUser = opr[0]
	conn.writeMessage(331, "User name ok, password required")
}

func (conn *Conn) handlePassCommand(opc string, opr []string) {
	ok, err := conn.checkPasswd(conn.user, opr[0])
	if err != nil {
		conn.writeMessage(550, "Checking password error")
		return
	}
	if ok {
		conn.user = conn.reqUser
		conn.reqUser = ""
		conn.writeMessage(230, "Password ok, continue")
	} else {
		conn.writeMessage(530, "Incorrect password, not logged in")
	}
}

func (conn *Conn) handleCwdCommand(opc string, opr []string) {
	path := conn.buildPath(opr[0])
	if f, err := os.Stat(conn.rootDir + path); err != nil || !f.IsDir() {
		conn.writeMessage(550, "Failed to change directory.")
		return
	}
	conn.currentDir = path
	conn.writeMessage(250, "Directory changed to "+path)
}

func (conn *Conn) handlePwdCommand(opc string, opr []string) {
	message := fmt.Sprintf("\"%s\" is the current directory", conn.currentDir)
	conn.writeMessage(257, message)
}

func (conn *Conn) handleRetlCommand(opc string, opr []string) {
}

func (conn *Conn) handleQuitCommand(opc string, opr []string) {
	conn.writeMessage(221, "Goodbye")
	conn.ctrlConn.Close()
}

func (conn *Conn) handlePortCommand(opc string, opr []string) {
	nums := strings.Split(opr[0], ",")
	portOne, _ := strconv.Atoi(nums[4])
	portTwo, _ := strconv.Atoi(nums[5])
	port := (portOne * 256) + portTwo
	host := nums[0] + "." + nums[1] + "." + nums[2] + "." + nums[3]
	dataConn, err := newSocket(host, port)
	if err != nil {
		conn.writeMessage(425, "Data connection failed")
		return
	}
	conn.dataConn = dataConn
	conn.writeMessage(200, "Connection established ("+strconv.Itoa(port)+")")
}

func (conn *Conn) handleTypeCommand(opc string, opr []string) {
	if strings.ToUpper(opr[0]) == "A" {
		conn.writeMessage(200, "Type set to ASCII")
	} else if strings.ToUpper(opr[0]) == "I" {
		conn.writeMessage(200, "Type set to binary")
	} else {
		conn.writeMessage(500, "Invalid type")
	}
}

func (conn *Conn) handleStruCommand(opc string, opr []string) {
	if strings.ToUpper(opr[0]) == "F" {
		conn.writeMessage(200, "OK")
	} else {
		conn.writeMessage(504, "STRU is an obsolete command")
	}
}

func (conn *Conn) handleModeCommand(opc string, opr []string) {
	if strings.ToUpper(opr[0]) == "S" {
		conn.writeMessage(200, "OK")
	} else {
		conn.writeMessage(504, "MODE is an obsolete command")
	}
}

func (conn *Conn) handleRetrCommand(opc string, opr []string) {
	fmt.Println("test")
	path := conn.buildPath(opr[0])
	f, err := os.Open(conn.rootDir + path)
	if err != nil {
		conn.writeMessage(550, "Failed to open file.")
		return
	}
	defer f.Close()
	conn.writeMessage(150, "Data transfer starting")
	fmt.Println("test2")
	io.Copy(conn.dataConn, f)
	fmt.Println("test3")
	conn.writeMessage(226, "Transfer complete.")
	return
}

func (conn *Conn) checkPasswd(user string, pass string) (bool, error) {
	return true, nil
}

func (conn *Conn) buildPath(filename string) string {
	fullPath := ""
	if len(filename) > 0 && filename[0:1] == "/" {
		fullPath = filepath.Clean(filename)
	} else if len(filename) > 0 && filename != "-a" {
		fullPath = filepath.Clean(conn.currentDir + "/" + filename)
	} else {
		fullPath = filepath.Clean(conn.currentDir)
	}
	fullPath = strings.Replace(fullPath, "//", "/", -1)
	return fullPath
}

func newSocket(host string, port int) (*net.TCPConn, error) {
	connectTo := buildTCPString(host, port)
	raddr, err := net.ResolveTCPAddr("tcp", connectTo)
	if err != nil {
		return nil, err
	}
	tcpConn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		return nil, err
	}
	return tcpConn, nil
}

func buildTCPString(hostname string, port int) (result string) {
	if strings.Contains(hostname, ":") {
		// ipv6
		if port == 0 {
			result = "[" + hostname + "]"
		} else {
			result = "[" + hostname + "]:" + strconv.Itoa(port)
		}
	} else {
		// ipv4
		if port == 0 {
			result = hostname
		} else {
			result = hostname + ":" + strconv.Itoa(port)
		}
	}
	return
}

func startSession(conn *net.TCPConn) {
	var c Conn
	c.ctrlConn = conn
	if rootDir, err := os.Getwd(); err == nil {
		fmt.Println(rootDir)
		c.rootDir = rootDir
	} else {
		log.Fatal(err)
	}
	c.writeMessage(220, "FTP Server Ready")
	go func() {
		defer func() {
			conn.Close()
			log.Println("close controll connection")
		}()
		input := bufio.NewScanner(conn)
		for input.Scan() {
			c.handleCommand(input.Text())
		}
	}()
}

func main() {
	port := flag.Int("port", 21, "port to bind")
	flag.Parse()
	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatal("err")
	}
	ctrlListener, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		log.Fatal("err")
	}
	for {
		conn, err := ctrlListener.AcceptTCP()
		if err != nil {
			log.Fatal("err")
			continue
		}
		startSession(conn)
	}
}