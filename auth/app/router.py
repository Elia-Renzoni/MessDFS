#
#   HTTP Router
#

from http.server import BaseHTTPRequestHandler, HTTPServer
import sys
import psycopg2
from psycopg2 import Error
import json
from urllib.parse import urlparse, parse_qs

TO_COMMIT = 1
READ = 0

class Router(BaseHTTPRequestHandler):
    def do_POST(self):
        parsed = urlparse(self.path)
        print(parsed)
        match parsed.path:
            case "/signup":
                self.signup()
            case "/add-friend":
                self.add_friend()
            case "/add-directory":
                self.add_directory()
            case _ :
                self.write_bad_request()
    
    def do_GET(self):
        parsed = urlparse(self.path)
        print(parsed.path)
        match parsed.path:
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
    
    def login(self):
        print("login", file=sys.stdout)
        parsed_url = urlparse(self.path)
        parsed_query = parse_qs(parsed_url.query)
        username = parsed_query.get("username")
        password = parsed_query.get("password")
        result = self.connect_db("SELECT username, password FROM users WHERE username = '%s' AND password = '%s'" % (username.pop(), password.pop()), READ)
        if len(result) == 0:
            self.send_response(400)
        else:
            self.send_response(200)
        self.end_headers()
    
    def signup(self):
        print("signup", file=sys.stdout)
        content_length = int(self.headers['Content-Length'])
        body = self.rfile.read(content_length)
        deser_data = json.loads(body.decode('UTF-8'))
        username = deser_data["username"]
        password = deser_data["password"]
        
        result = self.connect_db("INSERT INTO users (username, password) VALUES ('%s', '%s')" % (username, password), TO_COMMIT)
        if result != None:
            self.send_response(400)
            self.end_headers()
        else: 
            self.send_response(201)
            self.end_headers()
    

    def add_friend(self):
        print("add-friend", file=sys.stdout)
        content_length = int(self.headers['Content-Length'])
        body = self.rfile.read(content_length)
        deser_data = json.loads(body.decode('UTF-8'))
        username = deser_data["username"]
        friend_name = deser_data["friend_username"]

        result = self.connect_db("INSERT INTO friends (username, friend) VALUES ('%s', '%s')" % (username, friend_name), TO_COMMIT)
        print(result)
        if result != None: 
            self.send_response(400)
            self.end_headers()
        else: 
            self.send_response(201)
            self.end_headers()
        
    
    def add_directory(self):
        print("add-directory", file=sys.stdout)
        content_length = int(self.headers['Content-Length'])
        body = self.rfile.read(content_length)
        deser_data_add_directory = json.loads(body.decode('utf-8'))
        username = deser_data_add_directory["username"]
        directory_name = deser_data_add_directory["directory"]

        result = self.connect_db("INSERT INTO directories (username, directory) VALUES ('%s', '%s');" % (username, directory_name), TO_COMMIT)
        if result != None: 
            self.send_response(400)
            self.end_headers()
        else: 
            self.send_response(201)
            self.end_headers()        


    ''''
        TODO -> la query è sbagliata
    '''
    def check_friendship(self): 
        print("i'm checking the friendship", file=sys.stdout)    
        #http://127.0.0.1/friendship?txn=foo&dir=bar
        parsed_url = urlparse(self.path)
        parsed_query = parse_qs(parsed_url.query)
        txn = parsed_query.get("txn")
        directory = parsed_query.get("dir")
        result = self.connect_db("SELECT username, friend FROM friends WHERE username = '%s' AND friend = '%s'" % (txn.pop(), directory.pop()), READ)
        if len(result) == 0:
            self.send_response(400)
            self.end_headers()
        else: 
            self.send_response(200)
            self.end_headers()
    
    def get_directories(self):
        #http://127.0.0.1/directories?username=foo
        parsed_url = urlparse(self.path)
        parsed_query = parse_qs(parsed_url.query)
        username = parsed_query.get("username")
        result = self.connect_db("SELECT directory FROM directories WHERE username='%s'" % username.pop(), READ)
        if len(result) == 0:
            self.send_response(400)
            self.end_headers()
        else:
            resp = {
                "directories": result,
            }
            self.send_response(200)
            self.send_header('Content-Type', 'application/json')
            self.end_headers()
            self.wfile.write(json.dumps(resp).encode())
    
    def get_friends(self): 
        #http://127.0.0.1?friends?username=foo
        parsed_url = urlparse(self.path)
        parsed_query = parse_qs(parsed_url.query)
        username = parsed_query.get("username")
        result = self.connect_db("SELECT friend FROM friends WHERE username='%s'" % username.pop(), READ)
        if len(result) == 0:
            self.send_response(400)
            self.end_headers()
        else: 
            respo = {
                "friends": result,
            }
            self.send_response(200)
            self.send_header('Content-Type', 'application/json')
            self.end_headers()
            self.wfile.write(json.dumps(respo).encode())


    ''''
        TODO -> La query dovrebbe ritornare più valori quindi
        anche le directory e gli amici
    '''
    def search_buddy(self):
        #http://127.0.0.1/search_friend?username=foo
        parsed_url = urlparse(self.path)
        parsed_query = parse_qs(parsed_url.query)
        username = parsed_query.get("username")
        result = self.connect_db("SELECT username FROM users WHERE username='%s'" % username.pop(), READ)
        if len(result) == 0:
            self.send_response(404)
            self.end_headers()
        else: 
            response = {
                "user": result,
            }
            self.send_response(200)
            self.send_header('Content-Type', 'application/json')
            self.end_headers()
            self.wfile.write(json.dumps(response).encode())

    def write_bad_request(self):
        self.send_response(400)
        self.wfile.write(b"You Wrote a Bad Request Buddy!")
    
    def connect_db(self, sql_statement, opcode):
        sql_results = None
        try:
            connection = psycopg2.connect(
                dbname="messdfs",
                user="postgres",
                password="elia",
                host="localhost"
            )

            cursor_obj = connection.cursor()
            cursor_obj.execute(sql_statement)

            match opcode:
                # SELECT statement
                case 0:
                    sql_results = cursor_obj.fetchall()
                # WRITE statement
                case 1:
                    connection.commit()
            connection.close()
        except (Exception, Error) as error: 
            print('Occured Some Errors ', error)

        print(sql_results)
        return sql_results