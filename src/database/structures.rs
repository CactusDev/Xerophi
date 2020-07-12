
use crate::endpoints::{
	channel::PostChannel,
	command::PostCommand
};

use std::{
	vec::Vec
};

use chrono::prelude::*;

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Channel {
	#[serde(rename = "_id")]
	pub id: bson::oid::ObjectId,
	pub created_at: String,
	pub deleted_at: Option<String>,
	pub updated_at: String,
	pub token: String,
	pub enabled: bool,
	pub password: String
}

impl Channel {

	pub fn from_post(channel: PostChannel, password: String) -> Option<Self> {
		let the_time = Local::now().to_string();
		let id = bson::oid::ObjectId::new();

		match id {
			Ok(id) => {
				Some(Channel {
					id,
					created_at: the_time.clone(),
					deleted_at: None,
					updated_at: the_time,
					token: channel.name,
					enabled: true,
					password
   			    })
			},
			Err(_) => None
		}
	}
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct CommandMeta {
	pub added_by: String,
	pub cooldown: i32,
	pub count: i32,
	pub enabled: bool,
	pub role: String
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct EmojiMessageData {
	pub standard: String,
	pub alternatives: Vec<String>,	
}

#[derive(Serialize, Deserialize, Clone, Debug)]
#[serde(untagged)]
pub enum GenericMessageData {
	Basic(String),
	Emoji(EmojiMessageData)
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Message {
	pub data: GenericMessageData,
	#[serde(rename = "type")]
	pub message_type: String
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Command {
	pub channel: String,
	pub created_at: String,
	pub deleted_at: Option<String>,
	pub meta: CommandMeta,
	pub name: String,
	pub response: Vec<Message>,
	pub services: Vec<String>,
	pub updated_at: String
}

impl Command {

	pub fn from_post(cmd: PostCommand, channel: &str, name: &str) -> Command {
		let the_time = Local::now().to_string();

		Command {
			channel: channel.to_string(),
			created_at: the_time.clone(),
			deleted_at: None,
			meta: CommandMeta {
				added_by: "".to_string(),
				cooldown: 0,
				count: 0,
				enabled: true,
				role: cmd.role
			},
			name: name.to_string(),
			response: cmd.response,
			services: cmd.services,
			updated_at: the_time
		}
	}
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct RepeatConfig {
	pub disabled: bool,
	pub only_live: bool,
	pub default_minimum: i32
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct EventConfig {
	pub message: String,
	pub enabled: bool
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct EventsConfig {
	pub follow: EventConfig,
	pub subscribe: EventConfig,
	pub host: EventConfig,
	pub join: EventConfig,
	pub leave: EventConfig
}

#[derive(Serialize, Deserialize, Clone, Debug)]
#[serde(rename_all = "snake_case")]
pub enum SpamAction {
	Ignore, Purge, Timeout, Ban
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct SpamConfigs<T> {
	pub action: SpamAction,
	pub value: T,
	pub warnings: i32
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct SpamKeywords {
	pub blacklist: Vec<String>,
	pub whitelist: Vec<String>,
	pub urls: Vec<String>
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct SpamConfig {
	pub allow_urls: SpamConfigs<bool>,
	pub max_caps_score: SpamConfigs<i32>,
	pub max_emoji: SpamConfigs<i32>,
	pub keywords: SpamKeywords
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Config {
	pub repeat: RepeatConfig,
	pub events: EventsConfig,
	pub spam: SpamConfig,
	pub channel: String
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct BotAuthorization {
	pub access: String,
	pub refresh: Option<String>,
	pub expiration: Option<String>
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct ConnectedService {
	pub service: String,
	pub connected: bool,
	pub last_authorization: i32
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct BotState {
	pub services: Vec<ConnectedService>,
	pub token: String
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Quote {
	pub quote_id: i64,
	pub response: Vec<Message>,
	pub channel: String
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Trust {
	pub trusted: String,
	pub channel: String
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Alias {
	pub channel: String,
	pub command: String,
	pub alias: String
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct SocialService {
	pub channel: String,
	pub service: String,
	pub url: String
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct UserOffences {
	pub channel: String,
	pub service: String,
	pub user: String,
	pub caps: i32,
	pub emoji: i32,
	pub urls: i32
}

impl UserOffences {

	pub fn get_attribute(&self, name: &str) -> Option<i32> {
		match name {
			"caps" => Some(self.caps),
			"emoji" => Some(self.emoji),
			"urls" => Some(self.urls),
			_ => None
		}
	}

	pub fn set_attribute(mut self, name: &str, value: i32) -> Result<Self, ()>{
		match name {
			"caps" => self.caps = value,
			"emoji" => self.emoji = value,
			"urls" => self.urls = value,
			_ => return Err(())
		};
		Ok(self)
	}
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct UpdateCount {
	pub count: String
}
