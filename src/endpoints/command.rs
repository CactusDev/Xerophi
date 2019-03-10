
use rocket_contrib::json::Json;

use rocket::{
	State, Response,
	http::Status
};

use crate::{
	DbConn, endpoints::{generate_error, generate_response},
	database::{
		structures::Message,
		handler::HandlerError
	}
};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct PostCommand {
	pub response: Vec<Message>,
	pub services: Vec<String>
}

#[get("/<channel>")]
pub fn get_commands<'r>(handler: State<DbConn>, channel: String) -> Response<'r> {
	let commands = handler.lock().expect("db lock").get_commands(&channel);
	match commands {
		Ok(cmds) => generate_response(Status::Ok, json!({ "data": cmds })),
		Err(HandlerError::Error(_)) => generate_response(Status::Ok, json!({ "data": [] })),
		Err(e) => {
			println!("Internal error getting command: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}

#[get("/<channel>/<command>")]
pub fn get_command<'r>(handler: State<DbConn>, channel: String, command: String) -> Response<'r> {
	let command = handler.lock().expect("db lock").get_command(&channel, &command);
	match command {
		Ok(cmds) => generate_response(Status::Ok, json!({ "data": cmds })),
		Err(HandlerError::Error(_)) => generate_response(Status::NotFound, json!({})),
		Err(e) => {
			println!("Internal error getting command: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}

#[post("/<channel>/<name>", format = "json", data = "<command>")]
pub fn create_command<'r>(handler: State<DbConn>, channel: String, name: String, command: Json<PostCommand>) -> Response<'r> {
	let result = handler.lock().expect("db lock").create_command(&channel, &name, command.into_inner());
	match result {
		Ok(command) => generate_response(Status::Ok, json! ({ "data": command })),
		Err(HandlerError::Error(_)) => generate_response(Status::Conflict, generate_error(409, None)),
		Err(e) => {
			println!("Internal error creating command: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}

#[delete("/<channel>/<command>")]
pub fn delete_command<'r>(handler: State<DbConn>, channel: String, command: String) -> Response<'r> {
	let result = handler.lock().expect("db lock").remove_command(&channel, &command);
	match result {
		Ok(()) => generate_response(Status::Ok, json! ({
			"meta": json! ({
				"deleted": true
			})
		})),
		Err(HandlerError::Error(e)) => generate_response(Status::NotFound, generate_error(404, Some(e))),
		Err(e) => {
			println!("Internal error deleting command: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}

#[patch("/<channel>/<name>", format = "json", data = "<command>")]
pub fn edit_command<'r>(handler: State<DbConn>, channel: String, name: String, command: Json<PostCommand>) -> Response<'r> {
	let result = handler.lock().expect("db lock").update_command(&channel, &name, command.into_inner());
	match result {
		Ok(()) => generate_response(Status::Ok, json! ({
			"meta": json! ({
				"updated": true
			})
		})),
		Err(HandlerError::Error(e)) => generate_response(Status::NotFound, generate_error(404, Some(e))),
		Err(e) => {
			println!("Internal error updating command: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}
