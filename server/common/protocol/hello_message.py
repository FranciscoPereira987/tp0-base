from common.protocol.message import Message


class HelloMessage(Message):
    HELLO_OP = bytes([0x01])
    
    def __init__(self, id: int) -> None:
        super().__init__()
        self.id = id

    def serialize(self) -> bytes:
        stream = self.HELLO_OP
        stream = self._build_header(stream, 4)
        stream += self._serialize_uint32(self.id)

        return stream
    
    def deserialize(self, stream: bytes) -> bool:
        if len(stream) != 8:
            return False
        
        stream = stream[self.HEADER_SIZE:]
        clientID = self._deserialize_uint32(stream)
        self.id = clientID

        return clientID != -1
        

        
