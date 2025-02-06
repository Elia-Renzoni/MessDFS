# 
# MessDFS Authentication Microservice
#

from app.router import Router
from http.server import HTTPServer

def main(server_class=HTTPServer, handler_class=Router, port=8083):
    server_addr = ('127.0.0.1', 8083)
    httpd = server_class(server_addr, handler_class)
    print("Server Listening...")
    httpd.serve_forever()

main()

