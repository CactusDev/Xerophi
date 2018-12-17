
use rocket_contrib::json::{JsonValue};

use rocket::State;
use crate::DbConn;

#[get("/<name>")]
pub fn get_channel(handler: State<DbConn>, name: String) -> JsonValue {
	let channel = handler.lock().expect("db lock").get_channel(&name);
	json!({
		"data": channel.map_or_else(|e| json!({ "hate_that": format!("{:?}", e) }), |data| json!(data))
	})
}

#[get("/<name>/command")]
pub fn get_commands(handler: State<DbConn>, name: String) -> JsonValue {
	let command = handler.lock().expect("db lock").get_command(&name, None);
	json!({
		"data": command.map_or_else(|_e| json!({}), |data| json!(data))
	})
}

#[get("/<name>/command/<command>")]
pub fn get_command(handler: State<DbConn>, name: String, command: String) -> JsonValue {
	let command = handler.lock().expect("db lock").get_command(&name, Some(command));
	json!({
		"data": command.map_or_else(|_e| json!({}), |data| json!(data))
	})
}
