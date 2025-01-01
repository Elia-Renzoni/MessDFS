#
#   HTTP Router
#

import socket
import threading
import app.login as login
import app.signout as signout
import app.signup as signup

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
        while True:
            conn, client_addr = self.listen.accept()
            data = conn.recv(2024).decode()
            parsed_data = data.split('\r\n')
            url = parsed_data[0].split("/")
            print(url)
            endpoint = url[1].split(" ")
            match endpoint[0]:
                case "signout":
                   thread = threading.Thread(target=signout.Signout.handle_signout_requests, args=(conn, client_addr)) 
                   thread.start()
                case "login":
                    thread= threading.Thread(target=login.Login.handle_login_req, args=(conn, client_addr))
                    thread.start()
                case "signup":
                    thread = threading.Thread(target=signup.Signup.handle_signup_request, args=(conn, client_addr))
                    thread.start()
                case _:
                    response = "HTTP/1.1 500 INTERNAL SERVER ERROR\n\n" 
                    conn.send(bytes(response.encode()))
                    conn.close()