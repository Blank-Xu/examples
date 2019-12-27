package ftp

import (
	"log"
	"net"
	"sync"
	"time"
)

var _aliveTCPConnMap sync.Map

func startAliveCheck(interval uint32) {
	if interval < 10 {
		interval = 10
	}

	aliveChan := make(chan bool, 1)
	aliveFunc := func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println("AliveCheck failed, panic:", err)
				aliveChan <- true
			}
		}()

		for {
			now := time.Now().Unix()
			_aliveTCPConnMap.Range(func(key, value interface{}) bool {
				if value.(int64) <= now-int64(interval) {
					conn, ok := key.(*net.TCPConn)
					if ok && conn != nil {
						conn.Close()
					}
				}
				return true
			})

			time.Sleep(time.Second * time.Duration(interval))
		}
	}

	go func() {
		for {
			if <-aliveChan {
				go aliveFunc()
			}
		}
	}()

	aliveChan <- true
}

func addAliveCheck(conn *net.TCPConn) {
	_aliveTCPConnMap.LoadOrStore(conn, time.Now().Unix())
}

func deleteAliveCheck(conn *net.TCPConn) {
	_aliveTCPConnMap.Delete(conn)
}