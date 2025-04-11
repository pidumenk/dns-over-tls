## DNS OVER TLS

### Prerequisites
- docker
- golang 

### Configuration

How to run Go application locally in Docker:

```bash
# Set environmant vatiables to make bootstrapping handy
export N26_IMAGE_NAME=godns
export N26_CONTAINER_NAME=godns

# Build docker image
docker build -t $N26_IMAGE_NAME .

# Run dns-proxy container DNS over TLS by default
docker run --rm -d -p 53:53 --name $N26_CONTAINER_NAME $N26_IMAGE_NAME

# Run dns-proxy on standart UDP/53 port
docker run --rm -d -p 53:53/udp --name $N26_CONTAINER_NAME $N26_IMAGE_NAME udp

# Testing DNS over TLS requests by CLI
dig +short +tcp n26.com @localhost

# Testing standart DNS UDP requests by CLI
dig +short n26.com @localhost

# Check container logs (if necessary)
docker logs $N26_CONTAINER_NAME

# Stop dns-proxy service
docker stop $N26_CONTAINER_NAME
```
### FAQ

1) Imagine this proxy being deployed in an infrastructure. What would be the security
concerns you would raise?
    
    * Depends on many factors how it was deployed: on-premise, AWS, kubernetes and etc. 
    
    * Generally make sure all reliable and modern TLS ciphers and portocol are used. 

    * Logging and Monitoring: make sure the application captures DNS requests, responses or any potential errors and send them to a monitoring system (ELK, Datadog, Prometheus).

    * Authentication and autorization: in nutshell, keep principle of least privilages. Allow communications with DNS proxy only that services which really require this. That can be achieved by many different ways: firewall, NACL, security groups, kuberentes network policies and so on.

    * DDOS Protection: for instance, rate limiter to control the number of possible requests.

    * Application Security Testing: Make sure the app doesn't contain any modules or dependencies recently detected in CVE. 

    * Application and deployment redundancy.

    * DNS Cache Poisoning: https://www.cloudflare.com/learning/dns/dns-over-tls/

2) How would you integrate that solution in a distributed, microservices-oriented and
containerized architecture?

    * Use container orchestration solutions for microservices: Kubernetes, EKS, OpenShift, GKE and so on. 

    * Follow best practicies and IaC approach: keep config files as a code and wrap into variables. (Helm Charts, Terrafrom, OpenTofu).

    * In case of multiple instances of the DNS proxy service, use a load balancer to distribute traffic among them. By default Kubernetes provides out og the box functionality for service load balancing. Also can be the MetalLB for bare metal infrastructure. 

    * Levarage centralized logging solutions to capture DNS logs coming from all DNS-proxy instances: (ELK, prometheus, Datadog, Grafana).

    * Scalability and redundancy: design the DNS proxy service to scale horizontally by deploying multiple instances to handle increased load. Also trying to keep the application logic redundent (for example, allow multiple incoming requests). 

    * Follow best practices of SDLC. Especially, implement CI/CD pipelines to for your service to handle simply deployments, tests and reverting changes (GitHub Actions, ArgoCD, Jenkins).

    * Integrate application with secure 3d party application to store sensetive information in the code safely. For example, use HashiCorp Vault to keep secrets and automatically replace variables during deployment.

    * Consider versioning the DNS proxy service to enable seamless updates and rolling deployments.

    * Prepare and provide detailed techincal documentations about the service. 

3) What other improvements do you think would be interesting to add to the project?

    * Define k8s manifest files to bootstrap the application in Kuberentes (or Helm Charts). 
    
    * Caching mechanism for faster DNS response.

    * Implement black/white lists or dns filtering capabilities. 

    * Health Checks.

    * Feauture flags functionality to enable or disable a feature without modifying the source code or redeploying.

    * Custom logging format: it can be efficient in some cases. 
