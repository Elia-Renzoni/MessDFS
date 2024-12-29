#
#   HTTP Router
#

import socket
import threading
import login
import signout
import signup

class Router:
    def __init__(self, host, listen_port):
        self.host = host
        self.port = listen_port
        self.login = login.Login()
        self.signout = signout.Signout()
        self.signup = signup.Signup()
        self.conn = socket.create_server((self.host, self.port))

    def start_handler(self):
        pass
