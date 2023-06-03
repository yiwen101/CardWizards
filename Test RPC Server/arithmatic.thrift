namespace go api

struct Request {
    1: i32 firstArguement = 0
    2: i32 secondArguement = 0
    3: optional string message
}

struct Response {
    1: i32 result = 0
    2: optional string message
}

service Arithmatic {
    Response Add(1: Request req)
    Response Sub(1: Request req)
    Response Mul(1: Request req)
    Response Div(1: Request req)
}
