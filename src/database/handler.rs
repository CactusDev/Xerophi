
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
	user: String,
	password: String,
	database: Option<Database>,
	argon: ArgonConfig<'cfg>,
	salt: String
}

impl<'cfg> DatabaseHandler<'cfg> {

	pub fn new(username: &str, password: &str, host: &str, port: u32, db: &str, key: &'cfg str, salt: &str) -> Self {
		let mut argon = ArgonConfig::default();
		argon.secret = key.as_bytes();

		DatabaseHandler {
			url: format!("mongodb://{}:{}/{}", host, port, db),
			user: username.to_string(),
			password: password.to_string(),
			database: None,
			argon,
			salt: salt.to_string()
		}
	}

	pub fn connect(&mut self, database: &str, _username: &str, _password: &str) -> HandlerResult<()> {
		match Client::with_uri(&self.url) {
			Ok(client) => {
				let database = client.db(database);
				database.auth(&self.user, &self.password).map_err(|e| HandlerError::DatabaseError(e))?;
				
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
			Some(doc) => match from_bson::<Command>(mongodb::Bson::Document(doc.clone())) {
				Ok(cmd) => Ok(cmd),
				_ => Err(HandlerError::Error("no command".to_string()))
			},
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
		// Also remove any aliases associated with it
		let alias_collection = db.collection("aliases");
		alias_collection.delete_one(doc! {
			"channel": channel,
			"command": command
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

	pub fn update_count(&self, channel: &str, command: &str, count: UpdateCount) -> HandlerResult<i32> {
		let db = self.database.as_ref().expect("no database");
		let command_collection = db.collection("commands");
		let cmd = self.get_command(channel, command)?;
		let mut current = cmd.meta.count;


		let (operator, remaining) = count.count.split_at(1);
		let remaining = remaining.parse::<i32>().map_err(|_| HandlerError::Error("invalid count".to_string()))?;

		match operator {
			"+" => current += remaining,
			"-" => current -= remaining,
			"@" => current = remaining,
			_ => return Err(HandlerError::Error("invalid operator".to_string()))
		}

		match command_collection.update_one(doc! {
			"channel": channel,
			"name": command
		}, doc! {
			"$set": doc! {
				"count": current
			}
		}, None) {
			Ok(_) => Ok(current),
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
			"$set": doc! {
				"refresh": &auth.refresh.unwrap_or("".into()),
				"expires": &auth.expiration.unwrap_or("".into()),
				"access": &auth.access
			}
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

	pub fn edit_quote(&self, channel: &str, id: u32, quote: PostQuote) -> HandlerResult<()> {
		let filter = doc! {
			"quote_id": id,
			"channel": channel
		};

		let db = self.database.as_ref().expect("no database");
		let quote_collection = db.collection("quotes");
		let quote_response = to_bson(&quote.response).unwrap();

		match quote_collection.update_one(filter, doc! {
			"$set": doc! {
				"response": quote_response
			}
		}, None) {
			Ok(_) => Ok(()),
			Err(e) => Err(HandlerError::DatabaseError(e))
		}
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
		if let Ok(_) = self.get_trust(channel, user) {
			return Err(HandlerError::Error("already trusted".to_string()));
		}

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

	pub fn get_trusts(&self, channel: &str) -> HandlerResult<Vec<Trust>> {
		let filter = doc! {
			"channel": channel,
		};

        let db = self.database.as_ref().expect("no database");
		let trust_collection = db.collection("trusts");

		match trust_collection.find(Some(filter), None) {
			Ok(mut trusts) => {
				let mut all_documents = Vec::new();
				while trusts.has_next().unwrap_or(false) {
					let doc = trusts.next_n(1);
					match doc {
						Ok(ref docs) => for doc in docs {
							let trust = from_bson::<Trust>(mongodb::Bson::Document(doc.clone())).unwrap();
				 	 		all_documents.push(trust);
						},
				 	 	Err(_) => break
					}
				}
                // TODO: no
				if all_documents.len() == 0 {
					return Err(HandlerError::Error("no trusts".to_string()))
				}
				Ok(all_documents)
			},
			Err(e) => Err(HandlerError::DatabaseError(e))
		}
	}

	pub fn delete_trust(&self, channel: &str, user: &str) -> HandlerResult<()> {
		if let Err(_) = self.get_trust(channel, user) {
			return Err(HandlerError::Error("trust does not exist".to_string()));
		}

		let db = self.database.as_ref().expect("no database");
		let trust_collection = db.collection("trusts");

		trust_collection.delete_one(doc! {
			"channel": channel,
			"trusted": user
		}, None).map_err(|e| HandlerError::DatabaseError(e))?;
		Ok(())
	}

	pub fn get_aliases(&self, channel: &str) -> HandlerResult<Vec<Command>> {
		let filter = doc! {
			"channel": channel,
		};

        let db = self.database.as_ref().expect("no database");
		let alias_collection = db.collection("aliases");

		match alias_collection.find(Some(filter), None) {
			Ok(mut commands) => {
				let mut all_documents = Vec::new();
				while commands.has_next().unwrap_or(false) {
					let doc = commands.next_n(1);
					match doc {
						Ok(ref docs) => for doc in docs {
							let alias = from_bson::<Alias>(mongodb::Bson::Document(doc.clone())).unwrap();
				 	 		let command = self.get_command(channel, &alias.command)?;
				 	 		all_documents.push(command);
						},
				 	 	Err(_) => break
					}
				}
                // TODO: no
				if all_documents.len() == 0 {
					return Err(HandlerError::Error("no aliases".to_string()))
				}
				Ok(all_documents)
			},
			Err(e) => Err(HandlerError::DatabaseError(e))
		}
	}

	pub fn get_alias(&self, channel: &str, command: &str) -> HandlerResult<Command> {
		let filter = doc! {
			"channel": channel,
			"command": command
		};

		let db = self.database.as_ref().expect("no database");
		let alias_collection = db.collection("aliases");
		
		let document = alias_collection.find_one(Some(filter), None).map_err(|e| HandlerError::DatabaseError(e))?;
		match document {
			Some(doc) => {
				let alias = from_bson::<Alias>(mongodb::Bson::Document(doc.clone())).unwrap();
				Ok(self.get_command(channel, &alias.command)?)
			},
			_ => Err(HandlerError::Error("no alias".to_string()))
		}
	}

	pub fn create_alias(&self, channel: &str, alias: &str, command: &str) -> HandlerResult<Alias> {
		let db = self.database.as_ref().expect("no database");
		let alias_collection = db.collection("aliases");

		let alias = Alias {
			alias: alias.to_string(),
			command: command.to_string(),
			channel: channel.to_string()
		};

		alias_collection.insert_one(to_bson(&alias).unwrap().as_document().unwrap().clone(), None)
			.map_err(|e| HandlerError::DatabaseError(e))?;
		Ok(alias)
	}

	pub fn delete_alias(&self, channel: &str, command: &str) -> HandlerResult<()> {
		let db = self.database.as_ref().expect("no database");
		let alias_collection = db.collection("aliases");

		alias_collection.delete_one(doc! {
			"channel": channel,
			"alias": command
		}, None).map_err(|e| HandlerError::DatabaseError(e))?;
		Ok(())
	}

	pub fn get_socials(&self, channel: &str) -> HandlerResult<Vec<SocialService>> {
		let filter = doc! { "channel": channel };

        let db = self.database.as_ref().expect("no database");
		let social_collection = db.collection("socials");
		let mut cursor = social_collection.find(Some(filter), None).map_err(|e| HandlerError::DatabaseError(e))?;

		let mut all_documents = vec! [];

		while cursor.has_next().unwrap_or(false) {
			let doc = cursor.next_n(1);
			match doc {
			 	Ok(ref docs) => for doc in docs {
			 		all_documents.push(from_bson::<SocialService>(mongodb::Bson::Document(doc.clone())).unwrap());
			 	},
			 	Err(_) => break
			 }
		}
		// TODO: no
		if all_documents.len() == 0 {
			return Err(HandlerError::Error("no socials".to_string()))
		}
		Ok(all_documents)		
	}

	pub fn get_social_service(&self, channel: &str, service: &str) -> HandlerResult<SocialService> {
		let filter = doc! {
			"channel": channel,
			"service": service
		};

		let db = self.database.as_ref().expect("no database");
		let social_collection = db.collection("socials");

		let document = social_collection.find_one(Some(filter), None).map_err(|e| HandlerError::DatabaseError(e))?;
		match document {
			Some(doc) => Ok(from_bson::<SocialService>(mongodb::Bson::Document(doc.clone())).unwrap()),
			_ => Err(HandlerError::Error("no social service".to_string()))
		}
	}

	pub fn create_social_service(&self, channel: &str, service: &str, url: &str) -> HandlerResult<SocialService> {
		let db = self.database.as_ref().expect("no database");
		let social_collection = db.collection("socials");

		let social = SocialService {
			channel: channel.to_string(),
			service: service.to_string(),
			url: url.to_string()
		};

		social_collection.insert_one(to_bson(&social).unwrap().as_document().unwrap().clone(), None)
			.map_err(|e| HandlerError::DatabaseError(e))?;
		Ok(social)
	}

	pub fn remove_social(&self, channel: &str, service: &str) -> HandlerResult<()> {
		let db = self.database.as_ref().expect("no database");
		let social_collection = db.collection("socials");

		social_collection.delete_one(doc! {
			"channel": channel,
			"service": service
		}, None).map_err(|e| HandlerError::DatabaseError(e))?;
		Ok(())
	}

	pub fn get_offences(&self, channel: &str, service: &str, user: &str) -> HandlerResult<UserOffences> {
		let filter = doc! {
			"channel": channel,
			"user": user,
			"service": service
		};

		let db = self.database.as_ref().expect("no database");
		let offences_collection = db.collection("offences");

		let document = offences_collection.find_one(Some(filter), None).map_err(|e| HandlerError::DatabaseError(e))?;
		match document {
			Some(doc) => Ok(from_bson::<UserOffences>(mongodb::Bson::Document(doc.clone())).unwrap()),
			_ => Err(HandlerError::Error("no offences".to_string()))
		}
	}

	pub fn get_offence(&self, channel: &str, service: &str, user: &str, ty: &str) -> HandlerResult<i32> {
		let filter = doc! {
			"channel": channel,
			"user": user,
			"service": service
		};

		let db = self.database.as_ref().expect("no database");
		let offences_collection = db.collection("offences");

		let document = offences_collection.find_one(Some(filter), None).map_err(|e| HandlerError::DatabaseError(e))?;
		match document {
			Some(doc) => Ok(from_bson::<UserOffences>(mongodb::Bson::Document(doc.clone())).unwrap().get_attribute(ty).unwrap_or(0)),
			_ => Err(HandlerError::Error("no offences".to_string()))
		}
	}

	pub fn create_offence(&self, channel: &str, user: &str, service: &str) -> HandlerResult<UserOffences> {
		let db = self.database.as_ref().expect("no database");
		let offences_collection = db.collection("offences");

		let offence = UserOffences {
			channel: channel.to_string(),
			service: service.to_string(),
			user: user.to_string(),
			caps: 0,
			emoji: 0,
			urls: 0
		};

		offences_collection.insert_one(to_bson(&offence).unwrap().as_document().unwrap().clone(), None)
			.map_err(|e| HandlerError::DatabaseError(e))?;
		Ok(offence)
	}

	fn get_or_create_offence(&self, channel: &str, user: &str, service: &str) -> UserOffences {
		self.get_offences(channel, user, service).unwrap_or_else(|_| self.create_offence(channel, user, service).unwrap())
	}

	pub fn update_offence(&self, channel: &str, user: &str, service: &str, offence_type: &str, count: UpdateCount) -> HandlerResult<UserOffences> {
		let db = self.database.as_ref().expect("no database");
		let offences_collection = db.collection("offences");
		let mut offence = self.get_or_create_offence(channel, user, service);
		let mut current = offence.get_attribute(offence_type).ok_or(HandlerError::Error("invalid offence type".to_string()))?;

		let (operator, remaining) = count.count.split_at(1);
		let remaining = remaining.parse::<i32>().map_err(|_| HandlerError::Error("invalid count".to_string()))?;

		match operator {
			"+" => current += remaining,
			"-" => current -= remaining,
			"@" => current = remaining,
			_ => return Err(HandlerError::Error("invalid operator".to_string()))
		}
		offence = offence.set_attribute(offence_type, current).map_err(|_| HandlerError::Error("invalid offence type".to_string()))?;

		match offences_collection.update_one(doc! {
			"channel": channel,
			"service": service,
			"user": user
		}, doc! {
			"$set": bson::to_bson(&offence).unwrap()
		}, None) {
			Ok(_) => Ok(offence),
			Err(e) => Err(HandlerError::DatabaseError(e))
		}
	}
}
