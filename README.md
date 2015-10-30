# Todo API

## Building

    $ cd $GOPATH/github.com/bradylove/todo
    $ ./build.sh

## Todo Endpoints

### Index

    GET /todos

### Show

    GET /todos/:id

### Create

    POST /todos

Accepted parameters are: `description`. Example:

    {
       "description": "Write an awesome todo API"
    }

### Updated

    PUT /todos/:id

Accepted parameters are: `description`, `status`. Example:

    {
       "description": "Write an even more awesome todo API",
       "status": "completed"
    }

### Todo#Delete

    DELETE /todos/:id
