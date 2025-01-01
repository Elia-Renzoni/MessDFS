#
#   Login Request Handler
#

class Login:
    def __init__(self):
        pass
    
    def handle_login_req(conn, client_addr):
        print("Hey")
        conn.send(bytes("HTTP/1.1 200 OK\n\nHello World".encode()))
        conn.close()