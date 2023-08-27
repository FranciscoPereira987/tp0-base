from common.protocol.message import Message


class AckMessage(Message):
    ACK_OP = bytes([0x02])
    
    def serialize(self) -> bytes:
        return self._build_header(self.ACK_OP, 0)
    
    def should_ack(self) -> bool:
        return False