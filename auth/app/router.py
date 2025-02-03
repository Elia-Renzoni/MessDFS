#
#   HTTP Router
#

from http.server import BaseHTTPRequestHandler, HTTPServer
import sys
import psycopg2
import json
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
        parsed_url = urlparse(self.path)
        parsed_query = parse_qs(parsed_url.query)
        username = parsed_query.get("username")
        password = parsed_query.get("password")
        result = self.connect_db("SELECT username, password FROM users WHERE username = %s AND password = %s", (username, password))
        if len(result) == 0:
            self.send_response(400)
        else:
            self.send_response(200)
    
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
        #http://127.0.0.1/friendship?txn=foo&dir=bar
        parsed_url = urlparse(self.path)
        parsed_query = parse_qs(parsed_url.query)
        txn = parsed_query.get("txn")
        directory = parsed_query.get("dir")
        result = self.connect_db("SELECT username, friend FROM friends WHERE username = %s AND friend = %s", (txn, directory))
        if len(result) == 0:
            response = {
                "result": False,
            }
            self.send_response(400)
            self.wfile.write(json.dumps(response))
        else: 
            response = {
                "result": True,
            }
            self.send_response(200)
            self.wfile.write(json.dumps(response))
    
    def get_directories(self):
        #http://127.0.0.1/directories?username=foo
        parsed_url = urlparse(self.path)
        parsed_query = parse_qs(parsed_url.query)
        username = parsed_query.get("username")
        result = self.connect_db("SELECT directory FROM directories WHERE username=%s", username)
        if len(result) == 0:
            self.send_response(400)
        else:
            resp = {
                "directories": result,
            }
            self.send_response(200)
            self.wfile.write(json.dumps(resp))
    
    def get_friends(self): 
        #http://127.0.0.1?friends?username=foo
        parsed_url = urlparse(self.path)
        parsed_query = parse_qs(parsed_url.query)
        username = parsed_query.get("username")
        result = self.connect_db("SELECT friend FROM friends WHERE username=%s", username)
        if len(result) == 0:
            self.send_response(400)
        else: 
            respo = {
                "friends": result,
            }
            self.send_response(200)
            self.wfile.write(json.dumps(respo))


    def search_buddy(self):
        #http://127.0.0.1/search_friend?username=foo
        parsed_url = urlparse(self.path)
        parsed_query = parse_qs(parsed_url.query)
        username = parsed_query.get("username")
        result = self.connect_db("SELECT username FROM users WHERE username=%s", username)
        if len(result) == 0:
            self.send_response(404)
        else: 
            response = {
                "user": result,
            }
            self.send_response(200)
            self.wfile.write(json.dumps(response))

    def write_bad_request(self):
        self.send_response(400)
        self.wfile.write(b"You Wrote a Bad Request Buddy!")
    
    def connect_db(self, sql_statement, values):
        connection = None
        sql_results = None
        try:
            connection = psycopg2.connect(
                dbname="messdfs",
                user="elia",
                password="elia",
                host="localhost"
            )

            cursor_obj = connection.cursor()
            cursor_obj.execute(sql_statement, values)
            sql_results = cursor_obj.fetchall()
        except: 
            print('Occured Some Errors')
        finally:
            connection.close()
        return sql_results