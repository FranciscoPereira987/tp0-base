import logging
import socket
from common.protocol.ack_message import AckMessage

from common.protocol.hello_message import HelloMessage
from common.protocol.message import Message


class BetConnListener():
    
    # Initialices a BetConn that can listen for connections
    def __init__(self, port: int, listen_backlog: int) -> None:
        self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.socket.bind(('', port))
        self.socket.listen(listen_backlog)

    def accept(self) -> ('BetConn', str):
        new_conn, addr = self.socket.accept()
        
        return BetConn(new_conn), addr
    
    def close(self):
        #TODO: Add log entries
        self.socket.shutdown(socket.SHUT_RDWR)
        self.socket.close()

class BetConn():
    
    def __init__(self, conn: socket.socket) -> None:
        self.socket = conn
        #TODO: Fix issues on connection
        self.id = 0
        self.__accept_connection()

    def __peak(self) -> bytes:
        return self.__read_bytes(4)
    
    def __read_bytes(self, size: int) -> bytes:
        #TODO: Fix short read
        return self.socket.recv(size)

    def read(self, message: Message) -> bool:
        header = self.__peak()
        length = message.get_needed_size(header)
        body = self.__read_bytes(length-4)
        
        result = message.deserialize(header + body)
        if result and message.should_ack():
            
            ack = AckMessage()
            result = self.write(ack)
        return result

    def write(self, message: Message) -> bool:
        serialized = message.serialize()
        #TODO: Fix short write
        _ = self.socket.send(serialized)
        if message.should_ack():
            ack = AckMessage()
            return self.read(ack)
        return True
    
    def close(self):
        #TODO: Add log entries
        self.socket.shutdown(socket.SHUT_RDWR)
        self.socket.close()

    def __accept_connection(self) -> bool:
        message = HelloMessage(0)
        logging.info("action: waiting_hello | result: in_progress ")
        if not self.read(message):
            logging.error("action: waiting_hello | result: failed ")
            return -1
        self.id = message.id
        ack = AckMessage()
        return self.write(ack)
        
        
        
