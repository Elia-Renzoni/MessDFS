#
#   HTTP Router
#

from http.server import BaseHTTPRequestHandler, HTTPServer
import sys
#import psycopg2
from urllib.parse import urlparse, parse_qs

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
            case "/signup":
                self.signup()
            case "/add-friend":
                self.add_friend()
            case "/add-directory":
                self.add_directory()
            case _ :
                self.write_bad_request()
    
    def do_GET(self):
        match self.path:
            case "/login": 
                self.login()
            case "/friendship": 
                self.check_friendship()
            case "/directories": 
                self.get_directories()
            case "/friends":
                self.get_friends()
            case "/search-friend":
                self.search_buddy()
            case _: 
                self.write_bad_request()
    
    def do_DELETE(self):
        match self.path:
            case "/signout": 
                self.signout()
            case _: 
                self.write_bad_request()

    
    def login(self):
        print("login", file=sys.stdout)
        # http:127.0.0.1/login?username=pippo&password="bar"    
        print(self.path, file=sys.stdout)
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
    def add_friend(self):
        print("add-friend", file=sys.stdout)
    
    def add_directory(self):
        print("add-directory", file=sys.stdout)


    def check_friendship(self): 
        print("i'm checking the friendship", file=sys.stdout)    
    
    def get_directories(self):
        pass
    
    def get_friends(self): 
        pass

    def search_buddy(self):
        pass

    def write_bad_request(self):
        self.send_response(400)
        self.wfile.write(b"You Wrote a Bad Request Buddy!")
    
    def connect_db(self, sql_statement):
        connection = None
        try:
            connection = psycopg2.connect(
                dbname="messdfs",
                user="elia",
                password="elia",
                host="localhost"
            )

            cursor_obj = connection.cursor()
            cursor_obj.execute(sql_statement)
            sql_results = cursor_obj.fetchall()
        except: 
            print('Occured Some Errors')
        finally:
            connection.close()