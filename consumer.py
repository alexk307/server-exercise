import socket
import struct
import base64
import json
from requests import post, RequestException

PORT = 6789
MAX_SIZE_UDP = 65535
HEADER_SIZE = 12

SERVER = 'http://localhost:1234/add'


def main():
    # Create the UDP socket
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

    # Set socket with SO_REUSEPORT flag which allows for multiple
    # listeners on the same port. See SO_REUSEPORT at,
    # http://man7.org/linux/man-pages/man7/socket.7.html
    sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEPORT, 1)

    # Bind to the local loopback
    sock.bind(('0.0.0.0', PORT))

    while True:
        data, addr = sock.recvfrom(MAX_SIZE_UDP)

        # Big endian 2 unsigned chars and an unsigned short
        header = data[:4]
        flag1, flag2, payload_size = struct.unpack('>2B1H', header)

        # Big endian one unsigned long
        offset = data[4:8]
        offset_val = struct.unpack('>L', offset)[0]

        # Big endian one unsigned long
        tran = data[8:12]
        trans_id = struct.unpack('>L', tran)[0]

        # Data is the rest of the message
        payload = data[HEADER_SIZE:]

        # Assert we received data and it's the correct size
        assert payload != ''
        assert payload_size + HEADER_SIZE == len(data)

        fragment = {
            'offset': offset_val,
            'trans_id': trans_id,
            'payload': base64.b64encode(payload),
            'size': payload_size
        }

        try:
            post(SERVER, json.dumps(fragment))
            print 'sent'
        except RequestException:
            print "[ERROR] Failed to send to database"
            continue

if __name__ == '__main__':
    main()
