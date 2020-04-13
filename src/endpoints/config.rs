
use rocket_contrib::json::Json;
use serde_json::Value;

use rocket::{
    State, Response,
    http::Status
};

use crate::{
    DbConn, endpoints::{generate_error, generate_response},
    database::handler::HandlerError
};

#[get("/<channel>")]
pub fn get_config<'r>(handler: State<DbConn>, channel: String) -> Response<'r> {
    let config = handler.lock().expect("db lock").get_config(&channel);

    match config {
        Ok(config) => generate_response(Status::Ok, json!({ "data": config })),
        Err(HandlerError::Error(_)) => generate_response(Status::NotFound, generate_error(404, None)),
        Err(e) => {
            println!("Internal error getting channel config: {:?}", e);
            generate_response(Status::InternalServerError, generate_error(500, None))
        }
    }
}

#[patch("/<channel>", format = "json", data = "<config>")]
pub fn update_config<'r>(handler: State<DbConn>, channel: String, config: Json<Value>) -> Response<'r> {
    println!("{:?}", config);
    generate_response(Status::Ok, json!({ }))
    // let result = handler.lock().expect("db lock").update_config(channel);
    // match result {
    //     Ok(()) => generate_response(Status::Ok, json!({ "meta": json!({ "updated": true }) })),
    //     Err(HandlerError::Error(_)) => generate_response(Status::NotFound, generate_error(404, None)),
    //     Err(e) => {
    //         println!("Internal error updating config: {:?}", e);
    //         generate_response(Status::InternalServerError, generate_error(500, None))
    //     }
    // }
}
