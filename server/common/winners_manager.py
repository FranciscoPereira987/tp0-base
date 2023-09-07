import logging
from common.protocol.message import Message
from common.protocol.err_message import ErrMessage
from common.protocol.message import Message
from common.protocol.message import Message
from common.utils import Bet, has_won, load_bets, winners_by_agency
from common.protocol.winners_message import WinnersMessage
from common.protocol.winners_response_message import WinnersResponseMessage

class WinnersManager():

    def __init__(self, lock):
        self.__winners = None
        self.__handler = WinnersErrHandler()
        self.__lock = lock

    def get_handler(self, missing: int) -> 'WinnersHandler':
        try:
            self.__lock.acquire()
            if missing == 0 and not self.__winners:
                self.__winners = winners_by_agency()
                self.__handler = WinnersRespHandler(self.__winners)
        finally:
            self.__lock.release()
        return self.__handler

class WinnersHandler():
    
    def check_winners(self, stream: bytes) -> bool:
        winners = WinnersMessage()
        return winners.deserialize(stream)

    def handle_winners(self, client_id: int) -> Message:
        pass

class WinnersErrHandler(WinnersHandler):

    def handle_winners(self, client_id: int) -> Message:
        return ErrMessage()   
    

class WinnersRespHandler(WinnersHandler):

    def __init__(self, winners_map: dict[int, list[Bet]]) -> None:
        self.__map = winners_map

    def handle_winners(self, client_id: int) -> Message:
        bets = self.__map.get(client_id, [])
        
        return WinnersResponseMessage(bets)
        