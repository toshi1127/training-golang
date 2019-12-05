package main

// ftp-responsecode
// http://www.atmarkit.co.jp/fnetwork/rensai/netpro10/ftp-responsecode.html

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

var remoteRoot string

type ftpConn struct {
	conn         net.Conn
	prevCmd      string
	dir          string
	sep          string
	username     string
	password     string
	addr         string
	pasvListener net.Listener
}

func newFTPConn(c net.Conn) *ftpConn {
	separator := "/"
	if runtime.GOOS == "windows" {
		separator = "\\"
	}
	return &ftpConn{conn: c, dir: remoteRoot, sep: separator}
}

func (c *ftpConn) println(s ...interface{}) {
	s = append(s, "\r\n")
	_, err := fmt.Fprint(c.conn, s...)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *ftpConn) user(cmds []string) {
	if len(cmds) < 2 {
		c.println("501 Syntax error in parameters or arguments.")
		return
	}
	c.username = cmds[1]
	c.println("331 User name okay, need password.")
}

func (c *ftpConn) pass(cmds []string) {
	if len(cmds) < 2 {
		c.println("501 Syntax error in parameters or arguments.")
		return
	}
	c.password = cmds[1]
	c.println("230 logged in, proceed.")
}

func (c *ftpConn) port(cmds []string) {
	if len(cmds) < 2 {
		c.println("501 Syntax error in parameters or arguments.")
		return
	}
	ips := strings.Split(cmds[1], ",")
	if len(ips) != 6 {
		c.println("501 Syntax error in parameters or arguments.")
		return
	}
	p1, err := strconv.Atoi(ips[4])
	if err != nil {
		c.println("501 Syntax error in parameters or arguments.")
		return
	}
	p2, err := strconv.Atoi(ips[5])
	if err != nil {
		c.println("501 Syntax error in parameters or arguments.")
		return
	}
	c.addr = fmt.Sprintf("%s.%s.%s.%s:%d", ips[0], ips[1], ips[2], ips[3], p1*256+p2)
	c.println("200 PORT Command okay.")
}

func (c *ftpConn) pasv() {
	var err error
	c.pasvListener, err = net.Listen("tcp4", "")
	if err != nil {
		c.println("451 aborted. Local error in processing.")
		return
	}
	_, port, err := net.SplitHostPort(c.pasvListener.Addr().String())
	if err != nil {
		c.println("451 aborted. Local error in processing.")
		c.pasvListener.Close()
		return
	}
	host, _, err := net.SplitHostPort(c.conn.LocalAddr().String())
	if err != nil {
		c.println("451 aborted. Local error in processing.")
		c.pasvListener.Close()
		return
	}
	ipAdder, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		c.println("451 aborted. Local error in processing.")
		c.pasvListener.Close()
		return
	}
	ips := ipAdder.IP.To4()
	pVal, err := strconv.Atoi(port)
	if err != nil {
		c.println("451 aborted. Local error in processing.")
		c.pasvListener.Close()
		return
	}
	adder := fmt.Sprintf("%d,%d,%d,%d", ips[0], ips[1], ips[2], ips[3]) + fmt.Sprintf(",%d,%d", pVal/256, pVal%256)
	c.println(fmt.Sprintf("227 Entering Passive Mode (%s).", adder))
}

func (c *ftpConn) syst() {
	var osType string
	switch runtime.GOOS {
	case "windows":
		osType = "Windows NT"
	default:
		osType = "UNIX"
	}
	c.println(fmt.Sprintf("215 %s", osType))
}

func (c *ftpConn) pwd() {
	relPath, err := filepath.Rel(remoteRoot, c.dir)
	if err != nil {
		log.Print(err)
		c.println("451 Requested action aborted. Local error in processing.")
		return
	}
	relPath = strings.Replace(relPath, "\\", "/", -1)
	if []rune(relPath)[0] == '.' {
		relPath = "/"
	} else {
		relPath = "/" + relPath
	}
	c.println(fmt.Sprintf("257 \"%s\" is current directory", relPath))
}

func (c *ftpConn) list(cmds []string) {
	var target string
	var err error
	switch len(cmds) {
	case 1:
		target = "."
	case 2:
		target = cmds[1]
	default:
		c.println("501 Syntax error in parameters or arguments.")
		return
	}
	conn, err := c.createDataConn()
	if err != nil {
		c.println("425 Can't open data connection.")
		return
	}
	defer conn.Close()
	c.println("150 Opening ASCII mode data connection")

	fullPath, err := c.resolveFullPath(target)
	if err != nil {
		c.println("451 Requested action aborted. Local error in processing.")
		return
	}
	stat, err := os.Stat(fullPath)
	if err != nil || !stat.IsDir() {
		c.println("451 Requested action aborted. Local error in processing.")
		return
	}
	if err != nil {
		fmt.Println(err)
		c.println("450 Requested file action not taken.")
		return
	}
	if stat.IsDir() {
		var files []os.FileInfo
		files, err = ioutil.ReadDir(fullPath)
		if err != nil {
			fmt.Println(err)
			c.println("550 Requested action not taken.")
			return
		}
		for _, f := range files {
			if strings.ToUpper(cmds[0]) == "LIST" {
				_, err = fmt.Fprintf(conn, "%s %d\t%s\t%s\r\n", f.Mode(), f.Size(), f.ModTime(), f.Name())
			} else {
				_, err = fmt.Fprintf(conn, "%s\r\n", f.Name())
			}
			if err != nil {
				c.println("426 Connection closed; transfer aborted.")
				return
			}
		}
	} else {
		if strings.ToUpper(cmds[0]) == "LIST" {
			_, err = fmt.Fprintf(conn, "%s %d\t%s\t%s\r\n", stat.Mode(), stat.Size(), stat.ModTime(), stat.Name())
		} else {
			_, err = fmt.Fprintf(conn, "%s\r\n", stat.Name())
		}
		if err != nil {
			c.println("426 Connection closed; transfer aborted.")
			return
		}
	}
	c.println("226 Closing data connection.")
}

func (c *ftpConn) createDataConn() (conn io.ReadWriteCloser, err error) {
	switch c.prevCmd {
	case "PORT":
		conn, err = net.Dial("tcp", c.addr)
	case "PASV":
		conn, err = c.pasvListener.Accept()
	default:
		return nil, fmt.Errorf("previuos command not Connection: %s", c.prevCmd)
	}
	return
}

func (c *ftpConn) cwd(cmds []string) {
	if len(cmds) != 2 {
		c.println("501 Syntax error, command unrecognized.")
		return
	}
	fullPath, err := c.resolveFullPath(cmds[1])
	if err != nil {
		c.println("451 Requested action aborted. Local error in processing.")
		return
	}
	stat, err := os.Stat(fullPath)
	if err != nil || !stat.IsDir() {
		c.println("451 Requested action aborted. Local error in processing.")
		return
	}
	c.dir = fullPath
	fmt.Println(c.dir)
	c.println("250 Requested file action okay, completed.")
}

func (c *ftpConn) resolveFullPath(p string) (fullPath string, err error) {
	if []rune(p)[0] == '/' {
		root := remoteRoot + string(filepath.Separator)
		fullPath = strings.Replace(p, "/", root, 1)
	} else {
		fullPath = filepath.Join(c.dir, p)
	}
	rel, err := filepath.Rel(remoteRoot, fullPath)
	if err != nil {
		return
	}
	if strings.Split(rel, string(filepath.Separator))[0] == ".." {
		err = fmt.Errorf("%s upper remote root", p)
		return
	}
	return
}

func (c *ftpConn) size(cmds []string) {
	if len(cmds) != 2 {
		c.println("501 Syntax error in parameters or arguments.")
		return
	}

	fullPath, err := c.resolveFullPath(cmds[1])
	if err != nil {
		c.println("451 Requested action aborted. Local error in processing.")
		return
	}
	stat, err := os.Stat(fullPath)
	if err != nil {
		log.Print(err)
		c.println("451 Requested action aborted. Local error in processing.")
		return
	}
	if stat.IsDir() {
		c.println("450 Requested file action not taken.")
		return
	}
	c.println(fmt.Sprintf("213 %d", stat.Size()))
}

func (c *ftpConn) retr(cmds []string) {
	if len(cmds) != 2 {
		c.println("501 Syntax error in parameters or arguments.")
		return
	}
	fullPath, err := c.resolveFullPath(cmds[1])
	if err != nil {
		c.println("451 Requested action aborted. Local error in processing.")
		return
	}
	file, err := os.Open(fullPath)
	if err != nil {
		c.println("450 Requested file action not taken.")
		return
	}
	stat, err := os.Stat(fullPath)
	if err != nil || stat.IsDir() {
		c.println("450 Requested file action not taken.")
		return
	}
	conn, err := c.createDataConn()
	if err != nil {
		c.println("425 Can't open data connection.")
		return
	}
	defer conn.Close()
	c.println(fmt.Sprintf("150 Opening BINARY mode data connection for '%s'(%d bytes).", stat.Name(), stat.Size()))
	_, err = io.Copy(conn, file)
	if err != nil {
		c.println("450 Requested file action not taken.")
		return
	}
	c.println("226 Closing data connection.")
}

func (c *ftpConn) mkd(cmds []string) {
	if len(cmds) != 2 {
		c.println("501 Syntax error in parameters or arguments.")
		return
	}
	fullPath, err := c.resolveFullPath(cmds[1])
	if err != nil {
		c.println("451 Requested action aborted. Local error in processing.")
		return
	}
	if err := os.Mkdir(fullPath, 0777); err != nil {
		c.println("450 	Requested file action not taken.")
		return
	}
	c.println(fmt.Sprintf("\"%s\" created.", cmds[1]))
}

func (c *ftpConn) stor(cmds []string) {
	if len(cmds) != 2 {
		c.println("501 Syntax error in parameters or arguments.")
		return
	}
	file, err := os.Create(cmds[1])
	if err != nil {
		c.println("550 Requested action not taken.")
		return
	}
	defer file.Close()
	c.println("150 File status okay; about to open data connection.")
	conn, err := c.createDataConn()
	if err != nil {
		c.println("425 Can't open data connection.")
		return
	}
	defer conn.Close()
	if _, err := io.Copy(file, conn); err != nil {
		c.println("450 	Requested file action not taken.")
		return
	}
	c.println("226 Closing data connection.")
}

func (c *ftpConn) handleConn() {
	defer c.conn.Close()
	s := bufio.NewScanner(c.conn)
	c.println("220 Ready.")
LOOP:
	for s.Scan() {
		log.Println(s.Text())
		cmds := strings.Fields(s.Text())
		if len(cmds) == 0 {
			continue
		}
		switch strings.ToUpper(cmds[0]) {
		case "USER":
			c.user(cmds)
		case "PASS":
			c.pass(cmds)
		case "PORT":
			c.port(cmds)
		case "SYST":
			c.syst()
		case "XPWD":
			c.pwd()
		case "PWD":
			c.pwd()
		case "LIST":
			c.list(cmds)
		case "NLST":
			c.list(cmds)
		case "PASV":
			c.pasv()
		case "CWD":
			c.cwd(cmds)
		case "SIZE":
			c.size(cmds)
		case "RETR":
			c.retr(cmds)
		case "MKD":
			c.mkd(cmds)
		case "XMKD":
			c.mkd(cmds)
		case "STOR":
			c.stor(cmds)
		case "QUIT":
			c.println("221 Goodbye.")
			break LOOP
		default:
			c.println("502 Command not implemented.")
		}
		c.prevCmd = strings.ToUpper(cmds[0])
	}
	log.Printf("Close: %s", c.conn.RemoteAddr().String())
}

func main() {
	var err error
	remoteRoot, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	port := flag.Int("p", 21, "FTP port")
	flag.Parse()
	listener, err := net.Listen("tcp4", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("run FTP Server")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go newFTPConn(conn).handleConn()
	}
}