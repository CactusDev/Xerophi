
use rocket_contrib::json::{JsonValue, Json};

use rocket::State;
use crate::{
	DbConn, endpoints::generate_error,
	database::{
		structures::Message,
		handler::HandlerError
	}
};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct PostCommand {
	pub name: String,
	pub response: Vec<Message>,
	pub services: Vec<String>
}

#[get("/<channel>")]
pub fn get_channel(handler: State<DbConn>, channel: String) -> JsonValue {
	let channel = handler.lock().expect("db lock").get_channel(&channel);
	match channel {
		Ok(channel) => json!({ "data": channel }),
		Err(HandlerError::Error(e)) => generate_error(404, Some(e)),
		Err(e) => {
			println!("Internal error getting channel: {:?}", e);
			generate_error(500, None)
		}
	}
}

#[get("/<channel>/command")]
pub fn get_commands(handler: State<DbConn>, channel: String) -> JsonValue {
	let commands = handler.lock().expect("db lock").get_command(&channel, None);
	match commands {
		Ok(cmds) => json!({ "data": cmds }),
		Err(HandlerError::Error(e)) => generate_error(404, Some(e)),
		Err(e) => {
			println!("Internal error getting command: {:?}", e);
			generate_error(500, None)
		}
	}
}

#[get("/<channel>/command/<command>")]
pub fn get_command(handler: State<DbConn>, channel: String, command: String) -> JsonValue {
	let command = handler.lock().expect("db lock").get_command(&channel, Some(command));
	match command {
		Ok(cmds) => json!({ "data": cmds[0] }),
		Err(HandlerError::Error(e)) => generate_error(404, Some(e)),
		Err(e) => {
			println!("Internal error getting command: {:?}", e);
			generate_error(500, None)
		}
	}
}

#[post("/<channel>/command/create", format = "json", data = "<command>")]
pub fn create_command(handler: State<DbConn>, channel: String, command: Json<PostCommand>) -> JsonValue {
	let result = handler.lock().expect("db lock").create_command(&channel, command.into_inner());
	match result {
		Ok(command) => json! ({
			"data": command
		}),
		Err(HandlerError::Error(e)) => generate_error(401, Some(e)),
		Err(e) => {
			println!("Internal error creating command: {:?}", e);
			generate_error(500, None)
		}
	}
}

#[delete("/<channel>/command/<command>")]
pub fn delete_command(handler: State<DbConn>, channel: String, command: String) -> JsonValue {
	let result = handler.lock().expect("db lock").remove_command(&channel, &command);
	match result {
		Ok(()) => json! ({
			"meta": json! ({
				"deleted": true
			})
		}),
		Err(HandlerError::Error(e)) => generate_error(404, Some(e)),
		Err(e) => {
			println!("Internal error deleting command: {:?}", e);
			generate_error(500, None)
		}
	}
}

#[get("/<channel>/config")]
pub fn get_config(handler: State<DbConn>, channel: String) -> JsonValue {
	let result = handler.lock().expect("db lock").get_config(&channel);
	match result {
		Ok(config) => json! ({
			"data": config
		}),
		Err(HandlerError::Error(e)) => generate_error(404, Some(e)),
		Err(e) => {
			println!("Internal error getting config: {:?}", e);
			generate_error(500, None)
		}
	}
}
