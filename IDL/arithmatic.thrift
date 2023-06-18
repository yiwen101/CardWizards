namespace go arithmatic

struct Request {
    1: i64 firstArguement (api.query = 'firstArguement')
    2: i64 SecondArguement 
    3: optional string message 
    4: optional map<string, string> Extra
}

struct Response {
    1: i64 firstArguement (api.body = 'firstArguement')
    2: i64 SecondArguement (api.body = 'secondArguement')
    3: optional string message
    4: i64 result (api.body = 'result')
}

service Calculator {
    Response Add(1: Request request ) ( api.get = "/arith/add" )
    Response Subtract(1: Request request) ( api.get = "/arith/subtract")
    Response Multiply(1: Request request) ( api.head = "/arith/multiply")
    Response Divide(1: Request request)
}
