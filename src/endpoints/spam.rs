
use rocket_contrib::json::Json;

use rocket::{
    State, Response,
    http::Status
};

use crate::{
    DbConn, endpoints::{generate_error, generate_response},
    database::{
        handler::HandlerError,
        structures::Message
    }
};

#[get("/<channel>/<platform>/<user>")]
pub fn get_user_offences<'r>(handler: State<DbConn>, channel: String, platform: String, user: String) -> Response<'r> {
    
}
