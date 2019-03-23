
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

#[get("/<channel>", rank = 1)]
pub fn get_trusts<'r>(handler: State<DbConn>, channel: String) -> Response<'r> {
	let trusts = handler.lock().expect("db lock").get_trusts(&channel);
	match trusts {
		Ok(trusts) => generate_response(Status::Ok, json!({ "data": trusts })),
		Err(HandlerError::Error(e)) => generate_response(Status::NotFound, generate_error(404, Some(e))),
		Err(e) => {
			println!("Internal error getting trusts: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}

#[get("/<channel>/<user>", rank = 1)]
pub fn get_trust<'r>(handler: State<DbConn>, channel: String, user: String) -> Response<'r> {
	let trust = handler.lock().expect("db lock").get_trust(&channel, &user);
	match trust {
		Ok(trust) => generate_response(Status::Ok, json!({ "data": trust })),
		Err(HandlerError::Error(e)) => generate_response(Status::NotFound, generate_error(404, Some(e))),
		Err(e) => {
			println!("Internal error getting trust: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}

#[post("/<channel>/<user>", rank = 2)]
pub fn create_trust<'r>(handler: State<DbConn>, channel: String, user: String) -> Response<'r> {
	let trust = handler.lock().expect("db lock").create_trust(&channel, &user);
	match trust {
		Ok(trust) => generate_response(Status::Ok, json!({ "data": trust })),
		Err(HandlerError::Error(e)) => generate_response(Status::NotFound, generate_error(404, Some(e))),
		Err(e) => {
			println!("Internal error creating trust: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}

#[delete("/<channel>/<user>", rank = 3)]
pub fn delete_trust<'r>(handler: State<DbConn>, channel: String, user: String) -> Response<'r> {
	let result = handler.lock().expect("db lock").delete_trust(&channel, &user);
	match result {
		Ok(()) => generate_response(Status::Ok, json! ({
			"meta": json! ({
				"deleted": true
			})
		})),
		Err(HandlerError::Error(e)) => generate_response(Status::NotFound, generate_error(404, Some(e))),
		Err(e) => {
			println!("Internal error deleting trust: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}
