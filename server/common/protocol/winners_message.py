from common.protocol.message import Message


class WinnersMessage(Message):
    WINN_OP = bytes([0x06])
    

    def serialize(self) -> bytes:
        return self._build_header(self.WINN_OP, 0)
    
    def should_ack(self) -> bool:
        return False