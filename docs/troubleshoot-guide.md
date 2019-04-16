# Troubleshooting Guide

The following can cause watchman not to start, crash or restart perptually:
- watchman is not scheduled to a manager.
- a managed service is not not properly labelled.
- a service hostname can't be resolved

To troubleshoot check the output of 
`docker service logs <watchman-svc-name>`.

### watchman not scheduled to a manager
```shell
docker service logs nginx                                                                                                             
nginx.1.05lg8a2y2vfg@insyt-worker-01    | 2019-04-16T10:05:01Z c2513e5a2ea0 /usr/local/bin/watchman[1]: INFO Starting watchman                                                     
nginx.1.05lg8a2y2vfg@insyt-worker-01    | 2019-04-16T10:05:01Z c2513e5a2ea0 /usr/local/bin/watchman[1]: INFO Starting NGINX server                                                 
nginx.1.05lg8a2y2vfg@insyt-worker-01    | 2019-04-16T10:05:01Z c2513e5a2ea0 /usr/local/bin/watchman[1]: INFO Generating NGINX service configs                                      
nginx.1.05lg8a2y2vfg@insyt-worker-01    | 2019-04-16T10:05:01Z c2513e5a2ea0 /usr/local/bin/watchman[1]: INFO Verify NGINX serivce confgs                                           
nginx.1.05lg8a2y2vfg@insyt-worker-01    | 2019-04-16T10:05:01Z c2513e5a2ea0 /usr/local/bin/watchman[1]: ERROR Error response from daemon: This node is not a swarm manager. Worker 
nodes can't be used to view or modify cluster state. Please run this command on a manager node or promote the current node to a manager.                                           
nginx.1.05lg8a2y2vfg@insyt-worker-01    | 2019-04-16T10:05:01Z c2513e5a2ea0 /usr/local/bin/watchman[1]: ERROR Error response from daemon: This node is not a swarm manager. Worker 
nodes can't be used to view or modify cluster state. Please run this command on a manager node or promote the current node to a manager.
```

Fix your service deployment and schedule watchman to a Swarm manager node.

### a service hostname can't be resolved
```shell
docker service logs nginx                                                                                                             
nginx.1.ai9s84y45b0j@insyt-manager-01    | 2019-04-16T10:10:30Z 0eadb867c4ea /usr/local/bin/watchman[1]: INFO Starting watchman                                                    
nginx.1.ai9s84y45b0j@insyt-manager-01    | 2019-04-16T10:10:30Z 0eadb867c4ea /usr/local/bin/watchman[1]: INFO Starting NGINX server                                                
nginx.1.ai9s84y45b0j@insyt-manager-01    | 2019-04-16T10:10:30Z 0eadb867c4ea /usr/local/bin/watchman[1]: INFO Generating NGINX service configs                                     
nginx.1.ai9s84y45b0j@insyt-manager-01    | 2019-04-16T10:10:31Z 0eadb867c4ea /usr/local/bin/watchman[1]: INFO Verify NGINX serivce confgs                                          
nginx.1.ai9s84y45b0j@insyt-manager-01    | 2019-04-16T10:10:31Z 0eadb867c4ea /usr/local/bin/watchman[1]: FATAL nginx: [emerg] host not found in upstream "guestbook:80" in /etc/ngi
nx/conf.d/guestbook.conf:2                                                                                                                                                         
nginx.1.ai9s84y45b0j@insyt-manager-01    | nginx: configuration file /etc/nginx/nginx.conf test failed                                                                             
nginx.1.ai9s84y45b0j@insyt-manager-01    |                                                                                                                                         
nginx.1.pfsjcskeimxn@insyt-manager-01    | 2019-04-16T10:10:46Z 32a04f52ba68 /usr/local/bin/watchman[1]: INFO Starting watchman                                                    
nginx.1.pfsjcskeimxn@insyt-manager-01    | 2019-04-16T10:10:46Z 32a04f52ba68 /usr/local/bin/watchman[1]: INFO Starting NGINX server                                                
nginx.1.pfsjcskeimxn@insyt-manager-01    | 2019-04-16T10:10:46Z 32a04f52ba68 /usr/local/bin/watchman[1]: INFO Generating NGINX service configs                                     
nginx.1.pfsjcskeimxn@insyt-manager-01    | 2019-04-16T10:10:46Z 32a04f52ba68 /usr/local/bin/watchman[1]: INFO Verify NGINX serivce confgs                                          
nginx.1.pfsjcskeimxn@insyt-manager-01    | 2019-04-16T10:10:46Z 32a04f52ba68 /usr/local/bin/watchman[1]: FATAL nginx: [emerg] host not found in upstream "guestbook:80" in /etc/ngi
nx/conf.d/guestbook.conf:2                                                                                                                                                         
nginx.1.pfsjcskeimxn@insyt-manager-01    | nginx: configuration file /etc/nginx/nginx.conf test failed
```

Fix this by creating an overlay network and attaching watchman and all services it manages to the network.

### Improperly labelled service
Improper or missing labels can cause NGINX config verification to fail and watchman to restart endlessly.
Make sure the service has set all watchman labels are set correctly :
- watchman.service.enable
- watchman.service.port
- watchman.service.url

### My service is not still not published
Check that your service has the `watchman.service.enable` label set to true.

As a last resort you can inspect the NGINX configs generated by watchman in the location specified
by --output-dir (defaults to /etc/nginx/conf.d).
