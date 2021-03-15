// Copyright (c) 2021 Qianyun, Inc. All rights reserved.
package main

import (
	log "github.com/sirupsen/logrus"
)

type RealServerMessageHandler struct {
	LpConnHandler *ConnHandler
	ConnPool      *ConnHandlerPool
	UserId        string
	ClientKey     string
}

func (messageHandler *RealServerMessageHandler) Encode(msg interface{}) []byte {
	if msg == nil {
		return []byte{}
	}

	return msg.([]byte)
}

func (messageHandler *RealServerMessageHandler) Decode(buf []byte) (interface{}, int) {
	return buf, len(buf)
}

func (messageHandler *RealServerMessageHandler) MessageReceived(connHandler *ConnHandler, msg interface{}) {
	if connHandler.NextConn != nil {
		data := msg.([]byte)
		message := Message{Type: P_TYPE_TRANSFER}
		message.Data = data
		log.Infof("Transfer [%s] to [%s].", connHandler.NextConn.conn.LocalAddr(), connHandler.NextConn.conn.RemoteAddr())
		connHandler.NextConn.Write(message)
	}
}

func (messageHandler *RealServerMessageHandler) ConnSuccess(connHandler *ConnHandler) {
	log.Info("get proxy connection:", messageHandler.UserId)
	proxyConnHandler, err := messageHandler.ConnPool.Get() // 192.168.88.141:4993
	if err != nil {
		log.Error("get proxy connection err:", err, "uri:", messageHandler.UserId)
		message := Message{Type: TYPE_DISCONNECT}
		message.Uri = messageHandler.UserId
		messageHandler.LpConnHandler.Write(message)
		connHandler.conn.Close()
	} else {
		proxyConnHandler.NextConn = connHandler // 127.0.0.1:5000
		connHandler.NextConn = proxyConnHandler
		message := Message{Type: TYPE_CONNECT}
		message.Uri = messageHandler.UserId + "@" + messageHandler.ClientKey // 一个Client对应一个ClientKey,一个tcp连接对应一个UserId
		proxyConnHandler.Write(message)
		log.Info("realserver connect success, notify proxyserver:", message.Uri)
	}
}

func (messageHandler *RealServerMessageHandler) ConnError(connHandler *ConnHandler) {
	conn := connHandler.NextConn
	if conn != nil {
		message := Message{Type: TYPE_DISCONNECT}
		message.Uri = messageHandler.UserId
		conn.Write(message)
		conn.NextConn = nil
	}

	connHandler.messageHandler = nil
}

func (messageHandler *RealServerMessageHandler) ConnFailed() {
	message := Message{Type: TYPE_DISCONNECT}
	message.Uri = messageHandler.UserId
	messageHandler.LpConnHandler.Write(message)
}
