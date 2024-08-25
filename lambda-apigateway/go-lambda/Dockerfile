FROM golang:1.22 AS build
WORKDIR /lambda
COPY go.mod go.sum ./
COPY main.go db.go  ./
RUN go build -tags lambda.norpc -o main .

FROM public.ecr.aws/lambda/provided:al2023
COPY --from=build /lambda/main ./main
ENTRYPOINT [ "./main" ]