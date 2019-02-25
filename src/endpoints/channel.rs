
use rocket_contrib::json::Json;

use rocket::{
	State, Response,
	http::Status
};

use crate::{
	DbConn, endpoints::{generate_error, generate_response},
	database::handler::HandlerError
};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct PostChannel {
	pub name: String,
	pub password: String
}

#[get("/<channel>")]
pub fn get_channel<'r>(handler: State<DbConn>, channel: String) -> Response<'r> {
	let channel = handler.lock().expect("db lock").get_channel(&channel);
	match channel {
		Ok(channel) => generate_response(Status::Ok, json!({ "data": channel })),
		Err(HandlerError::Error(e)) => generate_response(Status::NotFound, generate_error(404, Some(e))),
		Err(e) => {
			println!("Internal error getting channel: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}

#[post("/create", format = "json", data = "<channel>")]
pub fn create_channel<'r>(handler: State<DbConn>, channel: Json<PostChannel>) -> Response<'r> {
	let result = handler.lock().expect("db lock").create_channel(channel.into_inner());
	match result {
		Ok(channel) => generate_response(Status::Ok, json!({
			"created": true,
			"token": channel.token
		})),
		Err(HandlerError::Error(e)) => generate_response(Status::Conflict, generate_error(409, Some(e))),
		Err(e) => {
			println!("Internal error creating channel: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}
