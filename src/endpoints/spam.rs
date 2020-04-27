
use rocket_contrib::json::Json;

use rocket::{
    State, Response,
    http::Status
};

use crate::{
    DbConn, endpoints::{generate_error, generate_response},
    database::{
        handler::HandlerError,
        structures::UpdateCount
    }
};

#[get("/<channel>/<service>/<user>")]
pub fn get_user_offences<'r>(handler: State<DbConn>, channel: String, service: String, user: String) -> Response<'r> {
    let offences = handler.lock().expect("db lock").get_offences(&channel, &service, &user);
    match offences {
        Ok(offences) => generate_response(Status::Ok, json!({ "data": offences })),
        Err(HandlerError::Error(_)) => generate_response(Status::NotFound, json!({})),
        Err(e) => {
            println!("Internal error getting offenses: {:?}", e);
            generate_response(Status::InternalServerError, generate_error(500, None))
        }
    }
}

#[get("/<channel>/<service>/<user>/<key>")]
pub fn get_user_offence<'r>(handler: State<DbConn>, channel: String, service: String, user: String, key: String) -> Response<'r> {
    let offences = handler.lock().expect("db lock").get_offence(&channel, &service, &user, &key);
    match offences {
        Ok(offences) => generate_response(Status::Ok, json!({ "data": offences })),
        Err(HandlerError::Error(e)) => generate_response(Status::NotFound, json!({})),
        Err(e) => {
            println!("Internal error getting offenses: {:?}", e);
            generate_response(Status::InternalServerError, generate_error(500, None))
        }
    }
}

#[patch("/<channel>/<service>/<user>/<key>", format = "json", data = "<update>")]
pub fn update_user_offences<'r>(handler: State<DbConn>, channel: String, service: String, user: String, key: String, update: Json<UpdateCount>) -> Response<'r> {
    let offences = handler.lock().expect("db lock").update_offence(&channel, &user, &service, &key, update.into_inner());
    match offences {
        Ok(offences) => generate_response(Status::Ok, json!({ "data": offences })),
        Err(HandlerError::Error(e)) => generate_response(Status::NotFound, json!({})),
        Err(e) => {
            println!("Internal error getting command: {:?}", e);
            generate_response(Status::InternalServerError, generate_error(500, None))
        }
    }    
}
