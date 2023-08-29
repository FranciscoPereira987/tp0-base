from common.protocol.message import Message


class ErrMessage(Message):
    ERR_OP = bytes([0x03])
    

    def serialize(self) -> bytes:
        return self._build_header(self.ERR_OP, 0)
    
    def should_ack(self) -> bool:
        return False