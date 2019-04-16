# Configuration Guide
There are three components of watchman which are configuration:
- watchman
- Swarm services
- NGINX

## watchman configuration
```
  -conf-dir string       
        watchman conf directory (default "/etc/watchman")                                           
  -config-file string    
        the watchman config file (default "/etc/watchman/watchman.yaml")                            
  -interval int          
        backend polling interval (default 60)     
  -log-level string      
        level which watchman should log messages  
  -output-dir string     
        the NGINX configuration directory (default "/etc/nginx/conf.d")                             
  -service-template-file string                   
        the NGINX service template file (default "/etc/watchman/nginx/service.tmpl")                
```

## Swarm service configuration
watchman manages services that have the following service labels set:
- `watchman.service.enable` must be set to true, watchman will ignore any service that
  does not have this label set or that is set to any value other than true.
- `watchman.service.url` must be set. The url value must not include a protocol component
  for example, example.com is a valid value, while https://example.com is an invalid value.
- `watchman.service.port` must be set the port on which the service is listening.

## NGINX configuration
Both nginx.conf and service template files can be replaced, the only requirement
is that the nginx.conf file must include additional configs in a specified
directory like /etc/nginx/conf.d.

For example the default nginx.conf and service template file that comes with the
watchman Docker image do not support HTTPS and we must replace both to enable HTTPS
support for watchman managed services.

The default nginx.conf can be replaced with the following configuration assuming we
have already created our docker secrets for ssl.key and ssl.crt for NGINX SSL
configuration. 

Read on how to generate SSL certificates and create docker secrets
[here](https://docs.docker.com/engine/swarm/secrets/#intermediate-example-use-secrets-with-a-nginx-service).

```
cat > newnginx.conf <<EOF
user  nginx;
worker_processes  1;

error_log  /var/log/nginx/error.log warn;
pid        /var/run/nginx.pid;

events {
    worker_connections  1024;
}

http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';

    access_log  /var/log/nginx/access.log  main;

    sendfile        on;
    keepalive_timeout  65;

    server_tokens off;

    # SSL config
    ssl_prefer_server_ciphers on;
    # Use only TLS
    ssl_protocols TLSv1 TLSv1.1 TLSv1.2;
    ssl_ciphers ECDH+AESGCM:ECDH+AES256:ECDH+AES128:DH+3DES:!ADH:!AECDH:!MD5;
    # Enable HSTS
    add_header Strict-Transport-Security "max-age=31536000" always;
    # Optimize session cache
    ssl_session_cache   shared:SSL:40m;
    ssl_session_timeout 4h;
    # Enable session tickets
    ssl_session_tickets on;
    ssl_certificate     /run/secrets/ssl.crt;
    ssl_certificate_key /run/secrets/ssl.key;
    
    include /etc/nginx/conf.d/*.conf;
}
EOF
```

and the default service template file with:
```
cat > newsvc.tmpl <<EOF
upstream {{ .ServiceName }} {
    server {{ .ServiceName }}:{{ .ServicePort }};
}

server {
    listen      80;
    server_name {{ .ServiceURL }};
    return      301  https://{{ .ServiceURL }}/;   
}

server {
    listen      443 ssl http2;
    server_name {{ .ServiceURL }};

    location / {
        proxy_pass        http://{{ .ServiceName }}/;
        proxy_redirect    off;
        proxy_set_header  Host      $host;
        proxy_set_header  X-Real-IP $remote_addr;
        proxy_hide_header X-POWERED-BY;
    }
}
EOF
```

Then create docker configs for nginx.conf and service template file.
```
docker create config watchman-nginx.conf newnginx.conf
docker create config watchman-nginx-service.tmpl newsvc.tmpl
```

Now start the watchman service passing the watchman-nginx.conf and 
watchman-nginx-service.tmpl configs which will be mounted at /etc/nginx/nginx.conf
and /etc/watchman/nginx/service.tmpl.

The ssl.key and ssl.crt secrets are expected at /run/secrets by the new nginx.conf.

```
docker service create --name nginx \
  --mount type=bind,source=/var/run/docker.sock,destination=/var/run/docker.sock \
  --constraint 'node.role==manager' \
  --network shared-net --publish 80:80 --publish 443:443 \
  --config source=watchman-nginx.conf,target=/etc/nginx/nginx.conf \
  --config source=watchman-nginx-service.tmpl,target=/etc/watchman/nginx/service.tmpl \
  --secret ssl.key --secret ssl.crt \
  ttousai/watchman:v0.1.0-alpha.1

docker service ls
ID                  NAME                  MODE                REPLICAS            IMAGE                             PORTS
udunlls8do5p        nginx                 replicated          1/1                 ttousai/watchman:v0.1.0-alpha.1   *:80->80/tcp, *:443->443/tcp
```
