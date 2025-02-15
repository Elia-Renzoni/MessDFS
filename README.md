# MessDFS
MessDFS is a distributed file system built in a microservices environment that allows users to manage a remote space to work with their CSV files.

## Project Features
* Possibility to work with remote files and directories.
* Possibility to read friends' files.

## Microservices Communication
There are two microservices: the first one controls the remote resources of each user, such as directories and files, while the second one stores all the information on a PostgreSQL instance. This information is essential to ensure exclusive access for users to their own directories and files. Therefore, the communication between the two microservices revolves around this aspect, where the storage microservice communicates with the auth microservice when it needs to verify if the user of the transaction is the owner of the indicated directories or if a user can read another user's file based on friendship relations. Additionally, the storage microservice also communicates with auth when a user decides to delete a directory, removing both the logical and physical directory.

## API
### Storage
* <b>insert data into a csv file</b>
```json
{
    "txn_user": "<transaction user>",
    "query_type": "insert",
    "user": "<directory name>",
    "file_name": "<file name>",
    "query_content": ["<values>"]
}
```
* <b>directory creation</b>
```json
{
    "txn_user": "<transaction user>",
    "dir_to_create": "<directory>"
}
```
* <b>update csv file</b>
```json
{
    "txn_user": "<transaction user>",
    "query_type": "update",
    "user_name": "<directory name>",
    "file_name": "<file name>",
    "query_content": {
        "<column name>": ["id", "old data value", "new data value"]
    }
}
```
* <b>read data</b>
```
http://<IP address>:8081/csvr/{txn user}/{friend name}/{directory}/{file name}?id=<id value>
```
* <b>delete directory</b>
```
http://<IP address>:8081/ddir/{txn user}/{directory to delete}
```
* <b>delete file</b>
```
http://<IP address>:8081/dfile/{txn user}/{file to delete}/{directory}
```
* <b>delete data</b>
```
http://<IP address>:8081/csvd/{txn user}/{directory}/{file name}?id=<id value>
```

