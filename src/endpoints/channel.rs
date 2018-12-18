
use rocket_contrib::json::{JsonValue, Json};

use rocket::{State, FromForm};
use crate::{
	DbConn, endpoints::generate_error,
	database::structures::Message
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
		Err(_) => generate_error(404)
	}
}

#[get("/<channel>/command")]
pub fn get_commands(handler: State<DbConn>, channel: String) -> JsonValue {
	let commands = handler.lock().expect("db lock").get_command(&channel, None);
	match commands {
		Ok(cmds) => json!({ "data": cmds }),
		Err(_) => generate_error(404)
	}
}

#[get("/<channel>/command/<command>")]
pub fn get_command(handler: State<DbConn>, channel: String, command: String) -> JsonValue {
	let command = handler.lock().expect("db lock").get_command(&channel, Some(command));
	match command {
		Ok(cmds) => json!({ "data": cmds[0] }),
		Err(_) => generate_error(404)
	}
}

#[post("/<channel>/command/create", format = "json", data = "<command>")]
pub fn create_command(handler: State<DbConn>, channel: String, command: Json<PostCommand>) -> JsonValue {
	let result = handler.lock().expect("db lock").create_command(&channel, command.into_inner());
	json! ({})
}
