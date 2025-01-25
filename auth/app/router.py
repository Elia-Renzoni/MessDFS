#
#   HTTP Router
#

from http.server import BaseHTTPRequestHandler, HTTPServer
import sys

class Router(BaseHTTPRequestHandler):
    # in case of HTTP 1.1
    def do_OPTIONS(self):
        self.send_response(200)
        self.send_header("Access-Control-Allow-Origin", "*")  
        self.send_header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        self.send_header("Access-Control-Allow-Headers", "Content-Type")
        self.end_headers()
    
    def do_POST(self):
        match self.path:
            case "/signout":
                self.signout()
            case "/login":
                self.login()
            case "/signup":
                self.signup()
    
    def login(self):
        print("login", file=sys.stdout)
        self.send_response(200)  # Codice HTTP 200 OK
        self.send_header("Content-Type", "text/plain")
        self.end_headers()
        self.wfile.write(b"Login endpoint reached")
    
    def signup(self):
        print("signup", file=sys.stdout)
        self.send_response(200)  # Codice HTTP 200 OK
        self.send_header("Content-Type", "text/plain")
        self.end_headers()
        self.wfile.write(b"Signup endpoint reached")
    
    def signout(self):
        print("signout", file=sys.stdout)
        self.send_response(200)  # Codice HTTP 200 OK
        self.send_header("Content-Type", "text/plain")
        self.end_headers()
        self.wfile.write(b"Signout endpoint reached")
    
