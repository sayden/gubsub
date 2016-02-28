# WIP gubsub
A Go publish subscriber server

## How to use

### Using go get
```bash
$ go get github.com/sayden/gubsub
```

### Using Git
```bash
$ git clone https://github.com/sayden/gubsub.git
$ cd gubsub
```

### CLI options

#### Command line options

##### Sets the name of the default topic
```bash
$ gubsub --topic default
#or
$ gubsub -t default
```

##### Sets the server listening port
```bash
$ gubsub --port 8080
#or
$ gubsub -p 8080
```

#### Get a list of topics
```bash
$ gubsub topics
```

#### Get a list of listeners
```bash
$ gubsub listeners
```

#### Dispatch a message
```bash
$ gubsub dispatch -t topic [content]
```

## Contributing
Any contribution is welcome.