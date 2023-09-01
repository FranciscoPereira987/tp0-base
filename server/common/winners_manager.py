import logging
from common.protocol.message import Message
from common.protocol.err_message import ErrMessage
from common.protocol.message import Message
from common.protocol.message import Message
from common.utils import has_won, load_bets
from common.protocol.winners_message import WinnersMessage
from common.protocol.winners_response_message import WinnersResponseMessage


class WinnersHandler():
    
    def check_winners(self, stream: bytes) -> bool:
        winners = WinnersMessage()
        return winners.deserialize(stream)

    def handle_winners(self, client_id: int) -> Message:
        pass

class WinnersErrHandler(WinnersHandler):

    def handle_winners(self, client_id: int) -> Message:
        return ErrMessage().serialize()    
    

class WinnersRespHandler(WinnersHandler):

    def handle_winners(self, client_id: int) -> Message:
        mapped = map(lambda x: x if x.agency == client_id else None,load_bets())
        
        bets = filter(lambda x: x != None and has_won(x), mapped)
        
        return WinnersResponseMessage(bets)
        