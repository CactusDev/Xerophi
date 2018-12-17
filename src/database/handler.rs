
use mongodb::{
	Client, ThreadedClient, doc, bson,
	db::{Database, ThreadedDatabase}
};

use bson::from_bson;

use super::structures::*;

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

	pub fn connect(&mut self, database: &str, _username: &str, _password: &str) -> Result<(), mongodb::Error> {
		match Client::with_uri(&self.url) {
			Ok(client) => {
				let database = client.db(database);
				self.database = Some(database);
				Ok(())
			},
			Err(err) => Err(err)
		}
	}

	pub fn get_channel(&self, name: &str) -> Result<Channel, mongodb::Error> {
		match &self.database {
			Some(db) => {
				let filter = doc! { "token": name };

				let channel_collection = db.collection("channels");
				let cursor = channel_collection.find_one(Some(filter), None);

				match cursor {
					Ok(Some(channel)) => Ok(from_bson::<Channel>(mongodb::Bson::Document(channel)).unwrap()),
					Ok(None) => Err(mongodb::Error::DefaultError("no channel".to_string())),
					Err(e) => Err(e)
				}
			},
			None => Err(mongodb::Error::DefaultError("no database".to_string()))
		}
	}

	pub fn get_command(&self, channel: &str, command: Option<String>) -> Result<Vec<Command>, mongodb::Error> {
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
				let mut cursor = command_collection.find(Some(filter), None)?;

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
					return Err(mongodb::Error::DefaultError("no command".to_string()))
				}
				Ok(all_documents)
			},
			None => Err(mongodb::Error::DefaultError("no database".to_string()))
		}
	}
}
