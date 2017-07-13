from requests import post
from random import randrange
from uuid import uuid4
import base64
import json

PORT = 6789
MAX_SIZE_UDP = 65535
HEADER_SIZE = 12
NUM_TRANSACTIONS = 10

SERVER = 'http://localhost:1234/add'


def main():
    for i in range(NUM_TRANSACTIONS):

        # Psuedo-random transaction ID
        transaction_id = randrange(1, 100)
        payload = str(uuid4())

        # Break into random pieces pieces
        l = range(1000)
        pieces = randrange(1, 100)
        chunks = [l[i:i + pieces] for i in xrange(0, len(l), pieces)]

        for chunk in chunks:

            fragment = {
                'offset': chunk[-1],
                'trans_id': transaction_id,
                'payload': base64.b64encode(payload),
                'size': len(chunk)
            }

            post(SERVER, json.dumps(fragment))

if __name__ == '__main__':
    main()
