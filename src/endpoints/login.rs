
use rocket_contrib::json::Json;

use rocket::{
    State, Response,
    http::Status
};

use crate::{
    DbConn, RedisConn, endpoints::{generate_error, generate_response},
    database::handler::HandlerError
};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct PostLogin {
    pub username: String,
    pub password: String
}

#[get("/", format = "json", data = "<data>")]
pub fn login<'r>(handler: State<DbConn>, handler: State<RedisConn>, data: Json<PostLogin>) -> Response<'r> {

}

#[post("/", rank = 2, format = "json", data = "<data>")]
pub fn create_login<'r>(handler: State<DbConn>, data: Json<PostLogin>) -> Response<'r> {

}
