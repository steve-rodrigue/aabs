## Docker

### Build all services

Build all Docker images from scratch without using the build cache.

docker compose build --no-cache 

### Start all services

Start all services in detached mode (background).

docker compose up -d 

### Rebuild and restart

After modifying source code, dependencies, or Dockerfiles:

docker compose build --no-cache docker compose up -d 

### View logs

docker compose logs -f 

### View container resource usage

Display CPU, memory, network, and disk usage for running containers.

docker stats 

Display resource usage for a specific container.

docker stats aabs-embeddings 

### Stop all services

docker compose down 

### Stop all services and remove volumes

docker compose down -v 