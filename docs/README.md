Orbital 2023 Milestone 1
 


Team Name:  
The Card Wizards
 
Team ID: 
5451
 
Team Members: 
Wang Zihan, 
Wang Yiwen
 
Proposed Level of Achievement: 
Apollo 11

Links & Refereces
Link to Google doc: 
https://docs.google.com/document/d/1ozqRrvVdwpz_M76d3R9yiwmT_47Dc-hkO5nU0ls2Oj8/edit#
Link to repository: 
https://github.com/yiwen101/TiktokXOrbital-CardWizards
Link to doc (work in progress): 
https://yiwen101.github.io/TiktokXOrbital-CardWizards/
Link to Project log:
https://nusu-my.sharepoint.com/:x:/g/personal/e1075515_u_nus_edu/ERV1K2okJktLgZuz_TMjbVEBSY2KWaOg7e1FZSsck6jyeA?email=wang.yiwen.e0968880%40u.nus.edu&e=5hbEr2

Milestone 1 
Pitch Poster:
Todo
 
Pitch Video:
Todo



Project information

Project Motivation
Our team decided to take up the external proposal from tiktok. The following is a summary of important information from their proposal
According to tiktok’s proposal, “this project aims to equip students with essential knowledge and practical experience in HTTP, JSON, Thrift, Load Balancing, Service Register and Discovery, as well as building HTTP servers using Hertz and RPC servers using Kitex.”
Aim of Project
“The primary goal of this project is to implement an API Gateway that accepts HTTP requests encoded in JSON format and uses the Generic-Call feature of Kitex to translate these requests into Thrift binary format requests. The API Gateway will then forward the request to one of the backend RPC servers discovered from the registry center. “
Requirements and Deliverables
“Students participating in this project must fulfill the following requirements:
Implement an API Gateway that accepts HTTP requests with JSON-encoded bodies.
Use the Generic-Call feature in Kitex to translate JSON requests into Thrift binary format.
Integrate a load balancing mechanism to distribute requests among backend RPC servers.
Integrate a service registry and discovery mechanism for the API Gateway and RPC servers.
Develop backend RPC servers using Kitex for testing the API Gateway.
Document the project, including design decisions, implementation details, and usage instructions. “
Tech Stack:
Stipulated Tech Stack:
Golang: 
HTTP: 
JSON:
Thrift: Both as an interface definition language (IDL) and binary communication protocol.
Kitex framework for RPC, generic call, service registry and discovery and load balancing
Hertz framework to build http servers
Potential Tech Stack 
 Relational database (SQL) and Gorm
 Docker and crul
 Postman for testing and monitoring
 Other authentication and validation tech, or existing solution that support caching and other plugins



User Cases and Features
    
User stories:
Core user stories:
As a user, I want to be able to send HTTP requests to the API Gateway.
As a user, I want the API Gateway to translate my JSON-encoded HTTP requests into Thrift binary format.
As a user, I want the API Gateway to forward my request to one of the backend RPC servers discovered from the registry center.
As a user, I want the API Gateway to distribute requests among backend RPC servers using a load balancing mechanism.
As a user, I want to be able to integrate a service registry and discovery mechanism for the API Gateway and RPC servers.
As a user, I wish to be able to exposed to useful set of admin api to configure and deploy the gateway for my microservices
As a user, I wish to see documentation and tutorials for the API to enable me to quickly get started and make reference and fix simple problems
Extensions:
As a user, I wish to see sophisticated documentation and tutorials for the API 
As a user, I wish to be offered code examples on how to use different functions of the gateway
As a user, I wish to be offered business demons on what I could build with the gateway
As a user, I wish to have Graphical user interface for the ease of controlling and monitoring
As a user, I wish to be offered a plugins that monitor and trace how a request is handled
As a user, I wish to be offered plugins that boost security and defend against attacks
As a user, I wish to be offered plugins that support response caching to improve performance
As a user, I wish to be offered plugins that help with user authentication, certificate validations
As a user, I wish to be offered help that allow me to build my plugins for my customized needs






Features:
Basic Features (Simple proxy and routing):
Accept HTTP requests with JSON-encoded bodies with a server.
Forward requests to backend RPC servers
Offer admin API to control and configure the api gateway
Develop backend RPC servers using Kitex for testing the API Gateway.
Documentation
Intermediate Features
Implement support for multiple programming languages (e.g., Go, Java) by Translate JSON requests into Thrift binary format using the Generic-Call feature in Kitex.
Distribute requests among backend RPC servers using a load balancing mechanism.
 Integrate a service registry and discovery mechanism for the API Gateway and RPC servers.
Advanced Features
Implement proper error handling for HTTP requests.
Implement rate limiting to prevent abuse of the API Gateway.
Implement authentication and authorization mechanisms to secure the API Gateway.
Implement caching mechanisms to improve performance.
Implement a user interface (like kong manager) to facilitate the management of api gateway
Implement support for multiple protocols (e.g., gRPC, REST).
 Implement support for multiple data formats (e.g., XML, Protocol Buffers).
 Implement support for distributed tracing to help with debugging and performance optimization.



Technical documentation

Architecture and flow:
Overall Architecture:
Flow Chart:
When a Http request is catched by the listener, it will first decide whether it is from a client or from the Admin component. Then, it will ask the router to find the matched route and pick the most matched one based on priority policies. If there are not any matched ones, the Route component will respond with a default response. Otherwise, if the request is from the Admin component, it will call the configurer in the stipulated component. Otherwise, it will forward the request to the Service component with the route found. The service component will find the registered upstreams and then make a generic call to it, where upstream manager will eventually dispatch the call to a microservice server at the backend according the load balancing strategy set.
Route component
Service component

Upstream Manager component



Plan for Implementation

SE

Testing
