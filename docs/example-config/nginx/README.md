## NGINX config
`nginx.conf` is an example nginx configuration file that uses the [ngx_http_auth_request](http://nginx.org/en/docs/http/ngx_http_auth_request_module.html) module which allows NGINX to consult `auth-baton` before accepting a request.

Assuming `auth-baton` is reachable @ `http://localhost:6800/`, you can run nginx with the given configuration and send a few simple requests to see if the flow works:

Case (1): auth-baton rejects the request
```
curl http://localhost:8090 -H "Authorization: Basic xyz"
<html>
<head><title>403 Forbidden</title></head>
<body>
<center><h1>403 Forbidden</h1></center>
<hr><center>nginx/1.19.5</center>
</body>
</html>

```

Case (2) auth-baton authenticates the request
```
curl http://localhost:8090/edit/profile -H "Authorization: Basic dXNlcjpwYXNz"
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>
```

**Note:** Since the original request path is reused, for this last request, the 
underlying request to authbaton will look like:
```
curl http://localhost:6800/edit/profile -H "Authorization: Basic dXNlcjpwYXNz"
``` 

