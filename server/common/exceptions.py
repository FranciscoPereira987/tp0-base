from common.protocol.message import Message


class CloseException(Exception):

    def __init__(self) -> None:
        super().__init__("Connection_closed")

class MalformedMessageException(Exception):

    def __init__(self) -> None:
        super().__init__("Recieved invalid message")

class UnexpectedMessageException(Exception):
    def __init__(self, message: bytes) -> None:
        super().__init__(f"Unexpected message recieved: {message[:Message.HEADER_SIZE]}")

class ErrRecievedException(Exception):

    def __init__(self) -> None:
        super().__init__("Recived Error Message")

class BrokenConnectionException(Exception):

    def __init__(self) -> None:
        super().__init__("Connection broken")
