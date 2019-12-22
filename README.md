# NodeWrecker


# Run via docker
`docker start jaeg/nodewrecker --threads=4 --escalate=true --abuse-memory=true`

# Flags
- threads 
    - default:4
    - Number of threads to run
- sleep 
    - default:1
    - milliseconds to sleep
- escalate 
    - default:false
    - Keep creating threads
- escalate-rate 
    - default:1000
    - milliseconds between creating new threads
- string-length 
    - default:1000
    - length of randomly generated string
- abuse-memory
    - defailt:false
    - if true nodewrecker will store all generated values in memory
