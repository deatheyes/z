package main


import (
	"fmt"
	"strings"
	"path"
	"github.com/urfave/cli"
	"github.com/samuel/go-zookeeper/zk"
)

func LsAction(c *cli.Context) {
	zPath := Meta.Path
	if c.NArg() > 0 {
		zPath = c.Args().Get(0)
	}

	children, _, err := Meta.Conn.Children(zPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := strings.Join(children, "\n")
	fmt.Println(data)
}

func GetAction(c *cli.Context) {
	zPath := Meta.Path
	if c.NArg() > 0 {
		zPath = c.Args().Get(0)
	}

	content, _, err := Meta.Conn.Get(zPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(content[:]))
}

func SetAction(c *cli.Context) {
	zPath := Meta.Path
	val := ""
	if c.NArg() == 1 {
		val = c.Args().Get(0)
	} else if c.NArg() > 1 {
		zPath = c.Args().Get(0)
		val = c.Args().Get(1)
	}

	if _, err := Meta.Conn.Set(zPath, []byte(val), 0); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("OK")
}

func CreateAction(c *cli.Context) {
	if c.NArg() < 1 {
		fmt.Println("no path specified")
		return
	}

	zPath := path.Clean(c.Args().Get(0))
	flag, _, err := Meta.Conn.Exists(zPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	if flag {
		fmt.Println("path: " + zPath + " exists")
		return
	}

	val := []byte{}
	acl := []zk.ACL{}
	if c.NArg() > 1 {
		val = []byte(c.Args().Get(1))
	}

	if _, err := Meta.Conn.Create(zPath, val, 0 ,acl); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("OK")
}


func DelAction(c *cli.Context) {
	if c.NArg() < 1 {
		fmt.Println("no path specified")
		return
	}

	zPath := path.Clean(c.Args().Get(0))
	if err := Meta.Conn.Delete(zPath, 0); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("OK")
}

func CdAction(c *cli.Context) {
	if c.NArg() < 1 {
		return
	}

	zPath := path.Join(Meta.Path, c.Args().Get(0))
	zPath = path.Clean(zPath)
	flag, _, err := Meta.Conn.Exists(zPath)
	if err != nil {
		fmt.Println(err)
	}
	if flag {
		Meta.Path = zPath
	} else {
		fmt.Println("path: " + zPath + " not exists")
	}
}
