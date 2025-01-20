#
#   Signup Requests Handler
#

class Signup:

    def __init__(self):
        pass

    def handle_signup_request(conn, client_address):
        print("Hey Signup")
        conn.send(bytes("HTTP/1.1 200 OK\n\nHello From Signup".encode()))
        conn.close()