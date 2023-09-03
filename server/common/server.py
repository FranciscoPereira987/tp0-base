import socket
import logging
import signal

from common.bet_conn import BetConn, BetConnListener
from common.protocol.bet_message import BetMessage
from common.exceptions import BrokenConnectionException, CloseException
from common.utils import store_bets
from common.bet import BetReader
from common.runner import Runner
import multiprocessing as mp


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = BetConnListener(port, listen_backlog)
        self.running = True
        self.__queue = mp.SimpleQueue()
        self.__lock = mp.Lock()
        self.__workers = {}
        self.__set_shutdown()
        self.__set_sigchld()

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """
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
        runner = Runner(client_sock, self.__queue, self.__lock)
        handle = runner.run()

        self.__workers[client_sock.id] = handle

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
    
    def __set_sigchld(self):
        signal.signal(signal.SIGCHLD, self.__handle_sigchld)

    def __handle_sigchld(self, _s, _f):
        if not self.__queue.empty():
            result = self.__queue.get()

            if result.ok():
                self._server_socket.update_missing()
            self.__workers[result.agency].join()

    def __set_shutdown(self):
        def sigterm_handle( _s, _f):
            logging.info('action: SIGTERM | result: in_progress')
            self.__close_server_socket(_s, _f)
            logging.info('action: SIGTERM | result: success')
            
        signal.signal(signal.SIGTERM, sigterm_handle)

    def __close_server_socket(self, _s, _f):
        logging.info('action: closing_server_socket | result: in_progress')
        for worker in self.__workers:
            self.__workers[worker].terminate()
            self.__workers[worker].join()
        self._server_socket.close()
        self.running = False
        logging.info('action: closing_server_socket | result: success')
        
