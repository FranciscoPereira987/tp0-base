class Message():

    HEADER_SIZE = 4
    ENDIAN = "big"

    def serialize(self) -> bytes:
        pass

    def deserialize(self, stream: bytes) -> bool:
        return self._compare_streams(self.serialize(), stream)

    def should_ack(self) -> bool:
        return True
    
    def get_needed_size(self, stream: bytes) -> int:
        if len(stream) != 4:
            return -1
        
        return self._get_message_length(stream)
    
    def _build_header(self, header: bytes, size: int) -> bytes:
        size += self.HEADER_SIZE
        return header + size.to_bytes(3, self.ENDIAN, signed=False)

    def _compare_streams(self, a: bytes, b: bytes) -> bool:
        
        return len(a) == len(b) and all(map(lambda x: x[0] == x[1], zip(a, b)))

    def _get_message_length(self, stream: bytes) -> int:
        if len(stream) < self.HEADER_SIZE:
            return -1
        
        return int.from_bytes(stream[1:self.HEADER_SIZE], self.ENDIAN)

    def _check_header(self, stream: bytes, op_code: int) -> bool:
        length = self._get_message_length(stream)
        
        return length == len(stream) and stream[:1] == op_code
    
    def _deserialize_uint32(self, stream: bytes) -> int:
        
        if len(stream) != 4:
            return -1
        
        return int.from_bytes(stream, self.ENDIAN, signed=False)

    def _serialize_uint32(self, num: int) -> bytes:
        return num.to_bytes(4, self.ENDIAN, signed=False)