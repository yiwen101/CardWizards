namespace go arithmatic

struct Request {
    1: i64 firstArguement
    2: i64 SecondArguement
    3: optional string message
}

struct Response {
    1: i64 firstArguement
    2: i64 SecondArguement
    3: optional string message
    4: i64 result
}

service Calculator {
    Response Add(1: Request request)
    Response Subtract(1: Request request)
    Response Multiply(1: Request request)
    Response Divide(1: Request request)
}
