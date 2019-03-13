

use rocket::{
	State, Response,
	http::Status
};

use crate::{
	DbConn, endpoints::{generate_error, generate_response},
	database::{
		handler::HandlerError
	}
};

#[get("/<channel>")]
pub fn get_aliases<'r>(handler: State<DbConn>, channel: String) -> Response<'r> {
	let aliases = handler.lock().expect("db lock").get_aliases(&channel);
	match aliases {
		Ok(aliases) => generate_response(Status::Ok, json!({ "data": aliases })),
		Err(HandlerError::Error(e)) => generate_response(Status::NotFound, generate_error(404, Some(e))),
		Err(e) => {
			println!("Internal error getting aliases: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}

#[get("/<channel>/<command>")]
pub fn get_alias<'r>(handler: State<DbConn>, channel: String, command: String) -> Response<'r> {
	let alias = handler.lock().expect("db lock").get_alias(&channel, &command);
	match alias {
		Ok(alias) => generate_response(Status::Ok, json!({ "data": alias })),
		Err(HandlerError::Error(e)) => generate_response(Status::NotFound, generate_error(404, Some(e))),
		Err(e) => {
			println!("Internal error getting alias: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}

#[post("/<channel>/<command>/<alias>")]
pub fn create_alias<'r>(handler: State<DbConn>, channel: String, command: String, alias: String) -> Response<'r> {
	let result = handler.lock().expect("db lock").create_alias(&channel, &command, &alias);
	match result {
		Ok(alias) => generate_response(Status::Ok, json!({ "data": alias })),
		Err(HandlerError::Error(e)) => generate_response(Status::NotFound, generate_error(404, Some(e))),
		Err(e) => {
			println!("Internal error creating alias: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}

#[delete("/<channel>/<command>")]
pub fn delete_alias<'r>(handler: State<DbConn>, channel: String, command: String) -> Response<'r> {
	let result = handler.lock().expect("db lock").delete_alias(&channel, &command);
	match result {
		Ok(()) => generate_response(Status::Ok, json!({ "data": json!({ "deleted": true }) })),
		Err(HandlerError::Error(e)) => generate_response(Status::NotFound, generate_error(404, Some(e))),
		Err(e) => {
			println!("Internal error creating trust: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}
