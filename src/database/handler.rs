
use mongodb::{
	Client, ThreadedClient, doc, bson,
	db::{Database, ThreadedDatabase},
	coll::options::{UpdateOptions}
};

use bson::{to_bson, from_bson};

use argon2::{Config as ArgonConfig, hash_encoded};

use super::structures::*;
use crate::endpoints::{
	channel::PostChannel,
	quote::PostQuote,
	authorization::PostServiceAuth,
	command::PostCommand
};

#[derive(Debug)]
pub enum HandlerError {
	DatabaseError(mongodb::Error),
	InternalError,
	Error(String)
}

type HandlerResult<T> = Result<T, HandlerError>;

pub struct DatabaseHandler<'cfg> {
	url: String,
	database: Option<Database>,
	argon: ArgonConfig<'cfg>,
	salt: String
}

impl<'cfg> DatabaseHandler<'cfg> {

	pub fn new(host: &str, port: u16, key: &'cfg str, salt: &str) -> Self {
		let mut argon = ArgonConfig::default();
		argon.secret = key.as_bytes();

		DatabaseHandler {
			url: format!("mongodb://{}:{}", host, port),
			database: None,
			argon,
			salt: salt.to_string()
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
        let db = self.database.as_ref().expect("no database");
		let filter = doc! { "token": name };

		let channel_collection = db.collection("channels");
		let cursor = channel_collection.find_one(Some(filter), None);

		match cursor {
			Ok(Some(channel)) => Ok(from_bson::<Channel>(mongodb::Bson::Document(channel)).unwrap()),
			Ok(None) => Err(HandlerError::Error("no channel".to_string())),
			Err(e) => Err(HandlerError::DatabaseError(e))
		}
	}

	pub fn create_channel(&self, channel: PostChannel) -> HandlerResult<Channel> {
		if let Ok(_) = self.get_channel(&channel.name) {
			return Err(HandlerError::Error("channel with name exists".to_string()));
		}

		let password = hash_encoded(&channel.password.clone().into_bytes(), &self.salt.clone().into_bytes(), &self.argon);
		if let Err(_) = password {
			return Err(HandlerError::Error("could not hash password".to_string()));
		}

        let db = self.database.as_ref().expect("no database");
		let channel_collection = db.collection("channels");
		match Channel::from_post(channel, password.unwrap()) {
			Some(channel) => {
				channel_collection.insert_one(to_bson(&channel).unwrap().as_document().unwrap().clone(), None)
					.map_err(|e| HandlerError::DatabaseError(e))?;
				Ok(channel)
			},
			None => Err(HandlerError::InternalError)
		}
	}

	pub fn get_commands(&self, channel: &str) -> HandlerResult<Vec<Command>> {
		let filter = doc! { "channel": channel };

        let db = self.database.as_ref().expect("no database");
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
	}

	pub fn get_command(&self, channel: &str, command: &str) -> HandlerResult<Command> {
		let filter = doc! {
			"name": command,
			"channel": channel
		};

        let db = self.database.as_ref().expect("no database");
		let command_collection = db.collection("commands");
		let document = command_collection.find_one(Some(filter), None).map_err(|e| HandlerError::DatabaseError(e))?;
		match document {
			Some(doc) => Ok(from_bson::<Command>(mongodb::Bson::Document(doc.clone())).unwrap()),
			_ => Err(HandlerError::Error("no command".to_string()))
		}
	}

	pub fn create_command(&self, channel: &str, name: &str, command: PostCommand) -> HandlerResult<Command> {
		if let Ok(_) = self.get_command(channel, name) {
			return Err(HandlerError::Error("command exists".to_string()));
		}

        let db = self.database.as_ref().expect("no database");
		let command_collection = db.collection("commands");
		let command = Command::from_post(command, channel, name);

		command_collection.insert_one(to_bson(&command).unwrap().as_document().unwrap().clone(), None).map_err(|e| HandlerError::DatabaseError(e))?;
		Ok(command)
	}

	pub fn remove_command(&self, channel: &str, command: &str) -> HandlerResult<()> {
		if let Err(_) = self.get_command(channel, command) {
			return Err(HandlerError::Error("command does not exist".to_string()));
		}

        let db = self.database.as_ref().expect("no database");
		let command_collection = db.collection("commands");
		command_collection.delete_one(doc! {
			"channel": channel,
			"name": command
		}, None).map_err(|e| HandlerError::DatabaseError(e))?;
		Ok(())
	}

	pub fn update_command(&self, channel: &str, name: &str, command: PostCommand) -> HandlerResult<()> {
		if let Err(_) = self.get_command(channel, name) {
			return Err(HandlerError::Error("command does not exist".to_string()));
		}

        let db = self.database.as_ref().expect("no database");
		let command_collection = db.collection("commands");
		let command = to_bson(&command.response).unwrap();

		match command_collection.update_one(doc! {
			"channel": channel,
			"name": name
		}, doc! {
			"$set": doc! {
				"response": command
			}
		}, None) {
			Ok(_) => Ok(()),
			Err(e) => Err(HandlerError::DatabaseError(e))
		}
	}

	pub fn get_config(&self, channel: &str) -> HandlerResult<Config> {
        let db = self.database.as_ref().expect("no database");
		let config_collection = db.collection("configs");
		match config_collection.find_one(Some(doc! {
			"channel": channel
		}), None) {
			Ok(Some(config)) => Ok(from_bson::<Config>(mongodb::Bson::Document(config)).unwrap()),
			Ok(None) => Err(HandlerError::Error("no channel".to_string())),
			Err(e) => Err(HandlerError::DatabaseError(e))
		}
	}

	pub fn get_channel_state(&self, channel: &str, service: Option<String>) -> HandlerResult<BotState> {
		let filter = match service {
			Some(service) => doc! {
				"service": service,
				"channel": channel
			},
			None => doc! { "channel": channel }
		};

        let db = self.database.as_ref().expect("no database");
		let state_collection = db.collection("state");
		match state_collection.find_one(Some(filter), None) {
			Ok(Some(state)) => Ok(from_bson::<BotState>(mongodb::Bson::Document(state)).unwrap()),
			Ok(None) => Err(HandlerError::Error("no state for provided channel".to_string())),
			Err(e) => Err(HandlerError::DatabaseError(e))
		}
	}

	pub fn get_service_auth(&self, channel: &str, service: &str) -> HandlerResult<BotAuthorization> {
		let filter = doc! {
			"service": service,
			"channel": channel
		};

        let db = self.database.as_ref().expect("no database");
		let authorization_collection = db.collection("authorization");
		match authorization_collection.find_one(Some(filter), None) {
			Ok(Some(auth)) => Ok(from_bson::<BotAuthorization>(mongodb::Bson::Document(auth)).unwrap()),
			Ok(None) => Err(HandlerError::Error("no auth for service".to_string())),
			Err(e) => Err(HandlerError::DatabaseError(e))
		}
	}

	pub fn update_service_auth(&self, channel: &str, service: &str, auth: PostServiceAuth) -> HandlerResult<()> {
		let filter = doc! {
			"service": service,
			"channel": channel
		};

		let mut opts = UpdateOptions::new();
		opts.upsert = Some(true);

        let db = self.database.as_ref().expect("no database");
		let authorization_collection = db.collection("authorization");
		match authorization_collection.update_one(filter, doc! {
			"refresh": &auth.refresh.unwrap_or("".into()),
			"expires": &auth.expiration.unwrap_or("".into()),
			"access": &auth.access
		}, Some(opts)) {
			Ok(_) => Ok(()),
			Err(e) => Err(HandlerError::DatabaseError(e))
		}
	}

	pub fn get_quote(&self, channel: &str, id: Option<u32>) -> HandlerResult<Vec<Quote>> {
		let filter = match id {
		    Some(id) => doc! {
				"channel": channel,
				"quote_id": id
			},
			None => doc! { "channel": channel }
		};

        let db = self.database.as_ref().expect("no database");
		let quote_collection = db.collection("quotes");

		match quote_collection.find(Some(filter), None) {
			Ok(mut quotes) => {
				let mut all_documents = Vec::new();
				while quotes.has_next().unwrap_or(false) {
					let doc = quotes.next_n(1);
					match doc {
						Ok(ref docs) => for doc in docs {
				 	 		all_documents.push(from_bson::<Quote>(mongodb::Bson::Document(doc.clone())).unwrap());
						},
				 	 	Err(_) => break
					}
				}
                // TODO: no
				if all_documents.len() == 0 {
					return Err(HandlerError::Error("no quote".to_string()))
				}
				Ok(all_documents)
			},
			Err(e) => Err(HandlerError::DatabaseError(e))
		}
	}

	pub fn get_random_quote(&self, channel: &str) -> HandlerResult<Quote> {
		let pipeline = vec![
		    doc! { "$match": doc! { "channel": channel } },
		    doc! { "$sample": doc! { "size": 1 } }
		];

        let db = self.database.as_ref().expect("no database");
		let quote_collection = db.collection("quotes");
		match quote_collection.aggregate(pipeline, None) {
			Ok(mut cursor) => match cursor.drain_current_batch() {
				Ok(batch) => match batch.as_slice() {
					[first] => Ok(from_bson::<Quote>(mongodb::Bson::Document(first.clone())).unwrap()),
					_ => Err(HandlerError::Error("no quotes".to_string()))
				}
				Err(_) => Err(HandlerError::InternalError)
			},
			Err(e) => Err(HandlerError::DatabaseError(e))
		}
	}

	pub fn create_quote(&self, channel: &str, quote: PostQuote) -> HandlerResult<i64> {
		let filter = doc! {
			"channel": channel
		};

		let db = self.database.as_ref().expect("no database");
		let quote_collection = db.collection("quotes");
		let mut count = quote_collection.count(Some(filter), None)
			.map_err(|e| HandlerError::DatabaseError(e))?;
		count += 1;  // Set the ID for the next quote. Also prevents quote 1 from having id 0.

		// Create the new quote.
		let quote = Quote {
			quote_id: count,
			response: quote.response,
			channel: channel.to_string()
		};
		quote_collection.insert_one(to_bson(&quote).unwrap().as_document().unwrap().clone(), None)
			.map_err(|e| HandlerError::DatabaseError(e))?;
		Ok(count)
	}

	pub fn delete_quote(&self, channel: &str, quote: u32) -> HandlerResult<()> {
		let db = self.database.as_ref().expect("no database");
		let quote_collection = db.collection("quotes");

		quote_collection.delete_one(doc! {
			"channel": channel,
			"quote_id": quote
		}, None).map_err(|e| HandlerError::DatabaseError(e))?;
		Ok(())
	}

	pub fn create_trust(&self, channel: &str, user: &str) -> HandlerResult<Trust> {
		let db = self.database.as_ref().expect("no database");
		let trust_collection = db.collection("trusts");

		let trust = Trust {
			trusted: user.to_string(),
			channel: channel.to_string()
		};

		trust_collection.insert_one(to_bson(&trust).unwrap().as_document().unwrap().clone(), None)
			.map_err(|e| HandlerError::DatabaseError(e))?;
		Ok(trust)
	}

	pub fn get_trust(&self, channel: &str, user: &str) -> HandlerResult<Trust> {
		let filter = doc! {
			"trusted": user,
			"channel": channel
		};

        let db = self.database.as_ref().expect("no database");
		let trust_collection = db.collection("trusts");
		match trust_collection.find_one(Some(filter), None) {
			Ok(Some(trust)) => Ok(from_bson::<Trust>(mongodb::Bson::Document(trust)).unwrap()),
			Ok(None) => Err(HandlerError::Error("no trust for user".to_string())),
			Err(e) => Err(HandlerError::DatabaseError(e))
		}
	}

	pub fn delete_trust(&self, channel: &str, user: &str) -> HandlerResult<()> {
		let db = self.database.as_ref().expect("no database");
		let trust_collection = db.collection("trusts");

		trust_collection.delete_one(doc! {
			"channel": channel,
			"trusted": user
		}, None).map_err(|e| HandlerError::DatabaseError(e))?;
		Ok(())
	}
}
