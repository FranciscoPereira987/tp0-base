from errno import ENOTCONN
import logging
import socket
from common.protocol.ack_message import AckMessage

from common.protocol.hello_message import HelloMessage
from common.protocol.message import Message
from common.exceptions import BrokenConnectionException, CloseException, ErrRecievedException, MalformedMessageException, UnexpectedMessageException
from common.protocol.end_message import EndMessage
from common.protocol.err_message import ErrMessage
from common.winners_manager import WinnersErrHandler, WinnersHandler, WinnersManager, WinnersRespHandler


class BetConnListener():
    
    # Initialices a BetConn that can listen for connections
    def __init__(self, port: int, listen_backlog: int) -> None:
        self.missing = listen_backlog
        self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.socket.bind(('', port))
        self.socket.listen(listen_backlog)
        self.__winners_manager = WinnersManager()

    def accept(self) -> ('BetConn', str):
        new_conn, addr = self.socket.accept()
        
        return BetConn(new_conn, self.__load_handler()), addr
    
    def update_missing(self):
        self.missing -= 1
    
    def close(self):
        #TODO: Add log entries
        self.socket.shutdown(socket.SHUT_RDWR)
        self.socket.close()

    def __load_handler(self) -> WinnersHandler:
        return self.__winners_manager.get_handler(self.missing)


class BetConn():
    
    def __init__(self, conn: socket.socket, handler: WinnersHandler) -> None:
        self.socket = conn
        self.id = 0
        self.active = True
        self.winners = handler
        self.__accept_connection()

    def read(self, message: Message) -> None:
        """
            Reads a message from the socket
        """
        readed = self.__read_message(message)
        self.__manage_recieved_data(message, readed)
        

    def write(self, message: Message) -> None:
        """
            Writes a message to the socket and validates and AckMessage if needed
        """
        serialized = message.serialize()
        self.__write_bytes(serialized)

        if message.should_ack():
            ack = AckMessage()
            self.read(ack)
        
    
    def close(self):
        """
            If the socket is active, close the connection
        """
        if self.active:
            logging.info(f"action: connection_close | result: in_progress")
            self.__end_connection()
            
    def __send_ack(self) -> None:
        ack = AckMessage()
        self.write(ack)

    def __send_err(self) -> None:
        err = ErrMessage()
        self.write(err)
   
    def __manage_recieved_data(self, expected: Message, recieved: bytes) -> None:
        """
            Raises an error if the message is not the one expected
            if the client has closed the connection by sending an end message
            it responds to the client and closes the socket.
        """
        result = expected.deserialize(recieved)

        if result and expected.should_ack():
            self.__send_ack()
        
        elif not result and self.__recover_end_message(recieved):
            raise CloseException()
        
        elif not result and self.__recover_err_message(recieved):
            self.__shutdown_connection()
            raise ErrRecievedException()
        
        elif not result and self.winners.check_winners(recieved):
            self.write(self.winners.handle_winners(self.id))
            self.__shutdown_connection()

        elif not result:
            self.__send_err()
            raise MalformedMessageException()
    
    def __end_connection(self) -> None:
        """
            Ends a connection with the client
        """
        try:
            self.__wait_end()
        except Exception as e:
            logging.error(f"action: connection_close | result: in_progress | error: {e}")
        finally:
            self.__shutdown_connection()
            
    def __recover_end_message(self, message: bytes) -> bool:
        """
            Tries recover an EndMessage from the byte stream
            If it does so, then closes the connection.
            Whenever an EndMessage is recovered, returns True
        """
        end = EndMessage()
        result = end.deserialize(message)
        if result:
            logging.info(f"action: connection_close  | result: in_progress | error: Recieved End message")
            self.__send_ack()
            self.__shutdown_connection()
        return result
    
    def __recover_err_message(self, message: bytes) -> bool:
        """
            Returns True if an ErrMessage is recovered from the stream
        """
        return ErrMessage().deserialize(message)
    
    def __recover_ack_message(self, message: bytes) -> bool:
        """
            Returns True if an AckMessage is recovered from the stream
        """
        return AckMessage().deserialize(message)

    def __wait_end(self) -> None:
        """
            Ends the connection on the Bet protocol side, if after sending the end
            message, an EndMessage is recovered instead, sends the ack to the client first.
        """
        end_message = EndMessage()
        
        self.__write_bytes(end_message.serialize())
        readed = self.__read_message(end_message)
        
        if end_message.deserialize(readed):
            self.__send_ack()
            readed = self.__read_message(end_message)

        if not self.__recover_ack_message(readed):
            raise UnexpectedMessageException(readed)

    def __shutdown_connection(self) -> None:
        """
            If the connection is active, then shutdowns the connection by closing the socket.
            If self.socket.shutdown raises an exception, i need to check for ENOTCONN, since this 
            means that the connection was already closed on the other side
        """
        if self.active:
            self.active = False
            try:
                self.socket.shutdown(socket.SHUT_RDWR)
            except OSError as e:
                #If the error is not ENOTCONN, then raise it again
                if e.errno != ENOTCONN:
                    raise e
            self.socket.close()
            logging.info(f"action: connection_close | result: success")

    def __accept_connection(self) -> bool:
        """
            Accepts a connection if the first message sent is a hello message
            if not, it closes the connection
        """
        message = HelloMessage(0)
        logging.info("action: waiting_hello | result: in_progress ")
        try:
            self.read(message)
            self.id = message.id
            logging.info(f"action: waiting_hello | result: success | client_id: {self.id}")
            return True
        except Exception as e:
            logging.error(f"action: waiting_hello | result: failed | error: {e}")
            self.__shutdown_connection()
            return False
    
    def __peak(self) -> bytes:
        """
            Tries to read a message header from the stream
        """
        return self.__read_bytes(Message.HEADER_SIZE)
    
    def __read_bytes(self, size: int) -> bytes:
        """
            Reads size bytes from the socket
            Ref: https://docs.python.org/3/howto/sockets.html
        """
        
        if size < 0:
            raise MalformedMessageException()
        chunks = []
        total_recv = 0
        while total_recv < size:
            chunk = self.socket.recv(size - total_recv)
            if not chunk:
                raise BrokenConnectionException()
            chunks.append(chunk)
            total_recv += len(chunk)
    
        return b''.join(chunks)

    def __read_message(self, message: Message) -> bytes:
        """
            Reads a message from the socket and returns the bytes readed
        """
        if not self.active:
            raise CloseException()
        
        header = self.__peak()
        length = message.get_needed_size(header)
        body = self.__read_bytes(length-4)
        return header + body

    def __write_bytes(self, message: bytes) -> None:
        """
            Write the bytes to the socket
            Ref: https://docs.python.org/3/howto/sockets.html
        """
        if not self.active:
            raise CloseException()
        
        total_sent = 0
        while total_sent < len(message):
            sent = self.socket.send(message[total_sent:])
            if sent == 0:
                raise BrokenConnectionException()
            total_sent += sent

        
