
import logging
from common.protocol.message import Message
from common.protocol.bet_message import BetMessage
from common.protocol.betbatch_message import BetBatchMessage
from common.utils import store_bets


class BetReader(Message):

    def __init__(self) -> None:
        self.bets = []

    def deserialize(self, stream: bytes) -> bool:
        return self.__process_stream(stream)
    
    def serialize(self) -> bytes:
        return super().serialize()
    
    def process_bets(self) -> None:
        """
            Stores the readed bets
        """
        store_bets(self.bets)
        for bet in self.bets:
            logging.info(f'action: apuesta_almacenada | result: success | dni: {bet.document} | number: {bet.number}')

    def __process_stream(self, stream: bytes) -> bool:
        """
            Processes a stream of bytes, storing the parsed bets
            the stream can be composed of either a single bet or a batch
        """
        self.bets = []
        result = False
        if len(stream) < Message.HEADER_SIZE:
            return False
        
        if stream[:1] == BetMessage.BET_OP:
            bet_message = BetMessage()
            result = bet_message.deserialize(stream)
            self.bets = [bet_message.bet]
            
        elif stream[:1] == BetBatchMessage.BETBATCH_OP:
            message = BetBatchMessage()
            result = message.deserialize(stream)
            self.bets = [bet.bet for bet in message.bets]
        
        return result
        