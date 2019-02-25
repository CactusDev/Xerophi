
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

#[get("/<channel>/command")]
pub fn get_commands<'r>(handler: State<DbConn>, channel: String) -> Response<'r> {
	let commands = handler.lock().expect("db lock").get_command(&channel, None);
	match commands {
		Ok(cmds) => generate_response(Status::Ok, json!({ "data": cmds })),
		Err(HandlerError::Error(_)) => generate_response(Status::Ok, json!([])),
		Err(e) => {
			println!("Internal error getting command: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}

#[get("/<channel>/command/<command>")]
pub fn get_command<'r>(handler: State<DbConn>, channel: String, command: String) -> Response<'r> {
	let command = handler.lock().expect("db lock").get_command(&channel, Some(command));
	match command {
		Ok(cmds) => generate_response(Status::Ok, json!({ "data": cmds[0] })),
		Err(HandlerError::Error(_)) => generate_response(Status::NotFound, json!({})),
		Err(e) => {
			println!("Internal error getting command: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}

#[post("/<channel>/command/<name>", format = "json", data = "<command>")]
pub fn create_command<'r>(handler: State<DbConn>, channel: String, name: String, command: Json<PostCommand>) -> Response<'r> {
	let result = handler.lock().expect("db lock").create_command(&channel, &name, command.into_inner());
	match result {
		Ok(command) => generate_response(Status::Ok, json! ({ "data": command })),
		Err(HandlerError::Error(e)) => generate_response(Status::NotFound, generate_error(404, Some(e))),
		Err(e) => {
			println!("Internal error creating command: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}

#[delete("/<channel>/command/<command>")]
pub fn delete_command<'r>(handler: State<DbConn>, channel: String, command: String) -> Response<'r> {
	let result = handler.lock().expect("db lock").remove_command(&channel, &command);
	match result {
		Ok(()) => generate_response(Status::Ok, json! ({
			"meta": json! ({
				"deleted": true
			})
		})),
		Err(HandlerError::Error(e)) => generate_response(Status::Conflict, generate_error(409, Some(e))),
		Err(e) => {
			println!("Internal error deleting command: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}
