#
#   Signout Requests Handler
#

class Signout:
    
    def __init__(self):
        pass

    def handle_signout_requests(conn, client_address):
        conn.send(bytes("HTTP/1.1 200 OK\n\nHello From Signout".encode()))
        conn.close()