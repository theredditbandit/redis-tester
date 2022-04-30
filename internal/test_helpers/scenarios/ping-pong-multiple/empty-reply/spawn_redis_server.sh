#!/usr/bin/env -S python3 -u
import socket
import time
print("hey, binding to 6379")
sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)

sock.bind(('', 6379))

sock.listen(1)

conn, cli_addr = sock.accept()
conn.sendall(b'+PONG\r\n\n') # The extra \n would be sent if a user uses Println.
time.sleep(5)
