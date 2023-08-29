from common.protocol.bet_message import BetMessage
from common.protocol.message import Message


class BetBatchMessage(Message):
    BETBATCH_OP = bytes([0x05])
    
    def __init__(self):
        self.bets = []

    def add_bets(self, bets: list):
        self.bets += bets

    def __eq__(self, __value: object) -> bool:
        try:    
            result = len(self.bets) == len(__value.bets) and all(map(lambda x: x[0] == x[1], zip(self.bets, __value.bets)))
        except:
            result = False
        finally:
            return result
        

    def serialize(self) -> bytes:
        serialized_bets = bytes()
        for bet in self.bets:
            serialized_bets += bet.serialize()
        header = self._build_header(self.BETBATCH_OP, len(serialized_bets))

        return header + serialized_bets
    
    def __deserialize_bet(self, stream: bytes) -> (bytes, bool):
        betsize = self._get_message_length(stream)
        if betsize < 0:
            return bytes(), False
        new_bet = BetMessage()
        ok = new_bet.deserialize(stream[:betsize])
        if not ok:
            return bytes(), False
        self.bets.append(new_bet)
        return stream[betsize:], ok
    
    def deserialize(self, stream: bytes) -> bool:
        if not self._check_header(stream, self.BETBATCH_OP):
            return False
        ok = True
        stream = stream[self.HEADER_SIZE:]
        while len(stream) > 0 and ok:
            stream, ok = self.__deserialize_bet(stream)
            

        return ok
    
