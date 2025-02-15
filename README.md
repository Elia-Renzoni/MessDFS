# MessDFS
MessDFS is a distributed file system built in a microservices environment that allows users to manage a remote space to work with their CSV files.

## Project Features
* Possibility to work with remote files and directories.
* Possibility to read friends' files.

## Microservices Communication
There are two microservices: the first one controls the remote resources of each user, such as directories and files, while the second one stores all the information on a PostgreSQL instance. This information is essential to ensure exclusive access for users to their own directories and files. Therefore, the communication between the two microservices revolves around this aspect, where the storage microservice communicates with the auth microservice when it needs to verify if the user of the transaction is the owner of the indicated directories or if a user can read another user's file based on friendship relations. Additionally, the storage microservice also communicates with auth when a user decides to delete a directory, removing both the logical and physical directory.