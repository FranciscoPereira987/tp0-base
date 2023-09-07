
import logging
from common.protocol.message import Message
from common.protocol.bet_message import BetMessage
from common.protocol.betbatch_message import BetBatchMessage
from common.utils import store_bets
from common.results import ErrResult, ExecutionResult, OkResult

class BetReader(Message):

    def __init__(self, lock) -> None:
        self.__bets = []
        self.__result = ErrResult()
        self.__lock = lock
        
    def deserialize(self, stream: bytes) -> bool:
        """
            If the stream was proccessed correctly sets the execution result to Ok
            subsequent streams do not change the result once it was changed to OkResult
        """
        result = self.__process_stream(stream)
        if result:
            self.__result = OkResult()
        return result
    
    def serialize(self) -> bytes:
        return super().serialize()
    
    def process_bets(self) -> bool:
        """
            Stores the readed bets
        """
        try:
            self.__store_bets()
            for bet in self.__bets:
                logging.info(f'action: apuesta_almacenada | result: success | dni: {bet.document} | number: {bet.number}')
            return True
        except Exception as e:
            logging.info(f"action: apuesta_almacenada | result: fail | error: {e}")
            return False
        finally:
            self.__lock.release()

    def executed(self, agency: int) -> ExecutionResult:
        self.__result.agency = agency
        return self.__result

    def __process_stream(self, stream: bytes) -> bool:
        """
            Processes a stream of bytes, storing the parsed bets
            the stream can be composed of either a single bet or a batch
        """
        self.__bets = []
        result = False
        if len(stream) < Message.HEADER_SIZE:
            return False
        
        if stream[:1] == BetMessage.BET_OP:
            bet_message = BetMessage()
            result = bet_message.deserialize(stream)
            self.__bets = [bet_message.bet]
            
        elif stream[:1] == BetBatchMessage.BETBATCH_OP:
            message = BetBatchMessage()
            result = message.deserialize(stream)
            self.__bets = [bet.bet for bet in message.bets]
        
        return result and self.process_bets()
    
    def __store_bets(self):
        self.__lock.acquire()
        store_bets(self.__bets)
        

