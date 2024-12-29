#
#   HTTP Router
#

import socket
import threading
import login
import signout
import signup

PORT = 8082
HOST = '172.0.0.1'

class Router:
    def __init__(self):
        self.login = login.Login()
        self.signout = signout.Signout()
        self.signup = signup.Signup()

    def start_handler(self):
        pass
