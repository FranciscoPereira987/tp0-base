package connection

import (
	"time"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/protocol"
	log "github.com/sirupsen/logrus"
)

type WinnersConn struct {
	addr string
	agency int
	backoffTime time.Duration
}

func  NewWinnersConn(addr string, agency int, backoff time.Duration) *WinnersConn {
	
	return &WinnersConn{
		addr,
		agency,
		backoff,
	}
}

func (wc WinnersConn) connect() (*BetConn, error) {
	return NewBetConn(wc.addr, wc.agency)
}


func (wc WinnersConn) trySendWinners(connection *BetConn) error {
	winners := new(protocol.Winners)
	return connection.Write(winners)
}

func (wc WinnersConn) recoverWinners(connection *BetConn) (*protocol.WinnersResponse, error) {
	winners := new(protocol.WinnersResponse)
	return winners, connection.Read(winners)
}

func (wc WinnersConn) queryWinners(connection *BetConn) (*protocol.WinnersResponse, error) {
	err := wc.trySendWinners(connection)
	if err != nil {
		return nil, err
	}
	return wc.recoverWinners(connection)
}

func (wc *WinnersConn) backoff() {
	time.Sleep(wc.backoffTime)
	wc.backoffTime *= 2
}  

func (wc *WinnersConn) WaitForWinners() (*protocol.WinnersResponse, error){
	connection, err := wc.connect()
	var winners *protocol.WinnersResponse
	if err != nil {
		log.Infof("activity: consulta_ganadores | result: failed | err: %s", err)
		return nil, err
	}
	
	for  winners, err = wc.queryWinners(connection) ;err != nil;  {
		log.Infof("activity: consulta_ganadores | result: in_progress | err: server_not_ready")
		connection.shutdown()
		wc.backoff()
		connection, err = wc.connect()
		if err != nil {
			log.Infof("activity: consulta_ganadores | result: failed | err: %s", err)
			return nil, err
		}
		winners, err = wc.queryWinners(connection)
	}
	connection.shutdown()
	return winners, err
}