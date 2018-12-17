
use std::{
	collections::HashMap,
	vec::Vec
};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Service {
	authorization: HashMap<String, String>,
	bot_name: String,
	channel: String,
	enabled: bool
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Channel {
	created_at: String,
	deleted_at: Option<String>,
	updated_at: String,
	token: String,
	id: String,
	enabled: bool,
	services: HashMap<String, Service>
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct CommandMeta {
	added_by: String,
	cooldown: u32,
	count: u32,
	enabled: bool
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Message {
	data: String,
	#[serde(rename = "type")]
	message_type: String
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Command {
	channel: String,
	created_at: String,
	deleted_at: Option<String>,
	id: String,
	meta: CommandMeta,
	name: String,
	response: Vec<Message>,
	services: Vec<String>,
	updated_at: String
}
