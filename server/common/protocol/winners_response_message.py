from common.protocol.message import Message
from common.utils import Bet


class WinnersResponseMessage(Message):
    WINNRESP_OP = bytes([0x07])

    def __init__(self, bets: list[Bet] = []) -> None:
        self.winners = list(map(lambda bet: bet.document, bets))
    
    def __eq__(self, __value: object) -> bool:
        try:
            result = all(map(lambda x: x[0]==x[1], zip(self.winners, __value.winners)))
        except Exception:
            result = False
        return result
            

    def __serialize_document(self, document: str) -> bytes:
        serialized_document = bytes([len(document)])
        serialized_document += document.encode()
        return serialized_document

    def serialize(self) -> bytes:
        body = bytes()

        for document in self.winners:
            body += self.__serialize_document(document)

        return self._build_header(self.WINNRESP_OP, len(body)) + body
    
    def __deserialize_field(self, field_length: int, stream: bytes) -> (str, int):
        if len(stream) < field_length:
            
            return "", -1
        
        field = stream[:field_length].decode()
        return field, field_length

    def __deserialize_field_length(self, stream: bytes) -> int:
        if len(stream) == 0:
            return -1
        return int.from_bytes(stream[:1], self.ENDIAN, signed=False)

    def __get_field_from_stream(self, stream: bytes) -> (str, int):
        field_length = self.__deserialize_field_length(stream)
        if field_length < 0:
            
            return "", field_length
        stream = stream[1:]
        field, err = self.__deserialize_field(field_length, stream)
        return field, err
    
    def deserialize(self, stream: bytes) -> bool:
        if not self._check_header(stream, self.WINNRESP_OP):
            return False

        stream = stream[self.HEADER_SIZE:]

        while len(stream) > 0:
            document, size = self.__get_field_from_stream(stream)
            if size < 0:
                return False
            self.winners.append(document)
            stream = stream[size:]

        return True

    def should_ack(self) -> bool:
        return False