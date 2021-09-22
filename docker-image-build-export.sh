./compile-to-linux-release.sh

(
  cd docker && \
  docker build -t blog-microservice . && \
  docker save blog-microservice -o blog-microservice.tar
)