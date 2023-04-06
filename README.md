# tracing-example

jaeger

https://www.jaegertracing.io/docs/1.6/architecture/

docker pull jaegertracing/all-in-one

docker image ls

docker run -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 -p 5775:5775/udp -p 6831:6831/udp -p 6832:6832/udp -p 5778:5778 -p 16686:16686 -p 14268:14268 -p 9411:9411 -d jaegertracing/all-in-one

=========================================================
Document:

https://hub.docker.com/r/jaegertracing/all-in-one

https://www.jaegertracing.io/docs/1.6/getting-started/

https://www.jaegertracing.io/docs/1.19/opentelemetry/

======
Install docker

Visit: https://docs.docker.com/desktop/install/mac-install/
Click "Docker Desktop for Mac with Apple silicon" to download .dmg

Pull docker image

docker pull jaegertracing/all-in-one

Run docker

docker run -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 -p 5775:5775/udp -p 6831:6831/udp -p 6832:6832/udp -p 5778:5778 -p 16686:16686 -p 14268:14268 -p 9411:9411 -d jaegertracing/all-in-one

UI
Visit: http://localhost:16686/search

// 
https://jessitron.com/2021/08/11/run-an-opentelemetry-collector-locally-in-docker/

https://uptrace.dev/get/opentelemetry-gin-gorm.html#instrumenting-gin

https://docs.docker.com/build/building/opentelemetry/

