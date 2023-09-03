import socket
import logging
import signal

from common.bet_conn import BetConn, BetConnListener
from common.protocol.bet_message import BetMessage
from common.exceptions import BrokenConnectionException, CloseException
from common.utils import store_bets
from common.bet import BetReader

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = BetConnListener(port, listen_backlog)
        self.running = True
        self.__set_shutdown()

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """
        # TODO: Modify this program to handle signal to graceful shutdown
        # the server
        while self.running:
            client_sock = self.__accept_new_connection()
            if self.running:
                self.__handle_client_connection(client_sock)

    def __handle_client_connection(self, client_sock: BetConn):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        bet = BetReader()
        while client_sock.active:
            try:
                client_sock.read(bet)
            except OSError as e:
                logging.error("action: receive_message | result: fail | error: {e}")
            except CloseException as e:
                logging.error(f"action: recieve_message | result: failed | connection: {e}")
            except BrokenConnectionException as e:
                logging.error(f"action: recieve_message | result: failed | connection: {e}")
            finally:
                pass
        client_sock.close()

    def __accept_new_connection(self) -> BetConn:
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        try:
            c, addr = self._server_socket.accept()
        except OSError:
            return
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c
    
    def __set_shutdown(self):
        def sigterm_handle( _s, _f):
            logging.info('action: SIGTERM | result: in_progress')
            self.__close_server_socket(_s, _f)
            logging.info('action: SIGTERM | result: success')
            
        signal.signal(signal.SIGTERM, sigterm_handle)

    def __close_server_socket(self, _s, _f):
        logging.info('action: closing_server_socket | result: in_progress')
        self._server_socket.close()
        self.running = False
        logging.info('action: closing_server_socket | result: success')
        
