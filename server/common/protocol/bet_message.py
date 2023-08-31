from common.protocol.message import Message
from common.utils import Bet


class BetMessage(Message):
    BET_OP = bytes([0x04])
    EXTRA_BET_BYTES = 8
    
    def __init__(self, first_name: str = "", last_name: str = "", document: str = "", birthdate: str = "1970-12-31", number: str = "0"):
        self.bet = Bet('0', first_name, last_name, document, birthdate, number)

    def __eq__(self, __value: object) -> bool:
        try:
            result = self.bet.agency == __value.bet.agency
            result &= self.bet.birthdate == __value.bet.birthdate
            result &= self.bet.document == __value.bet.document
            result &= self.bet.first_name == __value.bet.first_name
            result &= self.bet.last_name == __value.bet.last_name
            result &= self.bet.number == __value.bet.number
        except:
            result = False
        finally:
            return result
        
    def __make_header(self) -> bytes:
        length = self.EXTRA_BET_BYTES +\
            len(self.bet.first_name) +\
            len(self.bet.last_name) +\
            len(self.bet.document) +\
            len(self.bet.birthdate.isoformat())
        
        return self._build_header(self.BET_OP, length)

    def __add_field(self, field: str) -> bytes:
        stream = bytes([len(field)])
        stream += field.encode()
        return stream

    def __add_body(self, stream: bytes) -> bytes:
        stream += self.__add_field(self.bet.first_name)
        stream += self.__add_field(self.bet.last_name)
        stream += self.__add_field(self.bet.document)
        stream += self.__add_field(self.bet.birthdate.isoformat())
        stream += self._serialize_uint32(self.bet.number)

        return stream

    def serialize(self) -> bytes:
        serialized = self.__make_header()
        
        return self.__add_body(serialized)
    
    def __deserialize_field_length(self, stream: bytes) -> int:
        if len(stream) == 0:
            return -1
        return int.from_bytes(stream[:1], self.ENDIAN, signed=False)

    def __deserialize_field(self, field_length: int, stream: bytes) -> (str, int):
        if len(stream) < field_length:
            
            return "", -1
        
        field = stream[:field_length].decode()
        return field, field_length

    def __get_field_from_stream(self, stream: bytes) -> (str, int):
        field_length = self.__deserialize_field_length(stream)
        if field_length < 0:
            
            return "", field_length
        stream = stream[1:]
        field, err = self.__deserialize_field(field_length, stream)
        return field, err

    def deserialize(self, stream: bytes) -> bool:
        if not self._check_header(stream, self.BET_OP):
            
            return False
        
        stream = stream[self.HEADER_SIZE:]
        
        fields = ['0']
        for _ in range(4):
            field, field_size =\
                self.__get_field_from_stream(stream)
            if field_size < 0:
                
                return False
            fields.append(field)
            stream = stream[field_size+1:]

        bet_number = self._deserialize_uint32(stream)
        if bet_number < 0:
            return False
        fields.append(str(bet_number))
        
        self.bet = Bet(*fields)
        
        return True