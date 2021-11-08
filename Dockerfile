FROM k8deployment-base as build
WORKDIR /app
COPY CMakeLists.txt CMakeLists.txt
COPY conanfile.txt conanfile.txt
COPY go go
COPY main.c main.c
RUN cd go && go test && go build -buildmode=c-shared -ldflags="-w -s" -gcflags=all=-l -gcflags=all=-B -o ../libs/ready/libready.so libready.go
RUN mkdir conan && cd conan && conan install .. --build=missing
RUN mkdir build && cd build && cmake .. && cmake --build . && cd bin && staticx ready ready-static && rm -rf /tmp/*
#RUN cd build && ldd k8deployment | tr -s '[:blank:]' '\n' | grep '^/' | xargs -I % sh -c 'mkdir -p $(dirname deps%); cp % deps%;'

FROM scratch
#FROM gcr.io/distroless/static
#gcr.io/distroless/base
#gcr.io/distroless/static
#COPY --from=build /app/build/deps /
COPY --from=build /app/build/bin/ready-static  /ready
COPY --from=build /tmp /tmp
CMD ["./ready"]