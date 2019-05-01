package main

import (
	go_redis_orm "github.com/fananchong/go-redis-orm.v2"
	"github.com/fananchong/go-xserver/config"
	"github.com/fananchong/go-xserver/services/internal/db"
)

// Lobby : 大厅服务器
type Lobby struct {
	*db.IDGen
	accountMgr *AccountMgr
}

// NewLobby : 构造函数
func NewLobby() *Lobby {
	lobby := &Lobby{
		accountMgr: NewAccountMgr(),
	}
	lobby.IDGen = &db.IDGen{}
	return lobby
}

// Start : 启动
func (lobby *Lobby) Start() bool {
	if lobby.initRedis() == false {
		return false
	}
	Ctx.EnableMessageRelay(true)
	Ctx.RegisterFuncOnRelayMsg(lobby.onRelayMsg)
	Ctx.RegisterFuncOnLoseAccount(lobby.onLoseAccount)
	return true
}

// Close : 关闭
func (lobby *Lobby) Close() {

}

func (lobby *Lobby) onRelayMsg(source config.NodeType, account string, cmd uint64, data []byte) {
	switch source {
	case config.Client:
		lobby.accountMgr.PostMsg(account, cmd, data)
	default:
		Ctx.Errorln("Unknown source, type:", source, "(", int(source), ")")
	}
}

func (lobby *Lobby) onLoseAccount(account string) {
	lobby.accountMgr.DelAccount(account)
}

func (lobby *Lobby) initRedis() bool {
	// db account
	err := go_redis_orm.CreateDB(
		Ctx.Config.DbAccount.Name,
		Ctx.Config.DbAccount.Addrs,
		Ctx.Config.DbAccount.Password,
		Ctx.Config.DbAccount.DBIndex)
	if err != nil {
		Ctx.Errorln(err)
		return false
	}
	lobby.IDGen.Cli = go_redis_orm.GetDB(Ctx.Config.DbAccount.Name)
	return true
}
