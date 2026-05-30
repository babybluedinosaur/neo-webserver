Go webserver that exposes near earth asteroids by week using NASA API

Pull image for webserver from dockerhub:

```bash
docker pull sidondocker/nasa-webserver:latest
```

Run container (choose API key):

```bash
docker run -p 8080:8080 -e NASA_API_KEY=DEMO_KEY sidondocker/nasa-webserver:latest
```

Run health check:

```bash
curl http://localhost:8080/health
```

Get ids of 4th calendar week:

```bash
curl http://localhost:8080/neo/week/4
```

Get hazardous asteroids of 4th calendar week:

```bash
curl "http://localhost:8080/neo/week/4?filter=hazardous"
```

