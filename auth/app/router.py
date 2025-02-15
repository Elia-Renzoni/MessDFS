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
            case "/ownership": 
                self.check_ownership()
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
        parsed = urlparse(self.path)
        match parsed.path:
            case "/delete-dir":
                self.delete_directory()
            case "/delete-friend":
                self.delete_friend()
            case _:
                self.write_bad_request()
    
    
    '''
        This method is responsible for performing a login.
    '''
    def login(self):
        print("login", file=sys.stdout)
        parsed_url = urlparse(self.path)
        parsed_query = parse_qs(parsed_url.query)
        username = parsed_query.get("username")
        password = parsed_query.get("password")
        result = self.connect_db("SELECT username, password FROM users WHERE username = '%s' AND password = '%s'" % (username.pop(), password.pop()), READ)
        
        # if the response from the database server is empty
        if len(result) == 0:
            self.send_response(400)
        else:
            self.send_response(200)
        self.end_headers()
    
    '''
        This method adds a user to the system database.
    '''
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
    

    '''
        This method add a friend relationships between users
    '''
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


    
    '''
        This method is responsible for checking the friendship beetwen two users
        This method is very useful in the case of read operation. Indeed this method
        is called by the storage microservice when a client perform a read request
        to a csv file.
    '''
    def check_friendship(self):
        parsed_url = urlparse(self.path)
        parsed_query = parse_qs(parsed_url.query)
        txn = parsed_query.get("txn")
        friend = parsed_query.get("friend")

        check_friend = self.connect_db("SELECT username, friend FROM friends WHERE username = '%s' AND friend = '%s'" % (txn.pop(), friend.pop()), READ)
        if len(check_friend) == 0:
            self.send_response(400)
            self.end_headers()
        else:
            self.send_response(200)
            self.end_headers()


    '''
        This method is responsible for checking the ownership relation between a given
        user and directory. The following method is called by the storage microservice
        when performing every http request
    '''
    def check_ownership(self): 
        print("i'm checking the ownership", file=sys.stdout)    
        #http://127.0.0.1/friendship?txn=foo&dir=bar
        parsed_url = urlparse(self.path)
        parsed_query = parse_qs(parsed_url.query)
        txn = parsed_query.get("txn")
        directory = parsed_query.get("dir")
        
        check_directory_result = self.connect_db("SELECT username FROM directories WHERE directory = '%s'" % (directory.pop()), READ)
        owner_name = check_directory_result[0]
        print(owner_name[0])
        if owner_name[0] != txn.pop():
            self.send_response(400)
            self.end_headers()
        else: 
            self.send_response(200)
            self.end_headers()
    
    '''
        This method returns a list of directories owned by the given user
    '''
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
    
    def delete_directory(self):
        parsed_url = urlparse(self.path)
        parsed_query = parse_qs(parsed_url.query)
        directory_to_delete = parsed_query.get("directory")


        reuslt = self.connect_db("DELETE FROM directories WHERE directory = '%s'" % directory_to_delete.pop(), TO_COMMIT)
        if reuslt == None:
            self.send_response(200)
            self.end_headers()
        else:
            self.send_response(400)
            self.end_headers()

    def delete_friend(self):
        parsed_url = urlparse(self.path)
        parsed_query = parse_qs(parsed_url.query)
        friend_to_delete = parsed_query.get("friend")

        reuslt = self.connect_db("DELETE FROM friends WHERE friend = '%s'" % friend_to_delete.pop(), TO_COMMIT)
        if reuslt == None:
            self.send_response(200)
            self.end_headers()
        else:
            self.send_response(400)
            self.end_headers()


    def write_bad_request(self):
        badResponse = {
            "bad": "You Wrote a Bad Request Buddy!",
        }

        self.send_response(400)
        self.send_header('Content-Type', 'application/json')
        self.end_headers()
        self.wfile.write(json.dumps(badResponse).encode())
    
    '''
        This method is responsible for connect the server to the database
        the opcode paramter is very useful as is responsible to notify
        the aim of the SQL statement attached as parameter
    '''
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