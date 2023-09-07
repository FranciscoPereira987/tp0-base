
import unittest

from common.protocol.hello_message import HelloMessage
from common.protocol.ack_message import AckMessage
from common.protocol.err_message import ErrMessage
from common.protocol.bet_message import BetMessage
from common.protocol.end_message import EndMessage
from common.protocol.betbatch_message import BetBatchMessage
from common.protocol.winners_message import WinnersMessage
from common.protocol.winners_response_message import WinnersResponseMessage
from common.utils import Bet

class TestProtocol(unittest.TestCase):
  
    def test_hello_message_serialization(self):
        message = HelloMessage(3)
        expected = HelloMessage.HELLO_OP + bytes([0, 0, 8, 0, 0, 0, 3])

        result = message.serialize()

        self.assertListEqual(list(expected), list(result))

    def test_hello_message_deserialization(self):
        expected_id = 3
        message = HelloMessage(expected_id)

        serialized = message.serialize()
        result = HelloMessage(0)
        result.deserialize(serialized)

        self.assertEquals(expected_id, result.id)

    def test_ack_message_serialization(self):
        message = AckMessage()
        expected = AckMessage.ACK_OP + bytes([0, 0, 4])

        result = message.serialize()
        self.assertListEqual(list(expected), list(result))

    def test_ack_message_deserialization(self):
        message = AckMessage()
        serialized = message.serialize()

        self.assertTrue(message.deserialize(serialized))

    def test_err_message_serialization(self):
        message = ErrMessage()
        expected = ErrMessage.ERR_OP + bytes([0, 0, 4])

        result = message.serialize()

        self.assertListEqual(list(expected), list(result))

    def test_err_message_deserialization(self):
        message = ErrMessage()
        serialized = message.serialize()

        self.assertTrue(message.deserialize(serialized))

    def test_winners_message_serialization(self):
        message = WinnersMessage()
        expected = WinnersMessage.WINN_OP + bytes([0, 0, 4])

        result = message.serialize()

        self.assertListEqual(list(expected), list(result))

    def test_winners_message_deserialization(self):
        message = WinnersMessage()
        serialized = message.serialize()

        self.assertTrue(message.deserialize(serialized))

    def test_end_message_serialization(self):
        message = EndMessage()
        expected = EndMessage.END_OP + bytes([0, 0, 4])

        result = message.serialize()

        self.assertListEqual(list(expected), list(result))

    def test_end_message_deserialization(self):
        message = EndMessage()
        serialized = message.serialize()

        self.assertTrue(message.deserialize(serialized))

    def test_bet_message_serialization(self):
        message = BetMessage("2", "Francisco", "Pereira", "41797243", "1998-12-17", "12345")
        expected = BetMessage.BET_OP + bytes([0, 0, 50])

        result = message.serialize()[:4]

        self.assertListEqual(list(expected), list(result))

    def test_bet_message_deserialization(self):
        message = BetMessage("3", "Francisco", "Pereira", "41797243", "1998-12-17", "12345")
        serialized = message.serialize()

        result = BetMessage()
        
        result.deserialize(serialized)

        self.assertEqual(message, result, f"\n{message.bet} \nvs\n {result.bet}")
        

    def test_betbatch_message_serialization(self):
        message = BetBatchMessage()
        bet = BetMessage("5", "Francisco", "Pereira", "41797243", "1998-12-17", "12345")
        message.add_bets(8 * [bet])

        expected = BetBatchMessage.BETBATCH_OP + bytes([0, 1, 148])

        result = message.serialize()

        self.assertListEqual(list(result[:4]), list(expected))

    def test_betbatch_message_deserialization(self):
        message = BetBatchMessage()
        bet = BetMessage("6", "Francisco", "Pereira", "41797243", "1998-12-17", "12345")
        message.add_bets(20 * [bet])
        serialized = message.serialize()

        result = BetBatchMessage()
        result.deserialize(serialized)

        self.assertEquals(message, result)

    def test_winners_response_message_serialization(self):
        message = WinnersResponseMessage([Bet('0', "Francisco", "Pereira", "41797243", "1998-12-17", "12345")])
        expected = message.WINNRESP_OP + bytes([0, 0, 13])

        result = message.serialize()[:4]

        self.assertListEqual(list(expected), list(result))

    def test_bet_message_deserialization(self):
        message = WinnersResponseMessage([Bet('0', "Francisco", "Pereira", "41797243", "1998-12-17", "12345")])
        serialized = message.serialize()

        result = WinnersResponseMessage()
        
        result.deserialize(serialized)

        self.assertEqual(message, result)


if __name__ == "__main__":
    unittest.main()