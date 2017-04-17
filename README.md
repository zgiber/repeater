# repeater
A primitive http server which for a GET request returns the body of the latest POST request. 
Request body is kept in memory, so it will be lost on restart. It's mostly for small request bodies.
Use it for whatever you'd like though. I'm using it to expose public-safe data from a private network.

One way of running it with docker:

1. Build the image
`docker build -t repeater .`

2. Run the server
`docker run -d --name repeater -p80:8000 --restart on-failure repeater`

Use another port by changing -p value if you like. For example, running it on port 8765 requires to set: `-p8765:8000`
