#
#   Login Request Handler
#

class Login:
    def __init__(self):
        pass
    
    def handle_login_req(data, conn, client_addr):
        print("Hey")

        splitted_request = data.split("\r\n")

        body = splitted_request[len(splitted_request) - 1]
        print(body)

        conn.send(bytes("HTTP/1.1 200 OK\n\nHello World".encode()))
        conn.close()