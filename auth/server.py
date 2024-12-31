# 
# MessDFS Authentication Microservice
#

import app.router as router


def main():
    r = router.Router("127.0.0.1", 8082)
    r.start_handler()

main()

