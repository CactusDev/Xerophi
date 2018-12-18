
use mongodb::{
	Client, ThreadedClient, doc, bson,
	db::{Database, ThreadedDatabase}
};

use bson::{to_bson, from_bson};

use super::structures::*;
use crate::endpoints::channel::PostCommand;

#[derive(Debug)]
pub enum HandlerError {
	DatabaseError(mongodb::Error),
	InternalError,
	Error(String)
}

type HandlerResult<T> = Result<T, HandlerError>;

pub struct DatabaseHandler {
	url: String,
	database: Option<Database>
}

impl DatabaseHandler {

	pub fn new(host: &str, port: u16) -> Self {
		DatabaseHandler {
			url: format!("mongodb://{}:{}", host, port),
			database: None
		}
	}

	pub fn connect(&mut self, database: &str, _username: &str, _password: &str) -> HandlerResult<()> {
		match Client::with_uri(&self.url) {
			Ok(client) => {
				let database = client.db(database);
				self.database = Some(database);
				Ok(())
			},
			Err(err) => Err(HandlerError::DatabaseError(err))
		}
	}

	pub fn get_channel(&self, name: &str) -> HandlerResult<Channel> {
		match &self.database {
			Some(db) => {
				let filter = doc! { "token": name };

				let channel_collection = db.collection("channels");
				let cursor = channel_collection.find_one(Some(filter), None);

				match cursor {
					Ok(Some(channel)) => Ok(from_bson::<Channel>(mongodb::Bson::Document(channel)).unwrap()),
					Ok(None) => Err(HandlerError::Error("no channel".to_string())),
					Err(e) => Err(HandlerError::DatabaseError(e))
				}
			},
			None => Err(HandlerError::InternalError)
		}
	}

	pub fn get_command(&self, channel: &str, command: Option<String>) -> HandlerResult<Vec<Command>> {
		let filter = match command {
			Some(cmd) => doc! {
				"name": cmd,
				"channel": channel
			},
			None => doc! {
				"channel": channel
			}
		};

		match &self.database {
			Some(db) => {
				let command_collection = db.collection("commands");
				let mut cursor = command_collection.find(Some(filter), None).map_err(|e| HandlerError::DatabaseError(e))?;

				let mut all_documents: Vec<Command> = vec! [];

				while cursor.has_next().unwrap_or(false) {
					let doc = cursor.next_n(1);
					match doc {
					 	Ok(ref docs) => for doc in docs {
					 		all_documents.push(from_bson::<Command>(mongodb::Bson::Document(doc.clone())).unwrap());
					 	},
					 	Err(_) => break
					 }
				}
				// TODO: no
				if all_documents.len() == 0 {
					return Err(HandlerError::Error("no command".to_string()))
				}
				Ok(all_documents)
			},
			None => Err(HandlerError::InternalError)
		}
	}

	pub fn create_command(&self, channel: &str, command: PostCommand) -> HandlerResult<Command> {
		if let Ok(_) = self.get_command(channel, Some(command.name.clone())) {
			return Err(HandlerError::Error("command exists".to_string()));
		}

		match &self.database {
			Some(db) => {
				let command_collection = db.collection("commands");
				let command = Command::from_post(command, channel);

				command_collection.insert_one(to_bson(&command).unwrap().as_document().unwrap().clone(), None).map_err(|e| HandlerError::DatabaseError(e))?;
				Ok(command)
			},
			None => Err(HandlerError::InternalError)
		}
	}
}
