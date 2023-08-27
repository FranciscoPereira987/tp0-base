class Message():

    HEADER_SIZE = 4

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
        shift = 16
        while shift >= 0:
            header = header + bytes([(size>>shift) & 0xff])
            shift -= 8
        return header

    def _compare_streams(self, a: bytes, b: bytes) -> bool:
        
        return len(a) == len(b) and all(map(lambda x: x[0] == x[1], zip(a, b)))

    def _get_message_length(self, stream: bytes) -> int:
        if len(stream) < self.HEADER_SIZE:
            return -1
        
        length = 0
        for i in range(1, 4):
            length += int(stream[i]) << (8 * (3 - i))
        return length

    def _check_header(self, stream: bytes, op_code: int) -> bool:
        length = self._get_message_length(stream)
        
        return length == len(stream) and stream[:1] == op_code
    
    def _deserialize_uint32(self, stream: bytes) -> int:
        number = 0
        
        if len(stream) != 4:
            return -1
        
        for index, value in enumerate(stream):
            number |= int(value) << (8 * (3 - index))

        return number

    def _serialize_uint32(self, num: int) -> bytes:
        serialized = bytes()
        shift = 24
        while shift >= 0:
            serialized += bytes([(num>>shift) & 0xff])
            shift -= 8

        return serialized
