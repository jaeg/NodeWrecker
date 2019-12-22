# NodeWrecker
Stress test your cluster under sporadic high cpu and high memory load. 

# Run via docker
`docker start jaeg/nodewrecker --threads=4 --escalate=true --abuse-memory=true`

# Install via helm
`helm upgrade --install pi-wrecker ./helm-chart/`

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
    - default:false
    - if true nodewrecker will store all generated values in memory
- min-duration
    - default: 10
    - minimum seconds a test lasts
- max-duration
    - default: 60
    - max seconds a test lasts
- max-deay
    - default: 10
    - max seconds between tests
- min-delay
    - defaults: 10
    - min seconds between tests
- verbose
    - defaults: false
    - output everything from threads