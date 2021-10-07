rm ./docker/firebase.json
rm ./docker/blog-microservice.tar

export $(grep ^GOOGLE_APPLICATION_CREDENTIALS .env)
echo $GOOGLE_APPLICATION_CREDENTIALS

cp $GOOGLE_APPLICATION_CREDENTIALS ./docker/firebase.json

./compile-to-linux-release.sh

(
  cd docker && \
  docker build -t blog-microservice . && \
  docker save blog-microservice -o blog-microservice.tar
)