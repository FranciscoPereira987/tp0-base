import logging
import signal
from common.bet_conn import BetConn
import multiprocessing as mp

from common.bet import BetReader
from common.exceptions import CloseException
from common.results import ErrResult


class Runner():
    """
        Runs a BetConn connection on a process
    """
    def __init__(self, bet_conn: BetConn, queue: mp.SimpleQueue, lock) -> None:
        self.__conn = bet_conn
        self.__queue = queue
        self.__result = ErrResult(self.__conn.id)
        self.__reader = BetReader(lock)


    def run(self) -> mp.Process:
        """
            Sets up the process that's going to run and returns the process 
        """
        process = mp.Process(target=self.__process)
        process.start()
        return process

    def __run(self):
        
        while self.__conn.active:
            try:
                self.__conn.read(self.__reader)
            except CloseException as e:
                self.__result = self.__reader.executed(self.__conn.id)
                logging.error(f"action: recieve_message | result: conection_closed | error: {e}")
            except Exception as e:
                self.__conn.close()
                logging.error(f"action: recieve_message | result: failed | error: {e}")

        self.__queue.put(self.__result)

    def __process(self):
        self.__set_up_sigterm()
        self.__run()

    def __set_up_sigterm(self):
        signal.signal(signal.SIGTERM, self.__handle_sigterm)

    def __handle_sigterm(self, _s, _f):
        logging.info(f"action: shut_down_worker[{self.__conn.id}] | result: in_progress")
        self.__conn.close()
        logging.info(f"action: shut_down_worker[{self.__conn.id}] | result: success")
        


    