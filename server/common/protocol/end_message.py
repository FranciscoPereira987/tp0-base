from common.protocol.message import Message


class EndMessage(Message):
    END_OP = bytes([0xff])

    def serialize(self) -> bytes:
        return self._build_header(self.END_OP, 0)