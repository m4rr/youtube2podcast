docker pull mxpv/podsync:latest
docker run \
    -p 8080:8080 \
    -v $(pwd)/data:/app/data/ \
    -v $(pwd)/config.toml:/app/config.toml \
    mxpv/podsync:latest
