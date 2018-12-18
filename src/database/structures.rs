
use crate::endpoints::channel::PostCommand;
use rocket::data::FromData;

use std::{
	collections::HashMap,
	vec::Vec
};

use chrono::prelude::*;

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
	cooldown: i32,
	count: i32,
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
	id: Option<String>,
	meta: CommandMeta,
	name: String,
	response: Vec<Message>,
	services: Vec<String>,
	updated_at: String
}

impl Command {

	pub fn from_post(cmd: PostCommand, channel: &str) -> Command {
		let the_time = Local::now().to_string();

		Command {
			channel: channel.to_string(),
			created_at: the_time.clone(),
			deleted_at: None,
			id: None,
			meta: CommandMeta {
				added_by: "".to_string(),
				cooldown: 0,
				count: 0,
				enabled: true
			},
			name: cmd.name,
			response: cmd.response,
			services: cmd.services,
			updated_at: the_time
		}
	}
}
