#
#   HTTP Router
#

import socket
import threading
import login
import signout
import signup
import json

class Router:
    def __init__(self, host, listen_port):
        self.host = host
        self.port = listen_port
        self.login = login.Login()
        self.signout = signout.Signout()
        self.signup = signup.Signup()
        self.listen = socket.create_server((self.host, self.port))

    def start_handler(self):
        print("Server Listening...")
        thread = None
        while True:
            conn, client_addr = self.listen.accept()
            data = conn.recv(2024).decode()
            parsed_data = data.split('\r\n')
            # localhost:8080/insert/pippo
            url = parsed_data[0].split("/")
            # TODO: add parameters to the threadssss
            match url[1]:
                case "signout":
                   thread = threading.Thread(target=signout.Signout.handle_signout_requests) 
                case "login":
                    thread= threading.Thread(target=login.Login.handle_login_req)
                case "signup":
                    thread = threading.Thread(target=signup.Signup.handle_signup_request)
                case _:
                    payload = {
                        "err": "invalid endpoint name"
                    }
                    conn.send(json.dump(payload))
            thread.start()


    
