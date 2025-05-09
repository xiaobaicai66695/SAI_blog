package service

import (
	"SAI_blog/repository"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

var groupConnsMux sync.Mutex

// groupId -> uid ->conns
var groupConns = make(map[int64]map[int64]*websocket.Conn)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandlerChatConnection(w http.ResponseWriter, r *http.Request, groupId int64, uid int64) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return fmt.Errorf("升级HTTP到websocket时发生错误")
	}
	err = repository.AddUserToGroup(uid, groupId, conn)
	if err != nil {
		return fmt.Errorf("添加用户连接状态时出现错误")
	}
	groupConnsMux.Lock()
	if groupConns[groupId] == nil {
		groupConns[groupId] = make(map[int64]*websocket.Conn)
	} else if oldConn, ok := groupConns[groupId][uid]; ok {
		oldConn.Close()
	}
	groupConns[groupId][uid] = conn
	groupConnsMux.Unlock()
	go listenAndBroadcastMessage(conn, groupId, uid)
	go heartbeat(conn)
	return nil
}

func listenAndBroadcastMessage(conn *websocket.Conn, groupId int64, uid int64) {
	for {
		_, msgStr, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("连接断开(uid = %d ,groupId = %d) ; %v\n", uid, groupId, err)
			return
		}
		msg := packingMsg(string(msgStr), groupId, uid)
		_ = repository.SaveMsgToRedis(msg, groupId)
		err = broadcastMessage(msg, groupId)
		if err != nil {
			fmt.Println("广播消息时失败")
		}
	}
}

func broadcastMessage(msg *repository.Group, groupId int64) error {
	groupConnsMux.Lock()
	defer groupConnsMux.Unlock()
	uids, err := repository.GetUserIdsByGroupId(groupId)
	if err != nil {
		return fmt.Errorf("获取群组成员失败；%v", err)
	}
	//uids = append(uids, 6, 7)
	for _, uid := range uids {
		if userConns, ok := groupConns[groupId]; ok {
			if conn, ok := userConns[uid]; ok {
				msgBytes, _ := json.Marshal(msg)
				if err = conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
					fmt.Printf("用户(uid = %d)发送失败，准备删除连接\n", uid)
					conn.Close()
					//delete(userConns, uid)
					removeConnection(groupId, uid)
				}
			} else {
				//写入未读消息列表
			}
		}
	}
	return nil
}

func heartbeat(conn *websocket.Conn) {
	conn.SetPongHandler(func(string2 string) error {
		conn.SetReadDeadline(time.Now().Add(time.Second * 60))
		return nil
	})
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if err := conn.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
			conn.Close()
			break
		}
	}
}

func removeConnection(groupId int64, uid int64) {
	groupConnsMux.Lock()
	defer groupConnsMux.Unlock()
	if conns, ok := groupConns[groupId]; ok {
		delete(conns, uid)
		if len(conns) == 0 {
			delete(groupConns, groupId)
		}
	}
}

func packingMsg(message string, groupId int64, uid int64) *repository.Group {
	msg := &repository.Group{
		GroupId:   groupId,
		FromUid:   uid,
		Content:   message,
		Timestamp: time.Now().Unix(),
	}
	return msg
}
