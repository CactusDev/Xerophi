
use rocket_contrib::json::{JsonValue, Json};

use rocket::State;
use crate::{
	DbConn, endpoints::generate_error,
	database::{
		handler::HandlerError,
		structures::Message
	}
};

pub struct PostQuote {
	pub response: Vec<Message>,
	pub channel: String
}
